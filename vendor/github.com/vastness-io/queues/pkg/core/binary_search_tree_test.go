package core

import (
	"math/rand"
	"testing"
)

func TestInsert(t *testing.T) {
	seedValue := 5000
	tree := NewBinarySearchTree()
	tree.Insert(seedValue) //seed value

	for i := 1; i <= seedValue; i++ {
		r := rand.Intn(seedValue)
		tree.Insert(r)
		if _, exists := tree.Find(r); !exists {
			t.Fail()
		}
	}

}

func TestFind(t *testing.T) {
	seedValue := 5000
	tree := NewBinarySearchTree()

	if _, exists := tree.Find(seedValue); exists {
		t.Fail()
	}

	tree.Insert(seedValue) //seed value

	for i := 1; i <= seedValue; i++ {
		r := rand.Intn(seedValue)
		tree.Insert(r)

		if _, exists := tree.Find(r); !exists {
			t.Errorf("%v should exist", r)
		}
	}

	// edge cases
	if _, exists := tree.Find(seedValue + 1); exists {
		t.Fail()
	}
	if _, exists := tree.Find(seedValue - 1); exists {
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	seedValue := 5000
	tree := NewBinarySearchTree()

	tree.Delete(rand.Int()) //can be any value

	if tree.Root != nil {
		t.Fail()
	}

	tree.Insert(rand.Intn(seedValue))

	for i := 1; i <= seedValue; i++ {
		r := rand.Intn(seedValue)
		tree.Insert(r)
		tree.Delete(r)
		if _, exists := tree.Find(r); exists {
			t.Errorf("%v shouldn't exist", r)
		}
	}
}

func TestGetMaxInSubTree(t *testing.T) {
	node := &Node{}
	if GetMaxInSubTree(node) != node {
		t.Fail()
	}

	node.right = &Node{}

	if GetMaxInSubTree(node) != node.right {
		t.Fail()
	}

}
