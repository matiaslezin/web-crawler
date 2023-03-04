package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"web-crawler/internal"
	"web-crawler/internal/parser"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("missing base URI")
		os.Exit(1)
	}

	baseURI := args[0]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})

	c := crawler.New(buildParser())
	go c.Run(ctx, done, baseURI, baseURI)

	select {
	case <-done:
		fmt.Println("finished successfully")
	case <-ctx.Done():
		fmt.Println("process is taking too long, shutting down...")
	}
}

func buildParser() func(ctx context.Context, baseURI, uri string) []string {
	return parser.New(&http.Client{
		Timeout: 60 * time.Second,
	}).Parse
}
