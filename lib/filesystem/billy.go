package filesystem

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"gopkg.in/src-d/go-billy.v4"
	"indeed/gophers/rlog"
	"math"
	"os"
	"sync"
	"syscall"
)

// node functions
var _ fs.Node = &BillyNode{}          // Attr
var _ fs.NodeSetattrer = &BillyNode{} // Setattr

// directory functions
var _ fs.NodeStringLookuper = &BillyNode{} // Lookup
var _ fs.HandleReadDirAller = &BillyNode{} // HandleReadDirAller
var _ fs.NodeMkdirer = &BillyNode{}        // Mkdir
var _ fs.NodeCreater = &BillyNode{}        // Create
var _ fs.NodeRemover = &BillyNode{}        // Remove
var _ fs.NodeRenamer = &BillyNode{}        // Rename
var _ fs.NodeSymlinker = &BillyNode{}      // Symlink

// handle functions
var _ fs.NodeOpener = &BillyNode{}     // Open
var _ fs.HandleWriter = &BillyNode{}   // Write
var _ fs.HandleReader = &BillyNode{}   // Read
var _ fs.HandleFlusher = &BillyNode{}  // Flush
var _ fs.HandleReleaser = &BillyNode{} // Release

// symlink functions
var _ fs.NodeReadlinker = &BillyNode{} // Readlink

const defaultPerms = 0755
const allPerms = 0777
const maxFileSize = math.MaxUint64
const createFileFlags = os.O_RDWR | os.O_CREATE | os.O_TRUNC
const maxInt = uint64(int(^uint(0) >> 1))

type BillyUser struct {
	uid uint32
	gid uint32
}

type BillyNode struct {
	// common between directories and files
	bfs  billy.Filesystem
	path string

	// used only for symlinks
	target string

	// metadata about the underlying file / directory
	user BillyUser
	mode os.FileMode

	// data for files
	size uint64
	data []byte

	// support file level locking
	mu sync.Mutex

	// node cache for re-use
	cache map[string]*BillyNode
}

// symlink functions

func (n *BillyNode) Readlink(ctx context.Context, req *fuse.ReadlinkRequest) (string, error) {
	n.debug("Readlink", req)

	if !n.isSymlink() {
		return "", fuse.Errno(syscall.EINVAL)
	}

	return n.target, nil
}

// handle functions

func (n *BillyNode) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	n.debug("Release", req)

	n.mu.Lock()
	defer n.mu.Unlock()

	if n.data == nil {
		// nothing to release
		return nil
	}

	return nil
}

func (n *BillyNode) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	n.debug("Flush", req)

	if !n.isRegular() {
		return fuse.Errno(syscall.EINVAL)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	file, err := n.bfs.OpenFile(n.path, os.O_WRONLY, defaultPerms)
	if err != nil {
		return err
	}

	if _, err := file.Write(n.data); err != nil {
		return err
	}

	return nil
}

func (n *BillyNode) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	n.debug("Read", req)

	if !n.isRegular() {
		return fuse.Errno(syscall.EINVAL)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	if err := n.load(); err != nil {
		return err
	}

	fuseutil.HandleRead(req, resp, n.data)

	return nil
}

func (n *BillyNode) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	n.debug("Write", req)

	if !n.isRegular() {
		return fuse.Errno(syscall.EINVAL)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	newLen := uint64(req.Offset) + uint64(len(req.Data))
	if newLen > maxInt {
		return fuse.Errno(syscall.EFBIG)
	}

	if newLen > n.size {
		n.data = append(n.data, make([]byte, newLen-n.size)...)
		n.size = newLen
	}

	resp.Size = copy(n.data[req.Offset:], req.Data)
	return nil
}

func (n *BillyNode) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	n.debug("Open", req)

	if n.isSymlink() {
		return nil, fuse.Errno(syscall.EINVAL)
	}

	if n.isRegular() {
		n.mu.Lock()
		defer n.mu.Unlock()

		// force data
		if err := n.load(); err != nil {
			return nil, err
		}
	}

	return n, nil
}

// directory functions

