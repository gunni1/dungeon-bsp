package dbsp

import (
	"image"
	"image/color"
	"image/color/palette"
	"math/rand"
)

// MIN_ROOM must be smaller than MIN NODE!
const MIN_NODE_SIZE = 10
const MIN_ROOM_SIZE = 5

type Node struct {
	Left                *Node
	Right               *Node
	Room                *Room
	corridor            *Room
	X, Y, Width, Height int
}

// A Rect inside of a node. Only present in Leaf Nodes
type Room struct {
	x, y, width, height int
}

// Context to protocol interim results
type ProtocolCtx struct {
	InterimResults chan Node
	RootNode       *Node
}

// Push the current binary tree to the interim results channel
func (prtcCtx ProtocolCtx) takeSnapshot() {
	prtcCtx.InterimResults <- *prtcCtx.RootNode
}

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

func (node *Node) ConnectLeafs(rnd rand.Rand) {
	if !node.isLeaf() {
		if node.Left.isLeaf() && node.Right.isLeaf() {
			//Roll place for path, mind vertial or horizontal
			//Create room
		} else {
			node.Left.ConnectLeafs(rnd)
			node.Right.ConnectLeafs(rnd)
		}
	}
}

func (node Node) IsSplitHorizontal() bool {
	if node.isLeaf() {
		panic("is not split. ")
	}
	return node.Left.X == node.Right.X
}

func (node Node) isLeaf() bool {
	return node.Left == nil || node.Right == nil
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
		//Vertical split |
		pNode.Left = &Node{X: pNode.X, Y: pNode.Y, Width: split, Height: pNode.Height}
		pNode.Right = &Node{X: pNode.X + split, Y: pNode.Y, Width: pNode.Width - split, Height: pNode.Height}

	} else {
		//Horizontal split ---
		pNode.Left = &Node{X: pNode.X, Y: pNode.Y, Width: pNode.Width, Height: split}
		pNode.Right = &Node{X: pNode.X, Y: pNode.Y + split, Width: pNode.Width, Height: pNode.Height - split}
	}
}

func calcMaxSplit(splitVertical bool, node Node) int {
	if splitVertical {
		return node.Width - MIN_NODE_SIZE
	} else {
		return node.Height - MIN_NODE_SIZE
	}
}

// Determine split direction based on nodes ratio.
// If ratio is not significant, toss a coin.
func ShouldForceVerticalSplit(node Node, rnd rand.Rand) bool {
	if node.Width > node.Height && float64(node.Width)/float64(node.Height) > 1.5 {
		return true
	} else if node.Height > node.Width && float64(node.Height)/float64(node.Width) > 1.5 {
		return false
	} else {
		return rnd.Intn(2) == 1
	}
}

func (node *Node) SplitDeep(rnd rand.Rand, depth int, prtcCtx ProtocolCtx) {
	prtcCtx.takeSnapshot()
	if depth > 0 && !skipSplitForVariety(rnd) {
		node.Split(rnd)
		if node.Left != nil && node.Right != nil {
			node.Left.SplitDeep(rnd, depth-1, prtcCtx)
			node.Right.SplitDeep(rnd, depth-1, prtcCtx)
		}
	}
}

// Return true in 1 out of 10 random cases
func skipSplitForVariety(rnd rand.Rand) bool {
	return rnd.Intn(10) == 1
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

func (node Node) RenderNode() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, node.Width, node.Height))
	leafs := node.CollectLeafs()
	for _, leaf := range leafs {
		outline(img, leaf)
	}
	return img
}

func (node Node) RenderNodePaletted() image.Image {

	img := image.NewPaletted(image.Rect(0, 0, node.Width, node.Height), palette.WebSafe)
	leafs := node.CollectLeafs()
	for _, leaf := range leafs {
		outlinePaletted(img, leaf)
	}
	return img
}

func paintFilled(img *image.RGBA, room Room) {
	gray := color.RGBA{211, 211, 211, 255}
	for i := room.x; i < room.x+room.width; i++ {
		for j := room.y; j < room.y+room.height; j++ {
			img.Set(i, j, gray)
		}
	}
}

func outlinePaletted(img *image.Paletted, node Node) {
	green := color.RGBA{0, 100, 0, 255}
	for i := node.X; i < node.X+node.Width; i++ {
		img.Set(i, node.Y, green)
		img.Set(i, node.Y+node.Height-1, green)
	}
	for j := node.Y; j < node.Y+node.Height; j++ {
		img.Set(node.X, j, green)
		img.Set(node.X+node.Width-1, j, green)
	}
}

func outline(img *image.RGBA, node Node) {
	green := color.RGBA{0, 100, 0, 255}
	for i := node.X; i < node.X+node.Width; i++ {
		img.Set(i, node.Y, green)
		img.Set(i, node.Y+node.Height-1, green)
	}
	for j := node.Y; j < node.Y+node.Height; j++ {
		img.Set(node.X, j, green)
		img.Set(node.X+node.Width-1, j, green)
	}
}

// Collect Nodes without further child nodes. Recursion might be implemented in a better way...
func (node Node) CollectLeafs() []Node {
	if node.isLeaf() {
		return []Node{node}
	} else {
		return append(node.Left.CollectLeafs(), node.Right.CollectLeafs()...)
	}
}
