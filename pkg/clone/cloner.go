package clone

import (
	"crypto/sha256"
	"encoding/base32"
	"github.com/mjpitz/gitfs/pkg/config"
	"github.com/mjpitz/gitfs/pkg/sync"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"os"
	"regexp"
)

func NewCloner(cfg *config.CloneConfiguration) *Cloner {
	rootfs := make(map[string]billy.Filesystem)
	fscache := make(map[string]billy.Filesystem)

	return &Cloner{
		cfg:     cfg,
		rootfs:  rootfs,
		fscache: fscache,
	}
}

type Cloner struct {
	cfg     *config.CloneConfiguration
	rootfs  map[string]billy.Filesystem
	fscache map[string]billy.Filesystem
}

func (c *Cloner) Resolve(url string) (string, string, int32, error) {
	cfg := c.cfg
	root := ""
	depth := int32(1)

	hash := sha256.New()
	hash.Write([]byte(url))
	bucket := base32.HexEncoding.
		WithPadding(base32.NoPadding).
		EncodeToString(hash.Sum(nil))

	if cfg == nil {
		return root, bucket, depth, nil
	}

	if cfg.RepositoryRoot != nil {
		root = cfg.RepositoryRoot.Value
	}

	if cfg.Depth != nil {
		depth = cfg.Depth.Value
	}

	for key, value := range cfg.Overrides {
		if key != url {
			regex, err := regexp.Compile(key)

			if err != nil {
				return "", "", 0, err
			}

			if !regex.Match([]byte(url)) {
				continue
			}
		}

		if value.RepositoryRoot != nil {
			root = value.RepositoryRoot.Value
		}

		if value.Depth != nil {
			depth = value.Depth.Value
		}

		break
	}

	return root, bucket, depth, nil
}

func (c *Cloner) fs(root string) billy.Filesystem {
	fs, ok := c.rootfs[root]

	if !ok {
		if len(root) == 0 {
			fs = &sync.Filesystem{
				Delegate: memfs.New(),
			}
		} else {
			fs = osfs.New(os.ExpandEnv(root))
		}
		c.rootfs[root] = fs
	}

	return fs
}

func (c *Cloner) Clone(url string) (billy.Filesystem, error) {
	root, bucket, depth, err := c.Resolve(url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve url: %s", url)
	}

	// pull from a cache before chroot
	urlfs, ok := c.fscache[bucket]
	if !ok {
		urlfs, err = c.fs(root).Chroot(bucket)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create %s dir for %s", bucket, url)
		}

		gitfs, err := urlfs.Chroot(git.GitDirName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create .git dir")
		}

		storage := filesystem.NewStorage(gitfs, cache.NewObjectLRUDefault())

		_, err = git.Clone(storage, urlfs, &git.CloneOptions{
			URL:   url,
			Depth: int(depth),
		})

		if err == git.ErrRepositoryAlreadyExists {
			// gracefully handle the case where the repo already exists
			// intentionally empty block
			// probably could use a debug log
		} else if err != nil {
			// propagate other errors
			return nil, errors.Wrapf(err, "failed to clone repo: %s", url)
		} else {
			logrus.Infof("[clone.cloner] repo %s successfully cloned", url)
		}

		c.fscache[bucket] = urlfs
	}
	return urlfs, nil
}
