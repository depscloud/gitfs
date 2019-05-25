package filesystem

import (
	"bazil.org/fuse/fs"
	"github.com/deps-cloud/gitfs/pkg/tree"
	"github.com/deps-cloud/gitfs/pkg/urls"
)

// FileSystem defines the root of the directory tree.
type FileSystem struct {
	UID  uint32
	GID  uint32
	Tree *gitfstree.TreeNode
	FSA  *urls.FileSystemAdapter
}

var _ fs.FS = (*FileSystem)(nil)

// Root returns the root of the directory tree.
func (f *FileSystem) Root() (fs.Node, error) {
	return NewDirectory(f.UID, f.GID, f.Tree, f.FSA), nil
}
