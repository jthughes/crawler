package main

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"os"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	args := os.Args[1:]

	if len(args) != 3 {
		fmt.Printf("usage: %s URL maxConcurrency maxPages", os.Args[0])
		os.Exit(1)
	}

	rawURL := args[0]
	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("\"%s\" is not a valid int\n", args[1])
		os.Exit(1)
	}

	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Printf("\"%s\" is not a valid int\n", args[1])
		os.Exit(1)
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("")
		os.Exit(1)
	}

	cfg := config{
		pages:              make(map[string]PageData),
		baseURL:            u,
		concurrencyControl: make(chan struct{}, maxConcurrency),
		mu:                 &sync.Mutex{},
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	fmt.Println("starting crawl of:", rawURL)
	start := time.Now()
	cfg.wg.Add(1)
	go cfg.crawlPage(rawURL)
	cfg.wg.Wait()
	elapse := time.Since(start)

	// for url, count := range cfg.pages {
	// 	fmt.Printf("%s: %d\n", url, count)
	// }

	fmt.Printf("Pages scraped in %s:\n", elapse)
	for url := range cfg.pages {
		fmt.Printf(" - %s\n", url)
	}

}
