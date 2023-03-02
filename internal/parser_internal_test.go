package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_prepareLink(t *testing.T) {
	var cases = []struct {
		name     string
		baseURI  string
		href     string
		expected string
	}{
		{
			name:     "baseURI and href with same domain",
			baseURI:  "https://google.com/",
			href:     "https://google.com/news",
			expected: "https://google.com/news",
		},
		{
			name:     "href missing domain",
			baseURI:  "https://google.com/",
			href:     "/news",
			expected: "https://google.com/news",
		},
		{
			name:     "href with wrong domain",
			baseURI:  "https://google.com/",
			href:     "https://facebook.com/",
			expected: "",
		},
		{
			name:     "href with wrong sub domain",
			baseURI:  "https://google.com/",
			href:     "https://news.google.com/",
			expected: "",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result := prepareLink(tt.baseURI, tt.href)
			assert.Equal(t, tt.expected, result)
		})
	}
}
