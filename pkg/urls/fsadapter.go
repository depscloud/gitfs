package urls

import (
	"crypto/sha256"
	"encoding/base32"
	"os"
	"regexp"

	"github.com/deps-cloud/gitfs/pkg/config"
	"github.com/deps-cloud/gitfs/pkg/sync"

	"github.com/pkg/errors"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/osfs"
)

// NewFileSystemAdapter create a new FileSystemAdapter from the provided configuration.
// This method encapsulates creation of some caches.
func NewFileSystemAdapter(cfg *config.CloneConfiguration) *FileSystemAdapter {
	rootfs := make(map[string]billy.Filesystem)
	fscache := make(map[string]billy.Filesystem)
	cloner := NewCloner()

	return &FileSystemAdapter{
		cfg:     cfg,
		rootfs:  rootfs,
		fscache: fscache,
		cloner:  cloner,
	}
}

// FileSystemAdapter maintains two caches.
// One is a cache for the root filesystems.
// The other was intended to be a cache that was periodically pruned.
type FileSystemAdapter struct {
	cfg     *config.CloneConfiguration
	rootfs  map[string]billy.Filesystem
	fscache map[string]billy.Filesystem
	cloner  Cloner
}

// Resolve determines which bucket the url falls into.
// This function contains business logic around the configuration and is exposed for unit testing purposes.
func (fsa *FileSystemAdapter) Resolve(url string) (string, string, int32, error) {
	cfg := fsa.cfg
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

func (fsa *FileSystemAdapter) fs(root string) billy.Filesystem {
	fs, ok := fsa.rootfs[root]

	if !ok {
		if len(root) == 0 {
			fs = &sync.Filesystem{
				Delegate: memfs.New(),
			}
		} else {
			fs = osfs.New(os.ExpandEnv(root))
		}
		fsa.rootfs[root] = fs
	}

	return fs
}

// Clone accepts a url and clones it to an underlying filesystem.
func (fsa *FileSystemAdapter) Clone(url *URL) (billy.Filesystem, error) {
	rawurl := url.String()
	root, bucket, depth, err := fsa.Resolve(rawurl)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve url: %s", url)
	}

	// pull from a cache before chroot
	urlfs, ok := fsa.fscache[bucket]
	if !ok {
		urlfs, err = fsa.fs(root).Chroot(bucket)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create %s dir for %s", bucket, url)
		}

		if err = fsa.cloner.Clone(url, int(depth), urlfs); err != nil {
			return nil, errors.Wrapf(err, "failed to clone: %s url", rawurl)
		}

		fsa.fscache[bucket] = urlfs
	}

	// do something with urlfs
	return urlfs, nil
}
