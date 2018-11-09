package gitfstree

import "reflect"

func NewTreeNode() *TreeNode {
	return &TreeNode{
		children: make(map[string]*TreeNode),
	}
}

type TreeNode struct {
	URL      string
	children map[string]*TreeNode
}

func (t *TreeNode) Insert(path []string, url string) {
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

func (t *TreeNode) Children() []string {
	keys := reflect.ValueOf(t.children).MapKeys()
	stringKeys := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		stringKeys[i] = keys[i].String()
	}

	return stringKeys
}

func (t *TreeNode) Child(part string) (*TreeNode, bool) {
	node, ok := t.children[part]
	return node, ok
}
