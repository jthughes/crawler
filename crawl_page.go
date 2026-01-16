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
		fmt.Println("=> Reached max page count")
		return
	}

	fmt.Printf("Scraping: %s\n", rawCurrentURL)
	normCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println("=> Error normalizing url")
		return
	}
	normBaseURL, err := normalizeURL(cfg.baseURL.String())
	if strings.Index(normCurrentURL, normBaseURL) != 0 {
		fmt.Printf("=> \"%s\" not in \"%s\"\n", normCurrentURL, cfg.baseURL.String())
		return
	}
	cfg.mu.Lock()
	_, ok := cfg.pages[normCurrentURL]
	cfg.mu.Unlock()
	if ok {
		fmt.Println("=> Already visited page")
		return
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("=> Error getting html: %s\n", err)
		return
	}
	fmt.Printf("=> Got html from %s\n", rawCurrentURL)
	page := extractPageData(html, rawCurrentURL)
	cfg.mu.Lock()
	cfg.pages[normCurrentURL] = page
	cfg.mu.Unlock()
	for _, url := range page.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}
}
