//go:build !solution

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	c, _, err := websocket.DefaultDialer.Dial(*addr, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}
			fmt.Print(string(message))
		}
	}()
	mes := func(r io.Reader) <-chan string {
		lines := make(chan string)
		go func() {
			defer close(lines)
			scan := bufio.NewScanner(r)
			for scan.Scan() {
				lines <- scan.Text()
			}
		}()
		return lines
	}(os.Stdin)
	for {
		select {
		case <-done:
			return
		case <-stop:
			log.Println("interrupted")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return
			}
			select {
			case <-done:
			case <-time.After(time.Millisecond * 100):
			}
			return
		case text := <-mes:
			err := c.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				log.Println("write:", err)
				return
			}

		}
	}
}
