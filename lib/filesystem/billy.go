package filesystem

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
	"gopkg.in/src-d/go-billy.v4"
	"indeed/gophers/rlog"
	"io"
	"math"
	"os"
	"strings"
	"sync"
)

// directories

type BillyDirectory struct {
	path string
	fs   billy.Filesystem
}

func (b *BillyDirectory) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	fullPath := strings.Join([]string{
		b.path,
		req.Name,
	}, string(os.PathSeparator))

	rlog.Infof("attempting to remove: %s", fullPath)

	if err := b.fs.Remove(fullPath); err != nil {
		rlog.Errorf("failed to remove node at path: %s, %v", fullPath, err)
		return fuse.EPERM
	}

	return nil
}

func (b *BillyDirectory) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	fullPath := strings.Join([]string{
		b.path,
		req.Name,
	}, string(os.PathSeparator))

	rlog.Infof("attempting to create: %s", req.Name)

	file, err := createOrOpenFile(b.fs, fullPath)
	if err != nil {
		rlog.Errorf("failed to open file for writing: %s", fullPath)
		return nil, nil, fuse.EPERM
	}

	f := &BillyFile{
		path: fullPath,
		fs:   b.fs,
		file: file,
		refcount: 1,
	}

	return f, f, nil
}

func (b *BillyDirectory) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	fullPath := strings.Join([]string{
		b.path,
		req.Name,
	}, string(os.PathSeparator))

	rlog.Infof("attempting to make directory: %s", fullPath)

	if err := b.fs.MkdirAll(fullPath, 0755); err != nil {
		rlog.Errorf("failed to mkdir for path: %s, %v", fullPath, err)
		return nil, fuse.EPERM
	}

	return &BillyDirectory{
		path: fullPath,
		fs:   b.fs,
	}, nil
}

func (b *BillyDirectory) Lookup(ctx context.Context, name string) (fs.Node, error) {
	fullPath := strings.Join([]string{
		b.path,
		name,
	}, string(os.PathSeparator))

	rlog.Infof("attempting to lookup: %s", fullPath)

	finfo, err := b.fs.Stat(fullPath)
	if err != nil {
		return nil, fuse.ENOENT
	}

	if finfo.IsDir() {
		return &BillyDirectory{
			path: fullPath,
			fs:   b.fs,
		}, nil
	}

	return &BillyFile{
		path:     fullPath,
		fs:       b.fs,
		mu:       sync.Mutex{},
		refcount: 0,
	}, nil
}

func (b *BillyDirectory) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	rlog.Infof("attempting to readdir: %s", b.path)

	finfos, err := b.fs.ReadDir(b.path)

	if err != nil {
		rlog.Errorf("failed toreaddir: %s, %v", b.path, err)
		return nil, fuse.EPERM
	}

	dirents := make([]fuse.Dirent, len(finfos))
	for i := 0; i < len(finfos); i++ {
		finfo := finfos[i]

		direntType := fuse.DT_File
		if finfo.IsDir() {
			direntType = fuse.DT_Dir
		}

		dirents[i] = fuse.Dirent{
			Type: direntType,
			Name: finfo.Name(),
		}
	}

	return dirents, nil
}

func (b *BillyDirectory) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Mode = os.ModeDir | 0755

	return nil
}

// files

type BillyFile struct {
	path string
	fs   billy.Filesystem

	mu   sync.Mutex
	file billy.File

	refcount uint
}

func createOrOpenFile(fs billy.Filesystem, path string) (billy.File, error) {
	if _, err := fs.Stat(path); err != nil {
		// file does not exist, create it
		if _, err := fs.Create(path); err != nil {
			rlog.Errorf("failed to create file at path: %s, %v", path, err)
			return nil, fuse.EPERM
		}
	}

	file, err := fs.OpenFile(path, os.O_RDWR, 0755)

	if err != nil {
		rlog.Errorf("failed to open file at path: %s, %v", path, err)
		return nil, fuse.EPERM
	}

	return file, nil
}

func (b *BillyFile) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	rlog.Info("attempting to flush")

	return nil
}

func (b *BillyFile) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	rlog.Infof("attempting to open: %s", b.path)

	if b.file == nil {
		file, err := createOrOpenFile(b.fs, b.path)
		if err != nil {
			return nil, err
		}
		b.file = file
	}

	b.refcount++
	return b, nil
}

func (b *BillyFile) Attr(ctx context.Context, attr *fuse.Attr) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	info, err := b.fs.Stat(b.path)
	if err == nil {
		attr.Size = uint64(info.Size())
	} else {
		attr.Size = 0
	}

	attr.Mode = info.Mode()
	return nil
}

func (b *BillyFile) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	rlog.Infof("attempting to write: %s [offset=%d, length=%d]", b.path, req.Offset, len(req.Data))

	n, err := b.file.Write(req.Data)

	if err != nil {
		rlog.Errorf("failed to write data to file; %v", err)
		return fuse.EPERM
	}

	if n != len(req.Data) {
		rlog.Errorf("failed to write all data to file")
		return fuse.EPERM
	}

	resp.Size = n

	return nil
}

func (b *BillyFile) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	rlog.Infof("attempting to read: %s [offset=%d, length=%d]", b.path, req.Offset, req.Size)

	// seek to the start of the file
	// todo: refactor this to buffer in memory instead of reading from disk all the time
	pos, err := b.file.Seek(0, io.SeekStart)
	if err != nil || pos != 0 {
		return fuse.EPERM
	}

	// handle the case where we request a larger readableSize then what is available
	finfo, err := b.fs.Stat(b.path)
	if err != nil {
		return fuse.EPERM
	}

	reqSize := float64(req.Size)
	finfoSize := float64(finfo.Size())

	readableSize := int(math.Min(reqSize, finfoSize))

	bytes := make([]byte, readableSize)
	if _, err := b.file.ReadAt(bytes, req.Offset); err != nil {
		rlog.Errorf("failed to read data from file; %v", err)
		return fuse.EPERM
	}

	resp.Data = make([]byte, req.Size)
	copy(resp.Data, bytes)

	return nil
}

func (b *BillyFile) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	if b.file == nil {
		// nothing to release
		return nil
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.refcount--
	if b.refcount == 0 {
		b.file = nil
	}

	return nil
}
