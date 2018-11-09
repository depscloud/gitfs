package filesystem

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
	"gopkg.in/src-d/go-billy.v4"
	"indeed/gophers/3rdparty/p/github.com/pkg/errors"
	"os"
	"strings"
	"sync"
)

// directories

type BillyDirectory struct {
	path string
	fs   billy.Filesystem
}

func (b *BillyDirectory) Lookup(ctx context.Context, name string) (fs.Node, error) {
	fullPath := strings.Join([]string{
		b.path,
		name,
	}, string(os.PathSeparator))

	finfo, err := b.fs.Stat(fullPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to stat file at path: %s", fullPath)
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
	finfos, err := b.fs.ReadDir(b.path)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to readdir: %s", b.path)
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

func (b *BillyFile) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.file == nil {
		file, err := b.fs.Open(b.path)

		if err != nil {
			return nil, errors.Wrapf(err, "failed to open file at path: %s", b.path)
		}

		b.file = file
	}

	b.refcount++
	return b, nil
}

func (b *BillyFile) Attr(ctx context.Context, attr *fuse.Attr) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// readonly
	info, err := b.fs.Stat(b.path)
	if err != nil {
		return nil
	}

	attr.Size = uint64(info.Size())
	attr.Mode = 0755

	return nil
}

func (b *BillyFile) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	n, err := b.file.Write(req.Data)

	if err != nil {
		return errors.Wrap(err, "failed to write data to file")
	}

	if n != len(req.Data) {
		return errors.Wrap(err, "failed to write all data to file")
	}

	return nil
}

func (b *BillyFile) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	bytes := make([]byte, req.Size)
	if _, err := b.file.ReadAt(bytes, req.Offset); err != nil {
		return errors.Wrap(err, "failed to read data from file")
	}

	resp.Data = bytes
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
