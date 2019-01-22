package sync

import (
	"os"
	"sync"

	"gopkg.in/src-d/go-billy.v4"
)

var _ billy.Filesystem = &Filesystem{}
var _ billy.Change

// Filesystem implements a billy.Filesystem, but adds a mutex around all operations.
// The mutex around operations. The in-memory implementation is unsafe for concurrent use, so this
// wrapper makes it easy to synchronize
type Filesystem struct {
	sync.Mutex

	Delegate billy.Filesystem
}

// Create description.
func (f *Filesystem) Create(filename string) (billy.File, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Create(filename)
}

// Open description.
func (f *Filesystem) Open(filename string) (billy.File, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Open(filename)
}

// OpenFile description.
func (f *Filesystem) OpenFile(filename string, flag int, perm os.FileMode) (billy.File, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.OpenFile(filename, flag, perm)
}

// Stat description.
func (f *Filesystem) Stat(filename string) (os.FileInfo, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Stat(filename)
}

// Rename description.
func (f *Filesystem) Rename(oldpath, newpath string) error {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Rename(oldpath, newpath)
}

// Remove description.
func (f *Filesystem) Remove(filename string) error {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Remove(filename)
}

// Join description.
func (f *Filesystem) Join(elem ...string) string {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Join(elem...)
}

// TempFile description.
func (f *Filesystem) TempFile(dir, prefix string) (billy.File, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.TempFile(dir, prefix)
}

// ReadDir description.
func (f *Filesystem) ReadDir(path string) ([]os.FileInfo, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.ReadDir(path)
}

// MkdirAll description.
func (f *Filesystem) MkdirAll(filename string, perm os.FileMode) error {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.MkdirAll(filename, perm)
}

// Lstat description.
func (f *Filesystem) Lstat(filename string) (os.FileInfo, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Lstat(filename)
}

// Symlink description.
func (f *Filesystem) Symlink(target, link string) error {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Symlink(target, link)
}

// Readlink description.
func (f *Filesystem) Readlink(link string) (string, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Readlink(link)
}

// Chroot description.
func (f *Filesystem) Chroot(path string) (billy.Filesystem, error) {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Chroot(path)
}

// Root description.
func (f *Filesystem) Root() string {
	f.Lock()
	defer f.Unlock()

	return f.Delegate.Root()
}
