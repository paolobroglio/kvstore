package storage

import "sync"

type BTreeNode struct {
	keys []string
	locations []Location
	children []*BTreeNode
	isLeaf bool
	keyCount int
}

func newLeafNode(degree int) *BTreeNode {
	maxKeys := 2 * degree - 1
	return &BTreeNode{
		keys: make([]string, maxKeys),
		locations: make([]Location, maxKeys),
		children: make([]*BTreeNode, degree*2),
		isLeaf: true,
		keyCount: 0,
	}
}

func Search(node *BTreeNode, key string) *Location {
	i := 0
	for i < node.keyCount && key > node.keys[i] {
		i++
	}

	if i < node.keyCount && key == node.keys[i] {
		return &node.locations[i]
	}

	if node.isLeaf {
		return nil
	}

	return Search(node.children[i], key)
}

type BTreeIndex struct {
	root *BTreeNode
	degree int 
	mu sync.RWMutex
}

func NewBTreeIndex(degree int) *BTreeIndex {
	if degree < 2 {
		degree = 2
	}
	return &BTreeIndex{
		root: newLeafNode(degree),
		degree: degree,
	}
}

func (index *BTreeIndex) Get(key []byte) *Location {
	index.mu.RLock()
	defer index.mu.RUnlock()

	return Search(index.root, string(key))
}