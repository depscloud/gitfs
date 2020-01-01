package urls

import (
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

var _ Cloner = &gitcloner{}

type gitcloner struct {
}

func (gc *gitcloner) Clone(url *URL, depth int, urlfs billy.Filesystem) error {
	gitfs, err := urlfs.Chroot(git.GitDirName)
	if err != nil {
		return errors.Wrap(err, "failed to create .git dir")
	}

	storage := filesystem.NewStorage(gitfs, cache.NewObjectLRUDefault())

	_, err = git.Clone(storage, urlfs, &git.CloneOptions{
		URL:   url.String(),
		Depth: depth,
	})

	if err == git.ErrRepositoryAlreadyExists {
		// gracefully handle the case where the repo already exists
		// intentionally empty block
		// probably could use a debug log
	} else if err != nil {
		// propagate other errors
		return errors.Wrapf(err, "failed to clone repo: %s", url)
	} else {
		logrus.Infof("[urls.gitcloner] repo %s successfully cloned", url)
	}

	return nil
}
