package internal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jackdanger/collectlinks"
)

type Parser struct {
	http.Client
}

func NewParser(c http.Client) Parser {
	return Parser{
		c,
	}
}

func (p *Parser) Parse(ctx context.Context, baseURI, uri string, links chan string) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return
	}

	resp, err := p.Do(request)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	pageLinks := collectlinks.All(resp.Body)

	var hrefs []string

	for _, link := range pageLinks {
		l := prepareLink(baseURI, link)
		if l != "" {
			hrefs = append(hrefs, l)
			go func() { links <- l }()
		}
	}

	fmt.Println(fmt.Sprintf("Visited %s, pageLinks: %v \n", uri, hrefs))
}

func prepareLink(base, href string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}

	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}

	if uri.Hostname() != "" && uri.Hostname() != baseUrl.Hostname() {
		return ""
	}

	uri = baseUrl.ResolveReference(uri)

	return uri.String()
}
