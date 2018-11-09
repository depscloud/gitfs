package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"flag"
	"fmt"
	"github.com/indeedeng/gitfs/lib/filesystem"
	"github.com/indeedeng/gitfs/lib/remotes"
	"github.com/indeedeng/gitfs/lib/tree"
	"indeed/gophers/rlog"
	"os"
	"strings"
)

func usage() {
	rlog.Warn("usage: gitfs <mountpoint>")
	flag.PrintDefaults()
	os.Exit(2)
}

func fail(message string, data... interface{}) {
	fmt.Fprintf(os.Stderr, message, data...)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
	}

	mountpoint := flag.Args()[0]
	rlog.Infof("configured mount point: %s", mountpoint)

	darwin := remotes.NewDarwinRemote()
	tree := gitfstree.NewTreeNode()

	rlog.Info("fetching repositories from darwin")
	repositories, err := darwin.ListRepositories()
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
		path := repository[separator + 1:end]

		// expand into directory structure
		parts := strings.Split(path, "/")

		var fullPathParts []string
		fullPathParts = append(fullPathParts, host)
		fullPathParts = append(fullPathParts, parts...)

		rlog.Infof("parts: [%s]", strings.Join(fullPathParts, ","))
		tree.Insert(fullPathParts, repository)
	}

	rlog.Infof("attempting to mount %s", mountpoint)
	c, err := fuse.Mount(mountpoint)
	if err != nil {
		fail("failed to mount mountpoint: %v", err)
	}
	defer c.Close()

	filesys := &filesystem.FileSystem{
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
