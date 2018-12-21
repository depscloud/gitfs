package main

import (
	"flag"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/mjpitz/gitfs/pkg/config"
	"github.com/mjpitz/gitfs/pkg/filesystem"
	"github.com/mjpitz/gitfs/pkg/remotes"
	"github.com/mjpitz/gitfs/pkg/tree"
	"github.com/sirupsen/logrus"
)

func fail(message string, data ...interface{}) {
	logrus.Errorf("[main] " + message, data...)
	os.Exit(1)
}

func main() {
	current, err := user.Current()
	if err != nil {
		fail("failed to determine current user")
	}

	props := path.Join(current.HomeDir, ".gitfs/config.yml")
	flag.StringVar(&props, "config", props, "Specify the configuration path")
	flag.Parse()

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

	filesys := &filesystem.FileSystem{
		Uid:  uint32(uid),
		Gid:  uint32(gid),
		Tree: tree,
	}

	logrus.Info("[main] now serving file system")
	if err := fs.Serve(c, filesys); err != nil {
		fail("failed to serve filesystem: %v", err)
	}

	<-c.Ready
	if err := c.MountError; err != nil {
		fail("file system failed to start: %v", err)
	}
}
