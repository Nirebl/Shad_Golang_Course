//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	urls := os.Args[1:]
	for _, url := range urls {
		res, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		content, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fmt.Println(string(content))
	}
}
