package cmd

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/mjpitz/gitfs/pkg/clone"
	"github.com/mjpitz/gitfs/pkg/config"
	"github.com/mjpitz/gitfs/pkg/filesystem"
	"github.com/mjpitz/gitfs/pkg/remotes"
	"github.com/mjpitz/gitfs/pkg/tree"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"strconv"
	"strings"
)

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

		remote, err := remotes.ParseConfig(cfg)
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
			start := 4
			separator := strings.Index(repository, ":")
			end := len(repository) - 4

			// git@code.corp.indeed.com:squall/yellowhat-snapshots-prod.git
			// git@<<HOST>>:<<PATH>>.git
			host := repository[start:separator]
			path := repository[separator+1 : end]

			// expand into directory structure
			parts := strings.Split(path, "/")

			var fullPathParts []string
			fullPathParts = append(fullPathParts, host)
			fullPathParts = append(fullPathParts, parts...)

			tree.Insert(fullPathParts, repository)
		}

		logrus.Infof("[main] attempting to mount %s", mountpoint)
		c, err := fuse.Mount(mountpoint)
		if err != nil {
			fail("failed to mount mountpoint: %v", err)
		}
		defer c.Close()

		uid, _ := strconv.Atoi(current.Uid)
		gid, _ := strconv.Atoi(current.Gid)

		cloner := clone.NewCloner(cfg.Clone)

		filesys := &filesystem.FileSystem{
			Uid:    uint32(uid),
			Gid:    uint32(gid),
			Tree:   tree,
			Cloner: cloner,
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
