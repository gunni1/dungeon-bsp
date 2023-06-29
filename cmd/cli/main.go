package main

import (
	"flag"
	dbsp "gunni1/dungeon-bsp/dbsp"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func main() {
	rndSource := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(rndSource)

	var width int
	var height int
	var depth int
	flag.IntVar(&width, "width", 40, "Total width of the area.")
	flag.IntVar(&height, "height", 40, "total height of the area.")
	flag.IntVar(&depth, "depth", 4, "Number of splits / depth of binary tree")
	flag.Parse()

	root := dbsp.Node{X: 0, Y: 0, Width: width, Height: height}
	root.SplitDeep(*rnd, depth)
	root.CreateLeafRooms(*rnd)
	//TODO: Connect siblings

	img := root.RenderRooms()
	file, _ := os.Create("out.png")
	png.Encode(file, img)
}
