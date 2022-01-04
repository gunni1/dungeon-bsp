package main

import (
	"math/rand"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestCollectLeafsSimple(t *testing.T) {
	l := Node{x: 0, y: 0, width: 4, height: 10}
	r := Node{x: 4, y: 0, width: 6, height: 10}
	root := Node{x: 0, y: 0, width: 10, height: 10, left: &l, right: &r}

	result := root.CollectLeafs()

	Equal(t, len(result), 2)
	Contains(t, result, l)
	Contains(t, result, r)
}

func TestRollDirectionVert(t *testing.T) {
	rnd := rand.New(rand.NewSource(7))
	node := Node{x: 0, y: 0, width: 20, height: 8}
	result := RollDirection(node, *rnd)
	Equal(t, result, true)

}

func TestRollDirectionHori(t *testing.T) {
	rnd := rand.New(rand.NewSource(7))
	node := Node{x: 0, y: 0, width: 8, height: 20}
	result := RollDirection(node, *rnd)
	Equal(t, result, true)
}

func TestSplit(t *testing.T) {
	rnd := rand.New(rand.NewSource(7))
	root := Node{x: 0, y: 0, width: 20, height: 20}
	root.Split(*rnd)

	Equal(t, root.left, &Node{x: 0, y: 0, width: 20, height: 5})
	Equal(t, root.right, &Node{x: 0, y: 5, width: 20, height: 15})
}