func (n *BillyNode) Symlink(ctx context.Context, req *fuse.SymlinkRequest) (fs.Node, error) {
	n.debug("Symlink", req)

	if !n.isDir() {
		return nil, fuse.Errno(syscall.ENOTDIR)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	if node, ok := n.cache[req.NewName]; ok {
		return node, nil
	}

	fullPath := n.bfs.Join(n.path, req.NewName)

	if err := n.bfs.Symlink(req.Target, fullPath); err != nil {
		// assumes it already exists
		return nil, fuse.EEXIST
	}

	node := &BillyNode{
		bfs:    n.bfs,
		path:   fullPath,
		target: req.Target,
		user:   n.user,
		mode:   os.ModeSymlink | defaultPerms,
		size:   uint64(len(req.Target)),
		data:   nil,
		mu:     sync.Mutex{},
	}

	n.cache[req.NewName] = node
	return node, nil
}

func (n *BillyNode) Rename(ctx context.Context, req *fuse.RenameRequest, newDir fs.Node) error {
	n.debug("Rename", req)

	newBillyNode, ok := newDir.(*BillyNode)
	if !ok {
		return fmt.Errorf("newDir is not a BillyNode")
	}

	if !n.isDir() || !newBillyNode.isDir() {
		return fuse.Errno(syscall.ENOTDIR)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	oldFullPath := n.bfs.Join(n.path, req.OldName)
	newFullPath := n.bfs.Join(newBillyNode.path, req.NewName)

	return n.bfs.Rename(oldFullPath, newFullPath)
}

func (n *BillyNode) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	n.debug("Remove", req)

	if !n.isDir() {
		return fuse.Errno(syscall.ENOTDIR)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	fullPath := n.bfs.Join(n.path, req.Name)

	if err := n.bfs.Remove(fullPath); err != nil {
		return fuse.ENOENT
	}

	return nil
}

func (n *BillyNode) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	n.debug("Create", req)

	if !n.isDir() {
		return nil, nil, fuse.Errno(syscall.ENOTDIR)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	node, err := n.createOrOpen(req.Name, true)
	if err != nil {
		return nil, nil, err
	}

	// force load the data
	// ensure proper file handle
	node.load()

	return node, node, err
}

func (n *BillyNode) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	n.debug("Mkdir", req)

	if !n.isDir() {
		return nil, fuse.Errno(syscall.ENOTDIR)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	if node, ok := n.cache[req.Name]; ok {
		return node, nil
	}

	fullPath := n.bfs.Join(n.path, req.Name)

	if err := n.bfs.MkdirAll(fullPath, defaultPerms); err != nil {
		return nil, fuse.EEXIST
	}

	node := &BillyNode{
		bfs:    n.bfs,
		path:   fullPath,
		target: "",
		user:   n.user,
		mode:   os.ModeDir | defaultPerms,
		size:   0,
		data:   nil,
		mu:     sync.Mutex{},
		cache:  make(map[string]*BillyNode),
	}

	n.cache[req.Name] = node
	return node, nil
}

func (n *BillyNode) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	n.debug("ReadDirAll", nil)

	if !n.isDir() {
		return nil, fuse.Errno(syscall.ENOTDIR)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	finfos, err := n.bfs.ReadDir(n.path)
	if err != nil {
		return nil, fuse.ENOENT
	}

	dirents := make([]fuse.Dirent, len(finfos))

	for i := 0; i < len(finfos); i++ {
		finfo := finfos[i]
		dirents[i] = fuse.Dirent{
			Type: fuse.DT_Unknown,
			Name: finfo.Name(),
		}
	}

	return dirents, nil
}

func (n *BillyNode) Lookup(ctx context.Context, name string) (fs.Node, error) {
	n.debug("Lookup", name)

	if !n.isDir() {
		return nil, fuse.Errno(syscall.ENOTDIR)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	node, err := n.createOrOpen(name, false)

	return node, err
}

// node functions

func (n *BillyNode) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
	n.debug("Setattr", req)

	if !req.Valid.Size() {
		// only support setting the file size
		return nil
	}

	if !n.isRegular() {
		// Setting the size is only available on regular files.
		return fuse.Errno(syscall.EINVAL)
	}

	if req.Size > maxFileSize {
		// Too big.
		return fuse.Errno(syscall.EFBIG)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	if req.Size == n.size {
		return nil // Nothing to do.
	}

	// data not loaded, no need to resize array
	if n.data != nil {
		if req.Size > n.size {
			n.data = append(n.data, make([]byte, req.Size-n.size)...)
		} else {
			n.data = n.data[:req.Size]
		}
	}

	n.size = req.Size

	return nil
}

func (n *BillyNode) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Uid = n.user.uid
	attr.Gid = n.user.gid
	attr.Mode = n.mode

	if n.isRegular() {
		n.mu.Lock()
		defer n.mu.Unlock()

		attr.Size = n.size
	} else if n.isSymlink() {
		attr.Size = uint64(len(n.target))
	}

	return nil
}

// helper functions

func (n *BillyNode) createOrOpen(name string, create bool) (*BillyNode, error) {
	n.debug("createOrOpen", name)

	if node, ok := n.cache[name]; ok {
		return node, nil
	}

	fullPath := n.bfs.Join(n.path, name)

	// symlink
	if target, err := n.bfs.Readlink(fullPath); err == nil {
		node := &BillyNode{
			bfs:    n.bfs,
			path:   fullPath,
			target: target,
			user:   n.user,
			mode:   os.ModeSymlink | defaultPerms,
			size:   uint64(len(target)),
			data:   nil,
			mu:     sync.Mutex{},
			cache:  make(map[string]*BillyNode),
		}
		n.cache[name] = node
		return node, nil
	}

	finfo, err := n.bfs.Stat(fullPath)
	if err == nil {
		// file exists, create reference
		node := &BillyNode{
			bfs:    n.bfs,
			path:   fullPath,
			target: "",
			user:   n.user,
			mode:   finfo.Mode(),
			size:   uint64(finfo.Size()),
			data:   nil,
			mu:     sync.Mutex{},
			cache:  make(map[string]*BillyNode),
		}

		n.cache[name] = node
		return node, nil
	} else if !create {
		// file does not exist, not creating
		return nil, fuse.ENOENT
	}

	// don't use bfs.Create since it assigns 666 permissions
	if _, err := n.bfs.OpenFile(fullPath, createFileFlags, defaultPerms); err != nil {
		// shouldn't really happen but lets just account for it just in case
		return nil, fuse.EEXIST
	}

	node := &BillyNode{
		bfs:    n.bfs,
		path:   fullPath,
		target: "",
		user:   n.user,
		mode:   defaultPerms,
		size:   0,
		data:   make([]byte, 0),
		mu:     sync.Mutex{},
		cache:  make(map[string]*BillyNode),
	}

	n.cache[name] = node
	return node, nil
}

func (n *BillyNode) load() error {
	if n.data != nil {
		// already loaded, nothing to do
		return nil
	}

	n.debug("load", nil)

	file, err := n.bfs.OpenFile(n.path, os.O_RDONLY, defaultPerms)
	if err != nil {
		n.error("load", err)
		return fuse.ENOENT
	}

	data := make([]byte, n.size)
	if n.size > 0 {
		if _, err := file.Read(data); err != nil {
			n.error("load", err)
			return fuse.EPERM
		}
	}

	n.data = data
	return nil
}

func (n *BillyNode) isDir() bool {
	return n.mode.IsDir()
}

func (n *BillyNode) isRegular() bool {
	return n.mode.IsRegular()
}

func (n *BillyNode) isSymlink() bool {
	return n.mode&os.ModeSymlink != 0
}

func (n *BillyNode) error(method string, err error) {
	rlog.Errorf("[BillyNode#%s] [path=%s] %v", method, n.path, err)
}

func (n *BillyNode) debug(method string, req interface{}) {
	reqData, _ := json.Marshal(req)

	rlog.Infof(
		"[BillyNode#%s] [mode=%s, path=%s] [req=%s]",
		method, n.mode, n.path, string(reqData),
	)
}
