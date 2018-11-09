package filesystem

import (
	"bazil.org/fuse/fs"
	"github.com/indeedeng/gitfs/lib/tree"
)

type FileSystem struct {
	Tree *gitfstree.TreeNode
}

var _ fs.FS = (*FileSystem)(nil)

func (fs *FileSystem) Root() (fs.Node, error) {
	return &Directory{
		tree: fs.Tree,
	}, nil
}
