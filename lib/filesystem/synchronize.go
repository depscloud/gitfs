package filesystem

import (
	"os"
	"sync"

	"gopkg.in/src-d/go-billy.v4"
)

var _ billy.Filesystem = &SynchronizedFilesystem{}

type SynchronizedFilesystem struct {
	sync.Mutex

	Delegate billy.Filesystem
}

func (f *SynchronizedFilesystem) Create(filename string) (billy.File, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Create(filename)
}

func (f *SynchronizedFilesystem) Open(filename string) (billy.File, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Open(filename)
}

func (f *SynchronizedFilesystem) OpenFile(filename string, flag int, perm os.FileMode) (billy.File, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.OpenFile(filename, flag, perm)
}

func (f *SynchronizedFilesystem) Stat(filename string) (os.FileInfo, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Stat(filename)
}

func (f *SynchronizedFilesystem) Rename(oldpath, newpath string) error {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Rename(oldpath, newpath)
}

func (f *SynchronizedFilesystem) Remove(filename string) error {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Remove(filename)
}

func (f *SynchronizedFilesystem) Join(elem ...string) string {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Join(elem...)
}

func (f *SynchronizedFilesystem) TempFile(dir, prefix string) (billy.File, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.TempFile(dir, prefix)
}

func (f *SynchronizedFilesystem) ReadDir(path string) ([]os.FileInfo, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.ReadDir(path)
}

func (f *SynchronizedFilesystem) MkdirAll(filename string, perm os.FileMode) error {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.MkdirAll(filename, perm)
}

func (f *SynchronizedFilesystem) Lstat(filename string) (os.FileInfo, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Lstat(filename)
}

func (f *SynchronizedFilesystem) Symlink(target, link string) error {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Symlink(target, link)
}

func (f *SynchronizedFilesystem) Readlink(link string) (string, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Readlink(link)
}

func (f *SynchronizedFilesystem) Chroot(path string) (billy.Filesystem, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Chroot(path)
}

func (f *SynchronizedFilesystem) Root() string {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Root()
}
