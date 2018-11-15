package filesystem

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
	"encoding/json"
	"golang.org/x/net/context"
	"gopkg.in/src-d/go-billy.v4"
	"indeed/gophers/rlog"
	"io"
	"os"
	"reflect"
	"sync"
	"syscall"
)

const defaultPerms = 0755

func debug(obj interface{}, method, path string, req interface{}) {
	data, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	rlog.Infof(
		"%s#%s [path=%s] [req=%s]",
		reflect.TypeOf(obj),
		method, path,
		string(data),
	)
}

// directories

type BillyDirectory struct {
	uid uint32
	gid uint32
	path string
	fs   billy.Filesystem

	mu sync.Mutex
}

func (b *BillyDirectory) Rename(ctx context.Context, req *fuse.RenameRequest, newDir fs.Node) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	debug(b, "Rename", b.path, req)

	oldFullPath := b.fs.Join(b.path, req.OldName)
	newFullPath := b.fs.Join(b.path, req.NewName)

	return b.fs.Rename(oldFullPath, newFullPath)
}

func (b *BillyDirectory) Symlink(ctx context.Context, req *fuse.SymlinkRequest) (fs.Node, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	debug(b, "Symlink", b.path, req)

	fullPath := b.fs.Join(b.path, req.NewName)

	err := b.fs.Symlink(req.Target, fullPath)

	if err != nil {
		rlog.Errorf("failed to symlink %s to %s, %v", b.path, req.Target, err)
		return nil, fuse.EPERM
	}

	return &Symlink{
		uid: b.uid,
		gid: b.gid,
		path: fullPath,
		fs: b.fs,
		dir: true,
	}, nil
}

func (b *BillyDirectory) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	debug(b, "Remove", b.path, req)

	fullPath := b.fs.Join(b.path, req.Name)

	if err := b.fs.Remove(fullPath); err != nil {
		rlog.Errorf("failed to remove node at path: %s, %v", fullPath, err)
		return fuse.EPERM
	}

	return nil
}

func (b *BillyDirectory) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	debug(b, "Create", b.path, req)

	fullPath := b.fs.Join(b.path, req.Name)

	_, file, err := createOrOpenFile(b.fs, fullPath)
	if err != nil {
		rlog.Errorf("failed to open file for writing: %s", fullPath)
		return nil, nil, fuse.EPERM
	}

	f := &BillyFile{
		uid: b.uid,
		gid: b.gid,
		path:     fullPath,
		fs:       b.fs,
		file:     file,
		refcount: 1,
	}

	return f, f, nil
}

func (b *BillyDirectory) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	debug(b, "Mkdir", b.path, req)

	fullPath := b.fs.Join(b.path, req.Name)

	if err := b.fs.MkdirAll(fullPath, defaultPerms); err != nil {
		rlog.Errorf("failed to mkdir for path: %s, %v", fullPath, err)
		return nil, fuse.EPERM
	}

	return &BillyDirectory{
		uid: b.uid,
		gid: b.gid,
		path: fullPath,
		fs:   b.fs,
	}, nil
}

func (b *BillyDirectory) Lookup(ctx context.Context, name string) (fs.Node, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	debug(b, "Lookup", b.path, name)

	fullPath := b.fs.Join(b.path, name)

	if _, err := b.fs.Readlink(fullPath); err == nil {
		return &Symlink{
			uid: b.uid,
			gid: b.gid,
			path: fullPath,
			fs: b.fs,
		}, nil
	}

	finfo, err := b.fs.Stat(fullPath)
	if err != nil {
		return nil, fuse.ENOENT
	}

	if finfo.IsDir() {
		return &BillyDirectory{
			uid: b.uid,
			gid: b.gid,
			path: fullPath,
			fs:   b.fs,
		}, nil
	}

	return &BillyFile{
		uid: b.uid,
		gid: b.gid,
		path:     fullPath,
		fs:       b.fs,
		mu:       sync.Mutex{},
		refcount: 0,
	}, nil
}

func (b *BillyDirectory) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	debug(b, "ReadDirAll", b.path, nil)

	finfos, err := b.fs.ReadDir(b.path)

	if err != nil {
		rlog.Errorf("failed to readdir: %s, %v", b.path, err)
		return nil, fuse.EPERM
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

func (b *BillyDirectory) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Uid = b.uid
	attr.Gid = b.gid
	attr.Mode = os.ModeDir | defaultPerms

	return nil
}

// files

