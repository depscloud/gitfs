package filesystem

import (
	"github.com/mjpitz/gitfs/pkg/urls"
	"os"
	"sync"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/mjpitz/gitfs/pkg/tree"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func NewDirectory(uid, gid uint32, tree *gitfstree.TreeNode, cloner *urls.FileSystemAdapter) fs.Node {
	return &Directory{
		uid:  uid,
		gid:  gid,
		tree: tree,
		fsa:  cloner,
		lock: &sync.Mutex{},
	}
}

type Directory struct {
	uid  uint32
	gid  uint32
	tree *gitfstree.TreeNode
	fsa  *urls.FileSystemAdapter
	lock *sync.Mutex
}

func (d *Directory) Lookup(ctx context.Context, name string) (fs.Node, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	node, ok := d.tree.Child(name)
	if !ok {
		return nil, fuse.ENOENT
	}

	var directory fs.Node

	if node.URL != nil {
		cloned, err := d.fsa.Clone(node.URL)
		if err != nil {
			logrus.Errorf("[filesystem.directory] failed to clone url: %s, %v", node.URL, err)
			return nil, fuse.ENOENT
		}

		directory = &BillyNode{
			repourl: node.URL.String(),
			fs:      cloned,
			path:    "",
			target:  "",
			user: BillyUser{
				uid: d.uid,
				gid: d.gid,
			},
			mode: os.ModeDir | defaultPerms,
			size: 0,
			data: nil,
			mu:   &sync.Mutex{},
		}
	} else {
		directory = NewDirectory(d.uid, d.gid, node, d.fsa)
	}

	return directory, nil
}

func (d *Directory) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	// tree is immutable, no need to lock
	children := d.tree.Children()

	dirents := make([]fuse.Dirent, len(children))

	for i, child := range children {
		dirents[i] = fuse.Dirent{
			Name: child,
		}
	}

	return dirents, nil
}

func (d *Directory) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Mode = os.ModeDir | 0755
	return nil
}
