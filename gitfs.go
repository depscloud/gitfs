package main

import (
	"flag"
	"fmt"
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
	rlog "github.com/sirupsen/logrus"
)

func usage() {
	fmt.Println("usage: gitfs")
	flag.PrintDefaults()
	os.Exit(2)
}

func fail(message string, data ...interface{}) {
	fmt.Println(fmt.Sprintf(message, data...))
	os.Exit(1)
}

func main() {
	current, err := user.Current()
	if err != nil {
		fail("failed to determine current user")
	}

	props := path.Join(current.HomeDir, ".gitfs.yml")
	flag.StringVar(&props, "config", props, "Specify the configuration path")
	flag.Parse()

	if len(props) == 0 {
		fail("missing config file")
	}

	config, err := config.Load(os.ExpandEnv(props))
	if err != nil {
		fail("failed to parse configuration: %v", err)
	}

	remote, err := remotes.ParseConfig(config)
	if err != nil {
		fail("failed to parse remote configuration: %v", err)
	}

	mountpoint := os.ExpandEnv(config.Mount)
	rlog.Infof("configured mount point: %s", mountpoint)

	tree := gitfstree.NewTreeNode()

	rlog.Info("fetching repositories")
	repositories, err := remote.ListRepositories()
	if err != nil {
		fail("failed to fetch repositories: %v", err)
	}

	rlog.Info("parsing repositories into a directory structure")
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

		//rlog.Infof("parts: [%s]", strings.Join(fullPathParts, ","))
		tree.Insert(fullPathParts, repository)
	}

	rlog.Infof("attempting to mount %s", mountpoint)
	c, err := fuse.Mount(mountpoint)
	if err != nil {
		fail("failed to mount mountpoint: %v", err)
	}
	defer c.Close()

	currentUser, err := user.Current()
	if err != nil {
		fail("failed to determine current user: %v", err)
	}

	uid, _ := strconv.Atoi(currentUser.Uid)
	gid, _ := strconv.Atoi(currentUser.Gid)

	filesys := &filesystem.FileSystem{
		Uid:  uint32(uid),
		Gid:  uint32(gid),
		Tree: tree,
	}

	rlog.Info("now serving file system")
	if err := fs.Serve(c, filesys); err != nil {
		fail("failed to serve filesystem: %v", err)
	}

	<-c.Ready
	if err := c.MountError; err != nil {
		fail("file system failed to start: %v", err)
	}
}
