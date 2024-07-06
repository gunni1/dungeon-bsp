package main

import (
	"bytes"
	"encoding/base64"
	dbsp "gunni1/dungeon-bsp/dbsp"
	"html/template"
	"image"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var ImageTemplate string = `<!DOCTYPE html>
    <html lang="en"><head></head>
    <body>
	  <img src="data:image/png;base64,{{.Image}}">
	  <table>
	    <tbody><tr><td>seed</td><td>{{.Seed}}</td></tr></tbody>
	  </table>
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

		rndSource := rand.NewSource(seed)
		rnd := rand.New(rndSource)
		root := dbsp.Node{X: 0, Y: 0, Width: width, Height: height}
		root.SplitDeep(*rnd, depth)
		root.CreateLeafRooms(*rnd)
		//TODO: Connect siblings

		img := root.RenderRooms()

		data := map[string]interface{}{"Seed": seed}
		writeImageWithTemplate(w, &img, data)
	})
	http.ListenAndServe(":3000", nil)
}

func writeImageWithTemplate(w http.ResponseWriter, img *image.Image, data map[string]interface{}) {
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, *img); err != nil {
		log.Println("unable to encode image.")
	}

	encodedImage := base64.StdEncoding.EncodeToString(buffer.Bytes())
	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		log.Println("unable to parse image template.")
	} else {
		data["Image"] = encodedImage
		if err = tmpl.Execute(w, data); err != nil {
			log.Println("unable to execute template.")
		}
	}
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
