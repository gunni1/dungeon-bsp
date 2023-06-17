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
	min := 5
	max := 20
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
