package main

import (
	"image/png"
	"math/rand"
	"os"
	"time"
)

func main() {
	rndSource := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(rndSource)
	root := Node{x: 0, y: 0, width: 40, height: 40}
	root.SplitDeep(*rnd, 4)
	root.CreateLeafRooms(*rnd)
	//Connect siblings

	img := root.RenderRooms()
	file, _ := os.Create("out.png")
	png.Encode(file, img)
}
