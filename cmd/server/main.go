package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
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

	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {
		seed, err := parseSeed(r)
		if err != nil {
			fmt.Println("No seed found. Default to current time.")
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

func parseSeed(r *http.Request) (int64, error) {
	seedParam := r.URL.Query().Get("seed")
	return strconv.ParseInt(seedParam, 10, 64)
}
