package dbsp

import (
	"math/rand"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestCollectLeafsSimple(t *testing.T) {
	l := Node{X: 0, Y: 0, Width: 4, Height: 10}
	r := Node{X: 4, Y: 0, Width: 6, Height: 10}
	root := Node{X: 0, Y: 0, Width: 10, Height: 10, Left: &l, Right: &r}

	result := root.CollectLeafs()

	Equal(t, len(result), 2)
	Contains(t, result, l)
	Contains(t, result, r)
}

func TestShouldForceVerticalSplitert(t *testing.T) {
	rnd := rand.New(rand.NewSource(7))
	node := Node{X: 0, Y: 0, Width: 20, Height: 8}
	result := ShouldForceVerticalSplit(node, *rnd)
	Equal(t, result, true)

}

func TestForceVerticalSplitHori(t *testing.T) {
	rnd := rand.New(rand.NewSource(7))
	node := Node{X: 0, Y: 0, Width: 8, Height: 20}
	result := ShouldForceVerticalSplit(node, *rnd)
	Equal(t, result, false)
}

func TestSplit(t *testing.T) {
	rnd := rand.New(rand.NewSource(7))
	root := Node{X: 0, Y: 0, Width: 20, Height: 20}
	root.Split(*rnd)

	Equal(t, root.Left, &Node{X: 0, Y: 0, Width: 20, Height: 5})
	Equal(t, root.Right, &Node{X: 0, Y: 5, Width: 20, Height: 15})
}
