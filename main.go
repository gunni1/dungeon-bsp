package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

type Node struct {
	left                *Node
	right               *Node
	x, y, width, height int
}

func (node Node) isLeaf() bool {
	return node.left == nil || node.right == nil
}

func (pNode *Node) Split(rnd rand.Rand) {
	if pNode.left != nil || pNode.right != nil {
		panic("node already splitted")
	}
	var newLeft, newRight Node
	if rnd.Intn(2) == 1 {
		//Horizontal -> x
		splitMarker := rnd.Intn(pNode.x + pNode.width)
		newLeft = Node{x: pNode.x, y: pNode.y, width: splitMarker - pNode.x, height: pNode.height}
		newRight = Node{x: splitMarker, y: pNode.y, width: pNode.width - newLeft.width, height: pNode.height}

	} else {
		//Vertical -> y
		splitMarker := rnd.Intn(pNode.y + pNode.height)
		newLeft = Node{x: pNode.x, y: pNode.y, width: pNode.width, height: splitMarker - pNode.y}
		newRight = Node{x: pNode.x, y: pNode.y + splitMarker, width: pNode.width, height: pNode.height - newLeft.height}
	}
	pNode.left = &newLeft
	pNode.right = &newRight

}

func (node *Node) SplitDeep(rnd rand.Rand, depth int) {
	if depth > 0 {
		node.Split(rnd)
		node.left.SplitDeep(rnd, depth-1)
		node.right.SplitDeep(rnd, depth-1)
	}

}

func (node Node) Render() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, node.width, node.height))

	leafs := node.CollectLeafs()

	for _, leaf := range leafs {
		outline(img, leaf)
	}

	//Collect all leafs
	//Render Leafs as Rect
	return img
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

func (node Node) CollectLeafs() []Node {

	if node.isLeaf() {
		return []Node{node}
	} else {
		return append(node.left.CollectLeafs(), node.right.CollectLeafs()...)
	}
}

func main() {
	rndSource := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(rndSource)
	root := Node{x: 0, y: 0, width: 100, height: 100}
	root.SplitDeep(*rnd, 2)
	//root.Split(*rnd)
	//root.right.Split(*rnd)

	img := root.Render()
	file, _ := os.Create("out.png")
	png.Encode(file, img)
}
