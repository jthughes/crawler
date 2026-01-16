package main

import (
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	url := u.Host + u.Path
	url = strings.ToLower(url)
	url = strings.TrimSuffix(url, "/")

	return url, nil
}
