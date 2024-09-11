//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func fetch(url string, messages chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		messages <- fmt.Sprintf("Error fetching %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		messages <- fmt.Sprintf("Error reading %s: %v", url, err)
		return
	}

	duration := time.Since(start)
	messages <- fmt.Sprintf("%.2fs elapsed\t%s", duration.Seconds(), url)
}

func main() {
	start := time.Now()
	messages := make(chan string)

	urls := os.Args[1:]
	for _, url := range urls {
		go fetch(url, messages)
	}

	for range urls {
		fmt.Println(<-messages)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}
