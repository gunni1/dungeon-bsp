package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

const MIN_SIZE = 5

type Node struct {
	left                *Node
	right               *Node
	x, y, width, height int
}

func (node Node) isLeaf() bool {
	return node.left == nil || node.right == nil
}

func (pNode *Node) Split(rnd rand.Rand) {
	vert := RollDirection(*pNode, rnd)
	maxSplit := calcMaxSplit(vert, *pNode)

	//Dont Split if already to small
	if maxSplit <= MIN_SIZE {
		return
	}
	split := rnd.Intn(maxSplit-MIN_SIZE) + MIN_SIZE
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
		return node.width - MIN_SIZE
	} else {
		return node.height - MIN_SIZE
	}
}

func RollDirection(node Node, rnd rand.Rand) bool {
	if node.width > node.height && float64(node.width)/float64(node.height) > 2 {
		return true
	} else if node.height > node.width && float64(node.height)/float64(node.width) > 2 {
		return false
	} else {
		return rnd.Intn(2) == 1
	}
}

func (node *Node) SplitDeep(rnd rand.Rand, depth int) {
	if depth > 0 {
		node.Split(rnd)
		//TODO: Nicht immer splitten für mehr Variabilität
		if node.left != nil && node.right != nil {
			node.left.SplitDeep(rnd, depth-1)
			node.right.SplitDeep(rnd, depth-1)
		}
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
	root := Node{x: 0, y: 0, width: 40, height: 40}
	root.SplitDeep(*rnd, 5)
	//root.Split(*rnd)
	//root.right.Split(*rnd)

	img := root.Render()
	file, _ := os.Create("out.png")
	png.Encode(file, img)
}
