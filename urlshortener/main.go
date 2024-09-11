//go:build !solution

package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

type resBody struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

type reqBody struct {
	URL string
}

var (
	storage map[string]string
	mu      sync.Mutex
)

func rS(str string) (resource string, e error) {
	al := sha1.New()
	_, err := al.Write([]byte(str))

	if err != nil {
		e = err
		return
	}

	resource = hex.EncodeToString(al.Sum(nil))
	return
}

func main() {
	port := flag.String("port", "80", "http server port")
	flag.Parse()

	storage = make(map[string]string)

	http.HandleFunc(
		"/shorten",
		func(w http.ResponseWriter, r *http.Request) {
			var rb reqBody
			decoder := json.NewDecoder(r.Body)
			if decoder.Decode(&rb) != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}

			short, err := rS(rb.URL)

			if err != nil {
				http.Error(w, "server error", http.StatusInternalServerError)
			}

			mu.Lock()
			storage[short] = rb.URL
			mu.Unlock()

			data := resBody{URL: rb.URL, Key: short}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(data)

			if err != nil {
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		},
	)
	http.HandleFunc(
		"/go/",
		func(w http.ResponseWriter, r *http.Request) {
			p := strings.Split(r.URL.Path, "/")
			key := p[len(p)-1]

			mu.Lock()
			url, exist := storage[key]
			mu.Unlock()

			if !exist {
				http.Error(w, "key not found", http.StatusNotFound)
				return
			}

			w.Header().Add("Location", url)
			w.WriteHeader(http.StatusFound)
		},
	)

	host := fmt.Sprintf(":%s", *port)
	log.Fatal(http.ListenAndServe(host, nil))
}
