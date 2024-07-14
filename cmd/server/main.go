package main

import (
	"bytes"
	"encoding/base64"
	dbsp "gunni1/dungeon-bsp/dbsp"
	"html/template"
	"image"
	"image/gif"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var imagePage string = `<!DOCTYPE html>
    <html lang="en"><head></head>
    <body>
	  <img src="data:image/png;base64,{{.Image}}">
	  <table>
	    <tbody><tr><td>seed</td><td>{{.Seed}}</td></tr></tbody>
	  </table>
	  <img src="data:image/gif;base64,{{.BSPGif}}">
	</body>
	`

func main() {

	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {
		seed, err := parseSeed(r)
		if err != nil {
			seed = time.Now().UnixNano()
		}
		width := asInt(r.URL.Query().Get("width"))
		height := asInt(r.URL.Query().Get("height"))
		depth := asInt(r.URL.Query().Get("depth"))
		//FIXME: error handling on missing input parameters

		rndSource := rand.NewSource(seed)
		rnd := rand.New(rndSource)
		root := dbsp.Node{X: 0, Y: 0, Width: width, Height: height}
		prtcCtx := dbsp.ProtocolCtx{InterimResults: make(chan dbsp.Node), RootNode: &root}

		go func() {
			root.SplitDeep(*rnd, depth, prtcCtx)
			root.CreateLeafRooms(*rnd)
			close(prtcCtx.InterimResults)
		}()

		//Create result Gif
		bspGif := gif.GIF{}
		for result := range prtcCtx.InterimResults {
			bspGif.Image = append(bspGif.Image, result.RenderNodePaletted().(*image.Paletted))
			bspGif.Delay = append(bspGif.Delay, 100)
		}

		roomImg := root.RenderRooms()

		data := map[string]interface{}{"Seed": seed}

		templ := template.Must(template.New("image").Parse(imagePage))

		data["Image"] = encodePngB64(roomImg)
		data["BSPGif"] = encodeGifB64(&bspGif)
		if err = templ.Execute(w, data); err != nil {
			log.Println("unable to execute template")
		}
	})

	http.ListenAndServe(":3000", nil)
}

func encodeGifB64(input *gif.GIF) string {
	buffer := new(bytes.Buffer)
	gif.EncodeAll(buffer, input)
	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}

func encodePngB64(input image.Image) string {
	buffer := new(bytes.Buffer)
	png.Encode(buffer, input)
	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}

func asInt(value string) int {
	intVal, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return intVal
}

func parseSeed(r *http.Request) (int64, error) {
	seedParam := r.URL.Query().Get("seed")
	return strconv.ParseInt(seedParam, 10, 64)
}
