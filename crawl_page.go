package main

import (
	"fmt"
	"strings"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	normBaseURL, err := normalizeURL(rawBaseURL)
	if err != nil {
		return
	}
	normCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}
	if strings.Index(normCurrentURL, normBaseURL) != 0 {
		return
	}
	if count, ok := pages[normCurrentURL]; ok {
		pages[normCurrentURL] = count + 1
		return
	}
	pages[normCurrentURL] = 1
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}
	fmt.Printf("Got html from %s\n", rawCurrentURL)
	page := extractPageData(html, rawCurrentURL)
	for index := range page.OutgoingLinks {
		crawlPage(rawBaseURL, page.OutgoingLinks[index], pages)
	}
}
