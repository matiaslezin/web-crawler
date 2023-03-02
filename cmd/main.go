package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"web-crawler/internal"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("missing base URI")
		os.Exit(1)
	}

	baseURI := args[0]

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	links := make(chan string)
	pending := make(chan int)
	processing := make(chan string)

	done := make(chan interface{})

	go checkLinks(links, processing, pending)
	go checkPending(done, pending)

	g := runtime.GOMAXPROCS(0)
	for i := 0; i < g; i++ {
		crawler := internal.NewCrawler(baseURI)
		go crawler.Crawl(ctx, done, pending, links, processing)
	}

	fmt.Println("start crawling...")
	go func() { links <- baseURI }()

	select {
	case <-done:
		fmt.Println("finished successfully")
	case <-ctx.Done():
		fmt.Println("process is taking too long, shutting down...")
	}
}

func checkLinks(links chan string, processing chan string, pending chan int) {
	processed := make(map[string]bool)

	for l := range links {
		if !processed[l] {
			processed[l] = true
			processing <- l
			pending <- 1
		}
	}
}

func checkPending(done chan interface{}, pending chan int) {
	var count int

	for c := range pending {
		count += c
		if count == 0 {
			done <- struct{}{}
		}
	}
}
