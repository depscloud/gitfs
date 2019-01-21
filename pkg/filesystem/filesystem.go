package filesystem

import (
	"bazil.org/fuse/fs"
	"github.com/mjpitz/gitfs/pkg/tree"
	"github.com/mjpitz/gitfs/pkg/urls"
)

type FileSystem struct {
	Uid  uint32
	Gid  uint32
	Tree *gitfstree.TreeNode
	FSA  *urls.FileSystemAdapter
}

var _ fs.FS = (*FileSystem)(nil)

func (f *FileSystem) Root() (fs.Node, error) {
	return NewDirectory(f.Uid, f.Gid, f.Tree, f.FSA), nil
}
