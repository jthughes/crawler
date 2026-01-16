package main

import (
	"fmt"
	"net/url"
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
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	rawURL := args[0]

	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("")
		os.Exit(1)
	}

	cfg := config{
		pages:              make(map[string]PageData),
		baseURL:            u,
		concurrencyControl: make(chan struct{}, 5),
		mu:                 &sync.Mutex{},
		wg:                 &sync.WaitGroup{},
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
