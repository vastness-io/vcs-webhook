package core

import (
	"sync"
)

//Node is a Branch within the BinarySearchTree.
type Node struct {
	value       int
	left, right *Node
}

// BinarySearchTree (BST) are a particular type of container: data structures that store integers.
type BinarySearchTree struct {
	lock sync.RWMutex

	Root *Node //Root of the tree
}

// NewBinarySearchTree creates a new BinarySearchTree pointer.
func NewBinarySearchTree() *BinarySearchTree {
	return &BinarySearchTree{}
}

// Insert adds a new node into the tree, position is determined by the int value passed.
func (t *BinarySearchTree) Insert(v int) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.Root == nil {
		t.Root = &Node{
			value: v,
		}
		return
	}
	_insert(t.Root, v)
}

// Find will search for the Node with the given int value passed.
// Returns Node and true if the Node exists otherwise Nil and false.
func (t *BinarySearchTree) Find(v int) (*Node, bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if t.Root == nil {
		return nil, false
	}

	return _find(t.Root, v)
}

func _find(node *Node, v int) (*Node, bool) {
	if node == nil {
		return nil, false
	} else if v == node.value {
		return node, true
	} else if v > node.value {
		return _find(node.right, v)
	}
	return _find(node.left, v)
}

func _insert(node *Node, v int) {
	if v > node.value {
		if node.right == nil {
			node.right = &Node{value: v}
			return
		}
		_insert(node.right, v)
	} else if v < node.value {
		if node.left == nil {
			node.left = &Node{value: v}
			return
		}
		_insert(node.left, v)
	}
}

// Delete removes the Node with the given int value from the tree.
// This is a no-op if the Node does not exist.
func (t *BinarySearchTree) Delete(v int) {
	t.Root = _delete(t.Root, v)
}

func _delete(node *Node, v int) *Node {
	if node == nil {
		return nil
	} else if v < node.value {
		node.left = _delete(node.left, v)
	} else if v > node.value {
		node.right = _delete(node.right, v)
	} else {
		if node.left == nil && node.right == nil {
			node = nil
		} else if node.left == nil {
			node = node.right
		} else if node.right == nil {
			node = node.left
		} else {
			maxNode := GetMaxInSubTree(node.left)
			node.value = maxNode.value
			node.left = _delete(node.left, maxNode.value)
		}
	}
	return node
}

// GetMaxInSubTree finds the greatest int within the Subtree
func GetMaxInSubTree(node *Node) *Node {
	if node.right == nil {
		return node
	}
	return GetMaxInSubTree(node.right)
}

// GetMinInSubTree finds the lowest int within the Subtree
func GetMinInSubTree(node *Node) *Node {
	if node.left == nil {
		return node
	}
	return GetMinInSubTree(node.left)
}
