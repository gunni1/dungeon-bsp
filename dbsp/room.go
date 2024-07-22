package dbsp

import (
	"image"
	"math/rand"
)

// Create Room elements for all Leaf Nodes
func (node *Node) CreateLeafRooms(rnd rand.Rand) {
	if node.isLeaf() {
		node.GenerateRoom(rnd)
	} else {
		node.Left.CreateLeafRooms(rnd)
		node.Right.CreateLeafRooms(rnd)
	}
}

// Generate a random sized Room in the bounds of the Node
func (node *Node) GenerateRoom(rnd rand.Rand) {
	roomW := rnd.Intn(node.Width-MIN_ROOM_SIZE) + MIN_ROOM_SIZE
	roomH := rnd.Intn(node.Height-MIN_ROOM_SIZE) + MIN_ROOM_SIZE
	xOffset := rnd.Intn(node.Width - roomW)
	yOffset := rnd.Intn(node.Height - roomH)
	node.Room = &Room{node.X + xOffset, node.Y + yOffset, roomW, roomH}
}

func (node Node) RenderRooms() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, node.Width, node.Height))

	leafs := node.CollectLeafs()
	for _, leaf := range leafs {
		if leaf.Room != nil {
			paintFilled(img, *leaf.Room)
		}
	}
	return img
}
