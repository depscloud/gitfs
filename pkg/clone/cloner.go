package clone

import (
	"crypto/sha256"
	"github.com/mjpitz/gitfs/pkg/config"
	"github.com/mjpitz/gitfs/pkg/sync"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
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

func (c *Cloner) Resolve(url string) (string, int32, error) {
	cfg := c.cfg

	root := ""
	if cfg.RepositoryRoot != nil {
		root = cfg.RepositoryRoot.Value
	}

	depth := int32(1)
	if cfg.Depth != nil {
		depth = cfg.Depth.Value
	}

	for key, value := range cfg.Overrides {
		if key != url {
			regex, err := regexp.Compile(key)

			if err != nil {
				return "", 0, err
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

	return root, depth, nil
}

func (c *Cloner) fs(root, url string) billy.Filesystem {
	fs, ok := c.rootfs[root]

	if !ok {
		if len(root) == 0 {
			fs = &sync.Filesystem{
				Delegate: memfs.New(),
			}
		} else {
			fs = osfs.New(root)
		}
		c.rootfs[root] = fs
	}

	return fs
}

func (c *Cloner) Clone(url string) (billy.Filesystem, error) {
	root, depth, err := c.Resolve(url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve url: %s", url)
	}

	fs := c.fs(root, url)

	// use a sha256 in the url to insulate the paths
	urlPath := string(sha256.New().Sum([]byte(url)))

	// pull from a cache before chroot
	urlfs, ok := c.fscache[urlPath]
	if !ok {
		urlfs, err = fs.Chroot(urlPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create %s dir for %s", urlPath, url)
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

		c.fscache[urlPath] = urlfs
	}
	return urlfs, nil
}
