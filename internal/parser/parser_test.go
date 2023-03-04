package parser_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"web-crawler/internal/parser"

	"github.com/stretchr/testify/assert"
)

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func TestParser_Parse(t *testing.T) {
	var cases = []struct {
		name     string
		client   http.Client
		uri      string
		href     string
		expected []string
	}{
		{
			name: "find about link",
			client: http.Client{Transport: roundTripFunc(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`<li><a class="Button_button__30ukX Button_dark__hKxQ0 Button_link__Slw7D Button_left__P5ZhN" rel="" href="https://monzo.com/about/"><span class="Button_text__d9Ltb">About us</span></a></li>`)),
				}
			})},
			uri:      "https://monzo.com",
			expected: []string{"https://monzo.com/about/"},
		},
		{
			name: "find no link",
			client: http.Client{Transport: roundTripFunc(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`<div class="Row_Row__k7R1L Row_center__GJVq4"><a aria-label="Download on the App Store" href="https://app.adjust.com/ydi27sn?engagement_type=fallback_click&amp;fallback=https%3A%2F%2Fmonzo.com%2Fdownload&amp;redirect_macos=https%3A%2F%2Fmonzo.com%2Fdownload" class="AppStoreButton_appStoreButton__dXRYU undefined undefined"><img width="150" loading="lazy" src="https://images.ctfassets.net/ro61k101ee59/47R5HVs5mn8kaXIb4IgfMJ/678738cebf35e278307695136c1be015/Download_on_the_App_Store_Badge_US-UK_RGB_blk_092917.svg" alt=""/></a><a aria-label="Get it on Google Play" href="https://app.adjust.com/9mq4ox7?engagement_type=fallback_click&amp;fallback=https%3A%2F%2Fmonzo.com%2Fdownload&amp;redirect_macos=https%3A%2F%2Fmonzo.com%2Fdownload" class="AppStoreButton_playStoreButton__9Tg16 undefined undefined"><img width="191" loading="lazy" src="https://images.ctfassets.net/ro61k101ee59/hegpnNpXJdP4GkCOqHf5K/5a7d44685d8703460ca8db539f8401d5/google-play-badge.png" alt=""/></a></div>`)),
				}
			})},
			uri:      "https://monzo.com",
			expected: nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockParser := parser.Parser{
				Client: tt.client,
			}

			links := mockParser.Parse(context.Background(), "https://monzo.com", tt.uri)

			assert.Equal(t, tt.expected, links)
		})
	}
}
