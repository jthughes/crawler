package main

import (
	"fmt"

	"os"
)

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
	fmt.Println("starting crawl of:", rawURL)

	pages := map[string]int{}
	crawlPage(rawURL, rawURL, pages)
	for url, count := range pages {
		fmt.Printf("%s: %d\n", url, count)
	}
}
