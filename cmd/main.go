package main

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"web-crawler/internal"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	parser := internal.NewParser(http.Client{
		Timeout: 500 * time.Millisecond,
	})

	baseURI := "https://google.com/"

	links := make(chan string)
	pendingLinks := make(chan int)
	processingLinks := make(chan string)

	done := make(chan interface{})
	doneTotal := make(chan interface{})

	go checkLinks(links, processingLinks, pendingLinks)
	go checkPending(done, pendingLinks)

	g := runtime.GOMAXPROCS(0)
	for i := 0; i < g; i++ {
		go func(parser internal.Parser) {
			for {
				select {
				case uri := <-processingLinks:
					parser.Parse(ctx, baseURI, uri, links)
					pendingLinks <- -1
				case <-done:
					doneTotal <- struct{}{}
					return
				}
			}
		}(parser)
	}

	fmt.Println("start crawler...")
	go func() { links <- baseURI }()

	select {
	case <-doneTotal:
		fmt.Println("finished successfully")
	case <-ctx.Done():
		fmt.Println("process is taking too long, shutting down...")
	}
}

func checkLinks(links chan string, processingLinks chan string, pendingLinks chan int) {
	processed := make(map[string]bool)

	for l := range links {
		if !processed[l] {
			processed[l] = true
			processingLinks <- l
			pendingLinks <- 1
		}
	}
}

func checkPending(done chan interface{}, pendingLinks chan int) {
	var count int

	for c := range pendingLinks {
		count += c
		if count == 0 {
			done <- struct{}{}
		}
	}
}
