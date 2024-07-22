package dbsp

import (
	"math/rand"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestGenerateRoomAlwaysInBound(t *testing.T) {
	rnd := rand.New(rand.NewSource(7))

	for i := 0; i < 10; i++ {
		node := createRandomNode(*rnd)
		node.GenerateRoom(*rnd)
		isRoomInBounds(t, node)
	}
}

func createRandomNode(rnd rand.Rand) Node {
	min := 10
	max := 40
	width := rnd.Intn(max-min) + min
	height := rand.Intn(max-min) + min
	return Node{X: rnd.Intn(max), Y: rnd.Intn(max), Width: width, Height: height}
}

func isRoomInBounds(t *testing.T, node Node) {
	room := node.Room
	NotNil(t, room)
	GreaterOrEqual(t, room.x, node.X)
	GreaterOrEqual(t, room.y, node.Y)

	LessOrEqual(t, room.x+room.width, node.X+node.Width)
	LessOrEqual(t, room.y+room.height, node.Y+node.Height)
}

func TestCollectRoomParentNodes(t *testing.T) {
	root := createTree()
	results := root.CollectRoomParentNodes()

	Equal(t, 3, len(results))
	//Expect Nodes with specific X values as identifiers
	//See createTree for details
	True(t, containsNodeWithX(results, 10))
	True(t, containsNodeWithX(results, 30))
	True(t, containsNodeWithX(results, 40))
}

func TestCollectRoomParentNodesDontBreakWithoutRooms(t *testing.T) {
	n2 := Node{X: 2}
	n3 := Node{X: 3}
	root := Node{X: 1, Left: &n2, Right: &n3}
	result := root.CollectRoomParentNodes()
	Empty(t, result)
}

func containsNodeWithX(nodes []Node, x int) bool {
	for _, node := range nodes {
		if node.X == x {
			return true
		}
	}
	return false
}

//			         100
//		  10				 20
//	  11     12    	   30  	      40
//				    31    32    41    42
func createTree() Node {
	n11 := Node{X: 11, Room: &Room{}}
	n12 := Node{X: 12, Room: &Room{}}
	n10 := Node{X: 10, Left: &n11, Right: &n12}
	n31 := Node{X: 31, Room: &Room{}}
	n32 := Node{X: 32, Room: &Room{}}
	n41 := Node{X: 41, Room: &Room{}}
	n42 := Node{X: 42, Room: &Room{}}
	n30 := Node{X: 30, Left: &n31, Right: &n32}
	n40 := Node{X: 40, Left: &n41, Right: &n42}
	n20 := Node{X: 20, Left: &n30, Right: &n40}
	n100 := Node{X: 100, Left: &n10, Right: &n20}
	return n100
}
