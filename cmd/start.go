package cmd

import (
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/mjpitz/gitfs/pkg/config"
	"github.com/mjpitz/gitfs/pkg/filesystem"
	"github.com/mjpitz/gitfs/pkg/tree"
	"github.com/mjpitz/gitfs/pkg/urls"
	rds "github.com/mjpitz/rds/pkg/config"
	"github.com/mjpitz/rds/pkg/remotes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// StartCommand defines the cobra.Command used to start the system.
var StartCommand = &cobra.Command{
	Use:   "start",
	Short: "Starts the file system server.",
	Run: func(cmd *cobra.Command, args []string) {
		current, err := user.Current()
		if err != nil {
			fail("failed to determine current user")
		}

		props := ConfigPath
		if len(props) == 0 {
			fail("missing config file")
		}

		cfg, err := config.Load(os.ExpandEnv(props))
		if err != nil {
			fail("failed to parse configuration: %v", err)
		}

		remote, err := remotes.ParseConfig(&rds.Configuration{
			Accounts: cfg.Accounts,
		})
		if err != nil {
			fail("failed to parse remote configuration: %v", err)
		}

		mountpoint := os.ExpandEnv(cfg.Mount)
		logrus.Infof("[main] configured mount point: %s", mountpoint)

		tree := gitfstree.NewTreeNode()

		logrus.Info("[main] fetching repositories")
		repositories, err := remote.ListRepositories()
		if err != nil {
			fail("failed to fetch repositories: %v", err)
		}

		logrus.Info("[main] parsing repositories into a directory structure")
		for _, repository := range repositories {
			url, err := urls.ParseURL(repository)
			if err != nil {
				logrus.Warnf("[main] failed to parse url: %v", err)
				continue
			}

			host := url.URL.Host
			fullPath := url.URL.Path

			// remove suffix
			ext := path.Ext(fullPath)
			fullPath = strings.TrimSuffix(fullPath, ext)
			fullPath = strings.TrimPrefix(fullPath, "/")

			parts := strings.Split(fullPath, "/")

			var fullPathParts []string
			fullPathParts = append(fullPathParts, host)
			fullPathParts = append(fullPathParts, parts...)

			tree.Insert(fullPathParts, url)
		}

		logrus.Infof("[main] attempting to mount %s", mountpoint)
		c, err := fuse.Mount(mountpoint)
		if err != nil {
			fail("failed to mount mountpoint: %v", err)
		}
		defer c.Close()

		uid, _ := strconv.Atoi(current.Uid)
		gid, _ := strconv.Atoi(current.Gid)

		cloner := urls.NewFileSystemAdapter(cfg.Clone)

		filesys := &filesystem.FileSystem{
			UID:  uint32(uid),
			GID:  uint32(gid),
			Tree: tree,
			FSA:  cloner,
		}

		logrus.Info("[main] now serving file system")
		if err := fs.Serve(c, filesys); err != nil {
			fail("failed to serve filesystem: %v", err)
		}

		<-c.Ready
		if err := c.MountError; err != nil {
			fail("file system failed to start: %v", err)
		}
	},
}
