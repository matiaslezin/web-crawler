package crawler

import (
	"context"
	"fmt"
	"sync"
)

type (
	parserFunc func(ctx context.Context, baseURI, uri string) []string

	Crawler struct {
		parserFunc
	}
)

func New(parser func(ctx context.Context, baseURI, uri string) []string) Crawler {
	return Crawler{
		parser,
	}
}

func (c Crawler) Run(ctx context.Context, done chan<- struct{}, baseURI, uri string) {
	fmt.Println(fmt.Sprintf("parsing %s \n", uri))

	var (
		wg        sync.WaitGroup
		processed sync.Map
	)

	wg.Add(1)
	go c.crawl(ctx, baseURI, uri, &wg, &processed)

	wg.Wait()

	done <- struct{}{}
}

func (c Crawler) crawl(ctx context.Context, baseURI, uri string, wg *sync.WaitGroup, processed *sync.Map) {
	defer wg.Done()

	links := c.parserFunc(ctx, baseURI, uri)

	fmt.Println(fmt.Sprintf("Visited %s, links: %v \n", uri, links))

	for _, link := range links {
		if _, ok := processed.LoadOrStore(link, nil); ok {
			continue
		}

		wg.Add(1)
		go c.crawl(ctx, baseURI, link, wg, processed)
	}
}
