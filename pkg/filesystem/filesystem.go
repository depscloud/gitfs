package filesystem

import (
	"bazil.org/fuse/fs"
	"github.com/mjpitz/gitfs/pkg/clone"
	"github.com/mjpitz/gitfs/pkg/tree"
)

type FileSystem struct {
	Uid    uint32
	Gid    uint32
	Tree   *gitfstree.TreeNode
	Cloner *clone.Cloner
}

var _ fs.FS = (*FileSystem)(nil)

func (f *FileSystem) Root() (fs.Node, error) {
	return NewDirectory(f.Uid, f.Gid, f.Tree, f.Cloner), nil
}
