package main

import (
	"fmt"
	"strings"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	cfg.mu.Lock()
	page_count := len(cfg.pages)
	cfg.mu.Unlock()
	if page_count >= cfg.maxPages {
		return
	}

	normCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println("=> Error normalizing url")
		return
	}
	normBaseURL, err := normalizeURL(cfg.baseURL.String())
	if strings.Index(normCurrentURL, normBaseURL) != 0 {
		return
	}
	cfg.mu.Lock()
	_, ok := cfg.pages[normCurrentURL]
	cfg.mu.Unlock()
	if ok {
		return
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("=> Error getting html: %s\n", err)
		return
	}
	page := extractPageData(html, rawCurrentURL)
	cfg.mu.Lock()
	cfg.pages[normCurrentURL] = page
	cfg.mu.Unlock()
	for _, url := range page.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}
}
