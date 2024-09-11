//go:build !solution

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	files := os.Args[1:]
	for _, fileName := range files {
		file, _ := os.Open(fileName)
		countLines(file, counts)
		file.Close()
	}
	for line, count := range counts {
		if count >= 2 {
			fmt.Printf("%d\t%s\n", count, line)
		}
	}
}

func countLines(f *os.File, counts map[string]int) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		counts[scanner.Text()]++
	}
}
