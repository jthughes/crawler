package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s%s", u.Host, u.Path)
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return url, nil
}
