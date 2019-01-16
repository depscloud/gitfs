package filesystem

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sync"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gopkg.in/src-d/go-billy.v4"
)

var _ INode = &BillyNode{}

const defaultPerms = 0644
const maxFileSize = math.MaxUint64
const createFileFlags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
const maxInt = uint64(int(^uint(0) >> 1))

type BillyUser struct {
	uid uint32
	gid uint32
}

type BillyNode struct {
	// common between directories and files
	repourl string
	fs      billy.Filesystem
	path    string

	// used only for symlinks
	target string

	// metadata about the underlying file / directory
	user BillyUser
	mode os.FileMode

	// data for files
	size uint64
	data []byte

	// support file level locking
	mu *sync.Mutex
}

func (n *BillyNode) Fsync(ctx context.Context, req *fuse.FsyncRequest) error {
	n.debug("Fsync", req)
	// not quite sure what to do here, but it needs to be implemented and return nil.
	return nil
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

	file, err := n.fs.OpenFile(n.path, createFileFlags, n.mode)
	if err != nil {
		n.error("Flush", err)
		return err
	}

	if err = file.Lock(); err != nil {
		n.error("Flush", err)
		return fuse.ENOENT
	}

	defer func() {
		if err := file.Unlock(); err != nil {
			n.error("Flush", err)
		}

		if err := file.Close(); err != nil {
			n.error("Flush", err)
		}
	}()

	if _, err := file.Write(n.data); err != nil {
		n.error("Flush", err)
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

	fullPath := n.fs.Join(n.path, req.NewName)

	if err := n.fs.Symlink(req.Target, fullPath); err != nil {
		// assumes it already exists
		return nil, fuse.EEXIST
	}

	node := &BillyNode{
		repourl: n.repourl,
		fs:      n.fs,
		path:    fullPath,
		target:  req.Target,
		user:    n.user,
		mode:    os.ModeSymlink | defaultPerms,
		size:    uint64(len(req.Target)),
		data:    nil,
		mu:      &sync.Mutex{},
	}

	return node, nil
}

func (n *BillyNode) Rename(ctx context.Context, req *fuse.RenameRequest, newDir fs.Node) error {
	n.debug("Rename", req)

	newBillyNode, ok := newDir.(*BillyNode)
	if !ok {
		err := fmt.Errorf("newDir is not a BillyNode")
		n.error("Rename", err)
		return err
	}

	if !n.isDir() || !newBillyNode.isDir() {
		n.error("Rename", fmt.Errorf("n or newBillyNode is not a directory"))
		return fuse.Errno(syscall.ENOTDIR)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	oldFullPath := n.fs.Join(n.path, req.OldName)
	newFullPath := n.fs.Join(newBillyNode.path, req.NewName)

	err := n.fs.Rename(oldFullPath, newFullPath)

	if err != nil {
		n.error("Rename", err)
		return fuse.ENOENT
	}

	return nil
}

func (n *BillyNode) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	n.debug("Remove", req)

	if !n.isDir() {
		return fuse.Errno(syscall.ENOTDIR)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	// remove from the file system
	fullPath := n.fs.Join(n.path, req.Name)

	if err := n.fs.Remove(fullPath); err != nil {
		n.error("Remove", err)
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

	fullPath := n.fs.Join(n.path, req.Name)

	node := &BillyNode{
		repourl: n.repourl,
		fs:      n.fs,
		path:    fullPath,
		target:  "",
		user:    n.user,
		mode:    req.Mode,
		size:    uint64(0),
		data:    make([]byte, 0),
		mu:      &sync.Mutex{},
	}

	return node, node, nil
}

func (n *BillyNode) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	n.debug("Mkdir", req)

	if !n.isDir() {
		n.error("Mkdir", fmt.Errorf("not a directory"))
		return nil, fuse.Errno(syscall.ENOTDIR)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	fullPath := n.fs.Join(n.path, req.Name)

	if err := n.fs.MkdirAll(fullPath, defaultPerms); err != nil {
		n.error("Mkdir", err)
		return nil, fuse.EEXIST
	}

	node := &BillyNode{
		repourl: n.repourl,
		fs:      n.fs,
		path:    fullPath,
		target:  "",
		user:    n.user,
		mode:    os.ModeDir | defaultPerms,
		size:    0,
		data:    nil,
		mu:      &sync.Mutex{},
	}

	return node, nil
}

func (n *BillyNode) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	n.debug("ReadDirAll", nil)

	if !n.isDir() {
		return nil, fuse.Errno(syscall.ENOTDIR)
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	finfos, err := n.fs.ReadDir(n.path)
	if err != nil {
		n.error("ReadDirAll", err)
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

	fullPath := n.fs.Join(n.path, name)

	// symlink
	if target, err := n.fs.Readlink(fullPath); err == nil {
		node := &BillyNode{
			repourl: n.repourl,
			fs:      n.fs,
			path:    fullPath,
			target:  target,
			user:    n.user,
			mode:    os.ModeSymlink | 0755,
			size:    uint64(len(target)),
			data:    nil,
			mu:      &sync.Mutex{},
		}
		return node, nil
	}

	if finfo, err := n.fs.Stat(fullPath); err != nil {
		// assumed file does not exist
		return nil, fuse.ENOENT
	} else {
		return &BillyNode{
			repourl: n.repourl,
			fs:      n.fs,
			path:    fullPath,
			target:  "",
			user:    n.user,
			mode:    finfo.Mode(),
			size:    uint64(finfo.Size()),
			data:    nil,
			mu:      &sync.Mutex{},
		}, nil
	}
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

func (n *BillyNode) load() error {
	if n.data != nil {
		// already loaded, nothing to do
		return nil
	}

	n.debug("load", nil)

	file, err := n.fs.OpenFile(n.path, os.O_RDONLY, n.mode)
	if err != nil {
		n.error("load", err)
		return fuse.ENOENT
	}

	if err = file.Lock(); err != nil {
		n.error("load", err)
		return fuse.ENOENT
	}

	defer func() {
		if err := file.Unlock(); err != nil {
			n.error("load", err)
		}

		if err := file.Close(); err != nil {
			n.error("load", err)
		}
	}()

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
	logrus.Errorf(
		"[filesystem.billy] [repo=%s, path=%s] [BillyNode#%s] %v",
		n.repourl, n.path, method, err,
	)
}

func (n *BillyNode) debug(method string, req interface{}) {
	if os.Getenv("DEBUG") == "true" {
		reqData, _ := json.Marshal(req)

		logrus.Infof(
			"[filesystem.billy] [repo=%s, path=%s] [BillyNode#%s] [req=%s]",
			n.repourl, n.path, method, string(reqData),
		)
	}
}
