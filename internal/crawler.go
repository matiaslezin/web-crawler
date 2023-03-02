package internal

import (
	"context"
	"net/http"
	"time"
)

type Crawler struct {
	domain string
	Parser
}

func NewCrawler(domain string) Crawler {
	return Crawler{
		domain: domain,
		Parser: NewParser(http.Client{
			Timeout: 60 * time.Second,
		}),
	}
}

func (c *Crawler) Crawl(ctx context.Context, done chan interface{}, pending chan int, links, processing chan string) {
	for {
		select {
		case uri := <-processing:
			c.Parse(ctx, c.domain, uri, links)
			pending <- -1
		case <-done:
			return
		}
	}
}
