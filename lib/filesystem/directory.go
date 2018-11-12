package filesystem

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/indeedeng/gitfs/lib/tree"
	"golang.org/x/net/context"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	gitmemory "gopkg.in/src-d/go-git.v4/storage/memory"
	"indeed/gophers/rlog"
	"os"
)

type Directory struct {
	tree *gitfstree.TreeNode
}

func (d *Directory) Lookup(ctx context.Context, name string) (fs.Node, error) {
	node, ok := d.tree.Child(name)
	if !ok {
		return nil, fuse.ENOENT
	}

	if len(node.URL) > 0 {
		storage := gitmemory.NewStorage()
		fileSystem := memfs.New()

		rlog.Infof("cloning repository: %s", node.URL)

		// shallow clone for now since we only support read only
		repository, err := git.Clone(storage, fileSystem, &git.CloneOptions{
			URL:   node.URL,
			Depth: 1,
		})

		if err != nil {
			rlog.Error("failed to clone repository: %v", err)
			return nil, fuse.ENOENT
		}

		wt, err := repository.Worktree()

		if err != nil {
			rlog.Error("failed to obtain work tree: %v", err)
			return nil, fuse.ENOENT
		}

		return &BillyDirectory{
			path: "",
			fs:   wt.Filesystem,
		}, nil
	}

	return &Directory{tree: node}, nil
}

func (d *Directory) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
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
