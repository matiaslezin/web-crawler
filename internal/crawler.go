package crawler

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"web-crawler/internal/parser"
)

type Crawler struct {
	domain string
	parser.Parser
}

func New(domain string) Crawler {
	return Crawler{
		domain: domain,
		Parser: parser.New(&http.Client{
			Timeout: 60 * time.Second,
		}),
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

	links := c.Parse(ctx, baseURI, uri)

	fmt.Println(fmt.Sprintf("Visited %s, links: %v \n", uri, links))

	for _, link := range links {
		if _, ok := processed.LoadOrStore(link, nil); ok {
			continue
		}

		wg.Add(1)
		go c.crawl(ctx, baseURI, link, wg, processed)
	}
}
