package main

import "fmt"

// IndexTree represents a B-tree
type IndexTree struct {
	Root      *IndexTreeNode
	MinDegree int // Minimum degree of the B-tree.
}

type IndexValue struct {
	Key uint32
	Pos int64
}

// IndexTreeNode represents a node in the B-tree
type IndexTreeNode struct {
	IsLeaf bool
	Keys   []IndexValue
	Child  []*IndexTreeNode
}

// searchIndexValue returns the index of the first Key that is equal or greater than k in the Keys array.
// If all Keys are less than k, then it returns len(b.Keys)
func (b *IndexTreeNode) searchIndexValue(k uint32) int {
	idx := 0
	for idx < len(b.Keys) && b.Keys[idx].Key < k {
		idx++
	}
	return idx
}

func (idx *IndexTree) Search(k uint32) (*IndexTreeNode, int) {
	return idx.Root.Search(k)
}

func (idx *IndexTree) Get(k uint32) (int64, bool) {
	node, index := idx.Search(k)
	if node == nil {
		return -1, false
	}

	return node.Keys[index].Pos, true
}

// Search returns the node containing the Key and the index of the Key in the Keys array.
func (b *IndexTreeNode) Search(k uint32) (*IndexTreeNode, int) {
	idx := b.searchIndexValue(k)
	// if the Key is found in this node, return this node and the index of the Key.
	if idx < len(b.Keys) && b.Keys[idx].Key == k {
		return b, idx
	} else if b.IsLeaf {
		// if the node is a leaf node and the Key is not in this node, return nil.
		return nil, -1
	} else {
		// if the node is not a leaf, search the appropriate Child node.
		return b.Child[idx].Search(k)
	}
}

// Insert inserts a Key into the tree.
func (t *IndexTree) Insert(Key *IndexValue) {
	Root := t.Root
	if len(Root.Keys) == (2*t.MinDegree)-1 {
		// Create a new Root because the Root is full
		temp := &IndexTreeNode{}
		t.Root = temp
		temp.Child = append(temp.Child, Root)
		// Split the old Root and move 1 Key to the new Root
		t.splitChild(temp, 0)
		t.insertNonFull(temp, Key)
	} else {
		// If the Root is not full, call insertNonFull for the Root
		t.insertNonFull(Root, Key)
	}
}

// insertNonFull inserts a Key into a non-full node.
func (t *IndexTree) insertNonFull(x *IndexTreeNode, Key *IndexValue) {
	// Start from the rightmost Key in the node
	i := len(x.Keys) - 1

	// If the node is a leaf node
	if x.IsLeaf {
		// Append a new Key at the end of the Keys
		x.Keys = append(x.Keys, IndexValue{Key: 0, Pos: 0})

		// Shift all Keys greater than k to the right
		for i >= 0 && Key.Key < x.Keys[i].Key {
			x.Keys[i+1] = x.Keys[i]
			i--
		}

		// Insert the new Key at the found Position
		x.Keys[i+1] = *Key
	} else {
		// If the node is not a leaf, find the Child which is going to hold the new Key
		for i >= 0 && Key.Key < x.Keys[i].Key {
			i--
		}
		i++

		// If the found Child is full
		if len(x.Child[i].Keys) == (2*t.MinDegree)-1 {
			// Split the Child
			t.splitChild(x, i)

			// After split, the middle Key of the Child moves up and the Child is split into two.
			// Check which of the two Children is going to hold the new Key
			if Key.Key > x.Keys[i].Key {
				i++
			}
		}

		// Insert the Key into the Child
		t.insertNonFull(x.Child[i], Key)
	}
}

// splitChild splits the Child y of x.
func (t *IndexTree) splitChild(x *IndexTreeNode, i int) {
	// tt is the minimum degree of the B-tree
	tt := t.MinDegree

	// y is the i-th Child of x that is going to be split
	y := x.Child[i]

	// z is the new node, created to store tt-1 Keys of y
	z := &IndexTreeNode{IsLeaf: y.IsLeaf}

	// Make space for the new Child
	x.Child = append(x.Child, nil)
	copy(x.Child[i+2:], x.Child[i+1:])
	x.Child[i+1] = z

	// Make space for the new Key in x
	x.Keys = append(x.Keys, IndexValue{Key: 0, Pos: 0})
	copy(x.Keys[i+1:], x.Keys[i:])

	// Move the middle Key of y to x
	x.Keys[i] = y.Keys[tt-1]

	// Move the last tt-1 Keys of y to z
	z.Keys = append(z.Keys, y.Keys[tt:]...)
	y.Keys = y.Keys[:tt-1]

	// If y is not a leaf, move the last tt Children of y to z
	if !y.IsLeaf {
		z.Child = append(z.Child, y.Child[tt:]...)
		y.Child = y.Child[:tt]
	}
}

// Print prints the contents of the IndexTree.
func (t *IndexTree) Print() {
	t.Root.Print(0)
}

// Print prints the contents of the IndexTreeNode.
func (n *IndexTreeNode) Print(level int) {
	// Print the Keys in this node
	fmt.Printf("Level %d: ", level)
	for _, Key := range n.Keys {
		fmt.Printf("%d ", Key.Key)
	}
	fmt.Println()

	// If this node is not a leaf, print its Child nodes
	if !n.IsLeaf {
		for _, Child := range n.Child {
			Child.Print(level + 1)
		}
	}
}

func createIndexTree(Keys []IndexValue, t int) *IndexTree {
	tree := &IndexTree{
		Root: &IndexTreeNode{
			IsLeaf: true,
			Keys:   make([]IndexValue, 0),
			Child:  make([]*IndexTreeNode, 0),
		},
		MinDegree: t,
	}

	for _, Key := range Keys {
		tree.Insert(&Key)
	}

	return tree
}
