package parser

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jackdanger/collectlinks"
)

type Parser struct {
	*http.Client
}

func New(c *http.Client) Parser {
	if c == nil {
		panic("nil http.Client")
	}

	return Parser{
		Client: c,
	}
}

func (p Parser) Parse(ctx context.Context, baseURI, uri string) []string {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil
	}

	resp, err := p.Do(request)
	if err != nil {
		fmt.Println(fmt.Sprintf("request to %s, err: %s", uri, err.Error()))
		return nil
	}

	defer resp.Body.Close()

	pageLinks := collectlinks.All(resp.Body)

	var hrefs []string

	for _, link := range pageLinks {
		l := prepareLink(baseURI, link)
		if l != "" {
			hrefs = append(hrefs, l)
		}
	}

	return hrefs
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
