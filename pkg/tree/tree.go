package gitfstree

import (
	"reflect"

	"github.com/mjpitz/gitfs/pkg/urls"
)

// NewTreeNode constructs a node that is used during construction of the tree.
func NewTreeNode() *TreeNode {
	return &TreeNode{
		children: make(map[string]*TreeNode),
	}
}

// TreeNode represents a single node in a tree.
type TreeNode struct {
	URL      *urls.URL
	children map[string]*TreeNode
}

// Insert inserts a url element into the tree.
func (t *TreeNode) Insert(path []string, url *urls.URL) {
	ptr := t

	for _, part := range path {
		_, ok := ptr.children[part]

		if !ok {
			ptr.children[part] = NewTreeNode()
		}

		ptr = ptr.children[part]
	}

	ptr.URL = url
}

// Children returns the child keys of the tree.
func (t *TreeNode) Children() []string {
	keys := reflect.ValueOf(t.children).MapKeys()
	stringKeys := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		stringKeys[i] = keys[i].String()
	}

	return stringKeys
}

// Child retrieves a child at a given key.
func (t *TreeNode) Child(part string) (*TreeNode, bool) {
	node, ok := t.children[part]
	return node, ok
}