type BillyFile struct {
	uid uint32
	gid uint32
	path string
	fs   billy.Filesystem

	mu   sync.Mutex
	data []byte
	file billy.File

	refcount uint
}

func createOrOpenFile(fs billy.Filesystem, path string) (os.FileInfo, billy.File, error) {
	finfo, err := fs.Stat(path)
	if err != nil {
		// file does not exist, create it
		if _, err := fs.Create(path); err != nil {
			rlog.Errorf("failed to create file at path: %s, %v", path, err)
			return nil, nil, fuse.EPERM
		}
	}

	file, err := fs.OpenFile(path, os.O_RDWR, defaultPerms)

	if err != nil {
		rlog.Errorf("failed to open file at path: %s, %v", path, err)
		return nil, nil, fuse.EPERM
	}

	return finfo, file, nil
}

func (b *BillyFile) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	debug(b, "Open", b.path, req)

	if b.file == nil {
		finfo, file, err := createOrOpenFile(b.fs, b.path)
		if err != nil {
			rlog.Error("failed to createOrOpenFile: %s, %v", b.path, err)
			return nil, err
		}

		data := make([]byte, finfo.Size())

		if finfo.Size() > 0 {
			if _, err := file.Read(data); err != nil {
				rlog.Errorf("failed to read data from file: %s, %d, %v", b.path, finfo.Size(), err)
				return nil, fuse.EPERM
			}
		}

		b.data = data
		b.file = file
	}

	b.refcount++
	return b, nil
}

func (b *BillyFile) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	if b.file == nil {
		// nothing to release
		return nil
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	debug(b, "Release", b.path, req)

	b.refcount--
	if b.refcount == 0 {
		b.data = nil
		b.file = nil
	}

	return nil
}

func (b *BillyFile) Attr(ctx context.Context, attr *fuse.Attr) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	attr.Uid = b.uid
	attr.Gid = b.gid
	attr.Size = 0
	attr.Mode = defaultPerms

	if b.file == nil {
		info, err := b.fs.Stat(b.path)
		if err == nil {
			attr.Size = uint64(info.Size())
		}
	} else {
		attr.Size = uint64(len(b.data))
	}

	return nil
}

func (b *BillyFile) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	debug(b, "Read", b.path, req)

	fuseutil.HandleRead(req, resp, b.data)

	return nil
}

const maxInt = int(^uint(0) >> 1)

func (b *BillyFile) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	debug(b, "Write", b.path, req)

	// expand the buffer if necessary
	newLen := req.Offset + int64(len(req.Data))
	if newLen > int64(maxInt) {
		return fuse.Errno(syscall.EFBIG)
	}

	if newLen := int(newLen); newLen > len(b.data) {
		b.data = append(b.data, make([]byte, newLen-len(b.data))...)
	}

	n := copy(b.data[req.Offset:], req.Data)
	resp.Size = n
	return nil
}

func (b *BillyFile) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	debug(b, "Flush", b.path, req)

	_, err := b.file.Seek(0, io.SeekStart)

	if err != nil {
		rlog.Errorf("failed to seek to start of file: %v", err)
		return fuse.EPERM
	}

	n, err := b.file.Write(b.data)

	if err != nil {
		rlog.Errorf("failed to write data to file: %v", err)
		return fuse.EPERM
	}

	if n != len(b.data) {
		rlog.Error("failed to write all data to file")
		return fuse.EPERM
	}

	return nil
}

func (b *BillyFile) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if req.Valid.Size() {
		if req.Size > uint64(maxInt) {
			return fuse.Errno(syscall.EFBIG)
		}

		newLen := int(req.Size)

		switch {
		case newLen > len(b.data):
			b.data = append(b.data, make([]byte, newLen-len(b.data))...)
		case newLen < len(b.data):
			b.data = b.data[:newLen]
		}
	}

	return nil
}

// simple symlink struct

type Symlink struct {
	uid uint32
	gid uint32
	path string
	fs billy.Filesystem
	dir bool
}

func (s *Symlink) Readlink(ctx context.Context, req *fuse.ReadlinkRequest) (string, error) {
	debug(s, "Readlink", s.path, req)

	link, err := s.fs.Readlink(s.path)
	return link, err
}

func (s *Symlink) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Uid = s.uid
	attr.Gid = s.gid
	attr.Size = uint64(len(s.path))
	attr.Mode = os.ModeSymlink | defaultPerms

	if s.dir {
		attr.Mode = attr.Mode | os.ModeDir
	}

	return nil
}
