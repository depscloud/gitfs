package filesystem

import "bazil.org/fuse/fs"

// INode encapsulates all the functionality that should be implemented to
// emulate an inode in unix file systems. I managed to stumble across this
// in a cockroachdb project, but can no longer find the reference.
type INode interface {
	// node functions
	fs.Node
	fs.NodeSetattrer

	// directory functions
	fs.NodeStringLookuper
	fs.HandleReadDirAller
	fs.NodeMkdirer
	fs.NodeCreater
	fs.NodeRemover
	fs.NodeRenamer
	fs.NodeSymlinker

	// handle functions
	fs.NodeOpener
	fs.HandleWriter
	fs.HandleReader
	fs.NodeFsyncer
	fs.HandleFlusher
	fs.HandleReleaser

	// symlink functions
	fs.NodeReadlinker
}
