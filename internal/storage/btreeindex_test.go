package storage

import "testing"

func TestBTreeSearchOnlyRoot(t *testing.T) {
	root := newLeafNode(2)
	root.keys[0] = "1"
	root.keys[1] = "2"
	root.locations[0] = Location{}
	root.locations[1] = Location{}
	root.keyCount = 2
	
	location1 := Search(root, "1")
	location2 := Search(root, "2")
	location3 := Search(root, "3")

	if location1 == nil {
		t.Errorf("Search(root, 1) returned nil")
	}
	if location2 == nil {
		t.Errorf("Search(root, 2) returned nil")
	}
	if location3 != nil {
		t.Errorf("Search(root, 3) returned a Location")
	}
}

func TestBTreeSearch(t *testing.T) {
	//       [2, ]
	//      /     \
	//    [1, ]   [3,4]  
	root := newLeafNode(2)
	root.keys[0] = "2"
	root.locations[0] = Location{}
	root.keyCount = 1
	root.isLeaf = false

	left := newLeafNode(2)
	left.keys[0] = "1"
	left.locations[0] = Location{}
	left.keyCount = 1
	
	right := newLeafNode(2)
	right.keys[0] = "3"
	right.keys[1] = "4"
	right.locations[0] = Location{}
	right.locations[1] = Location{}
	right.keyCount = 2

	root.children[0] = left
	root.children[1] = right
	
	location1 := Search(root, "1")
	location2 := Search(root, "2")
	location3 := Search(root, "3")
	location4 := Search(root, "4")

	if location1 == nil {
		t.Errorf("Search(root, 1) returned nil")
	}
	if location2 == nil {
		t.Errorf("Search(root, 2) returned nil")
	}
	if location3 == nil {
		t.Errorf("Search(root, 3) returned nil")
	}
	if location4 == nil {
		t.Errorf("Search(root, 4) returned nil")
	}
}

