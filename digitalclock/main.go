//go:build !solution

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	port := flag.String("port", "80", "http server port")
	flag.Parse()

	http.HandleFunc("/", getImage)

	host := fmt.Sprintf(":%s", *port)
	log.Fatal(http.ListenAndServe(host, nil))
}

func getImage(w http.ResponseWriter, r *http.Request) {
	var sP = 1
	var tP string
	var err error

	q := r.URL.Query()

	if sslc, ok := q["k"]; ok && len(sslc) > 0 {
		sP, err = strconv.Atoi(sslc[0])
		if err != nil || sP < 1 || sP > 30 {
			http.Error(w, "invalid k", http.StatusBadRequest)
			return
		}
	}

	if tSlc, ok := q["time"]; ok {
		tP = tSlc[0]
		tPatt := "^(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$"
		ok, err = regexp.MatchString(tPatt, tP)

		if err != nil || !ok {
			http.Error(w, "invalid time", http.StatusBadRequest)
			return
		}
	} else {
		t := time.Now()
		tP = fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
	}

	fmt.Println(sP, tP)
	img := CreateImage(tP, sP)

	w.Header().Set("Content-Type", "image/png")
	err = png.Encode(w, img)

	if err != nil {
		http.Error(w, "render error", http.StatusBadRequest)
	}
}

func CreateImage(time string, size int) *image.RGBA {
	cW := 4
	dW := 8
	dH := 12

	iW := (6*dW + 2*cW) * size
	iH := dH * size

	uL := image.Point{X: 0, Y: 0}
	lR := image.Point{X: iW, Y: iH}

	img := image.NewRGBA(image.Rectangle{Min: uL, Max: lR})

	offset := 0
	d := ""

	for _, char := range time {
		switch string(char) {
		case "0":
			d = Zero
		case "1":
			d = One
		case "2":
			d = Two
		case "3":
			d = Three
		case "4":
			d = Four
		case "5":
			d = Five
		case "6":
			d = Six
		case "7":
			d = Seven
		case "8":
			d = Eight
		case "9":
			d = Nine
		default:
			d = Colon
		}

		lines := strings.Split(d, "\n")

		for h := 0; h < len(lines); h++ {
			for w, char := range lines[h] {
				x := w + offset
				y := h
				ch := string(char)
				for i := x * size; i < (x+1)*size; i++ {
					for j := y * size; j < (y+1)*size; j++ {
						if ch != "." {
							img.Set(i, j, Cyan)
						} else {
							img.Set(i, j, color.RGBA{R: 255, G: 255, B: 255, A: 255})
						}
					}
				}
			}
		}

		if d == Colon {
			offset += cW
		} else {
			offset += dW
		}
	}

	return img
}
