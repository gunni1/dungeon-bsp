package main

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
	return Node{x: rnd.Intn(max), y: rnd.Intn(max), width: width, height: height}
}

func isRoomInBounds(t *testing.T, node Node) {
	room := node.room
	NotNil(t, room)
	GreaterOrEqual(t, room.x, node.x)
	GreaterOrEqual(t, room.y, node.y)

	LessOrEqual(t, room.x+room.width, node.x+node.width)
	LessOrEqual(t, room.y+room.height, node.y+node.height)
}
