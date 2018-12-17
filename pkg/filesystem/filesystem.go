package filesystem

import (
	"bazil.org/fuse/fs"
	"github.com/mjpitz/gitfs/pkg/tree"
)

type FileSystem struct {
	Uid  uint32
	Gid  uint32
	Tree *gitfstree.TreeNode
}

var _ fs.FS = (*FileSystem)(nil)

func (fs *FileSystem) Root() (fs.Node, error) {
	return NewDirectory(fs.Uid, fs.Gid, fs.Tree), nil
}
