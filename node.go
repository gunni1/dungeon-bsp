package main

import (
	"image"
	"image/color"
	"math/rand"
)

const MIN_NODE_SIZE = 5

//MIN_ROOM muss immer kleiner sein als MIN NODE!
const MIN_ROOM_SIZE = 3

type Node struct {
	left  *Node
	right *Node
	room  *Room
	//TODO: Path between left & right (corridor)
	x, y, width, height int
}

//A Rect inside of a node. Only present in Leaf Nodes
type Room struct {
	x, y, width, height int
}

//Create Room elements for all Leaf Nodes
func (node *Node) CreateLeafRooms(rnd rand.Rand) {
	if node.isLeaf() {
		node.GenerateRoom(rnd)
	} else {
		node.left.CreateLeafRooms(rnd)
		node.right.CreateLeafRooms(rnd)
	}
}

//Generate a random sized Room in the bounds of the Node
func (node *Node) GenerateRoom(rnd rand.Rand) {
	roomW := rnd.Intn(node.width-MIN_ROOM_SIZE) + MIN_ROOM_SIZE
	roomH := rnd.Intn(node.height-MIN_ROOM_SIZE) + MIN_ROOM_SIZE
	xOffset := rnd.Intn(node.width - roomW)
	yOffset := rnd.Intn(node.height - roomH)
	node.room = &Room{node.x + xOffset, node.y + yOffset, roomW, roomH}
}

func (node Node) isLeaf() bool {
	return node.left == nil || node.right == nil
}

func (pNode *Node) Split(rnd rand.Rand) {
	vert := ShouldForceVerticalSplit(*pNode, rnd)
	maxSplit := calcMaxSplit(vert, *pNode)

	//Dont Split if already to small
	if maxSplit <= MIN_NODE_SIZE {
		return
	}
	split := rnd.Intn(maxSplit-MIN_NODE_SIZE) + MIN_NODE_SIZE
	if vert {
		//Vertikal split |
		pNode.left = &Node{x: pNode.x, y: pNode.y, width: split, height: pNode.height}
		pNode.right = &Node{x: pNode.x + split, y: pNode.y, width: pNode.width - split, height: pNode.height}

	} else {
		//Horizontal split ---
		pNode.left = &Node{x: pNode.x, y: pNode.y, width: pNode.width, height: split}
		pNode.right = &Node{x: pNode.x, y: pNode.y + split, width: pNode.width, height: pNode.height - split}
	}
}

func calcMaxSplit(splitVertical bool, node Node) int {
	if splitVertical {
		return node.width - MIN_NODE_SIZE
	} else {
		return node.height - MIN_NODE_SIZE
	}
}

func ShouldForceVerticalSplit(node Node, rnd rand.Rand) bool {
	if node.width > node.height && float64(node.width)/float64(node.height) > 1.5 {
		return true
	} else if node.height > node.width && float64(node.height)/float64(node.width) > 1.5 {
		return false
	} else {
		return rnd.Intn(2) == 1
	}
}

func (node *Node) SplitDeep(rnd rand.Rand, depth int) {
	if depth > 0 {
		node.Split(rnd)
		//FIXME: Nicht immer splitten für mehr Variabilität?
		if node.left != nil && node.right != nil {
			node.left.SplitDeep(rnd, depth-1)
			node.right.SplitDeep(rnd, depth-1)
		}
	}
}

func (node Node) RenderRooms() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, node.width, node.height))

	leafs := node.CollectLeafs()
	for _, leaf := range leafs {
		if leaf.room != nil {
			paintFilled(img, *leaf.room)
		}
	}

	return img
}

func (node Node) Render() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, node.width, node.height))

	leafs := node.CollectLeafs()

	for _, leaf := range leafs {
		outline(img, leaf)
	}
	return img
}

func paintFilled(img *image.RGBA, room Room) {
	gray := color.RGBA{211, 211, 211, 255}
	//roomRect := image.Rect(0, 0, room.width, room.height)
	//draw.Draw(img, roomRect, &image.Uniform{gray}, image.ZP, draw.Src)
	for i := room.x; i < room.x+room.width; i++ {
		for j := room.y; j < room.y+room.height; j++ {
			img.Set(i, j, gray)
		}
	}
}

func outline(img *image.RGBA, node Node) {
	green := color.RGBA{0, 100, 0, 255}
	for i := node.x; i < node.x+node.width; i++ {
		img.Set(i, node.y, green)
		img.Set(i, node.y+node.height-1, green)
	}
	for j := node.y; j < node.y+node.height; j++ {
		img.Set(node.x, j, green)
		img.Set(node.x+node.width-1, j, green)
	}
}

//Collect Nodes without further child nodes. Recursion might be implemented in a better way...
func (node Node) CollectLeafs() []Node {
	if node.isLeaf() {
		return []Node{node}
	} else {
		return append(node.left.CollectLeafs(), node.right.CollectLeafs()...)
	}
}
