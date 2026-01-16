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

	html, err := getHTML(rawURL)
	if err != nil {
		fmt.Printf("failed to get html: %v", err)
		os.Exit(1)
	}
	fmt.Println(html)
}
