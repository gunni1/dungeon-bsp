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
    <body><img src="data:image/png;base64,{{.Image}}"></body>`

func main() {
	rndSource := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(rndSource)

	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {

		width := asInt(r.URL.Query().Get("width"))
		height := asInt(r.URL.Query().Get("height"))
		depth := asInt(r.URL.Query().Get("depth"))
		root := dbsp.Node{X: 0, Y: 0, Width: width, Height: height}
		root.SplitDeep(*rnd, depth)
		root.CreateLeafRooms(*rnd)
		//TODO: Connect siblings

		img := root.RenderRooms()
		writeImageWithTemplate(w, &img)
	})
	http.ListenAndServe(":8080", nil)
}

func writeImageWithTemplate(w http.ResponseWriter, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, *img); err != nil {
		log.Println("unable to encode image.")
	}

	encodedImage := base64.StdEncoding.EncodeToString(buffer.Bytes())
	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		log.Println("unable to parse image template.")
	} else {
		data := map[string]interface{}{"Image": encodedImage}
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
