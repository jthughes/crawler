package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func writeCSVReport(pages map[string]PageData, filename string) error {

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	w := csv.NewWriter(file)
	w.Write([]string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"})
	for _, page := range pages {
		w.Write([]string{page.URL, page.H1, page.FirstParagraph, strings.Join(page.OutgoingLinks, ";"), strings.Join(page.ImageURLs, ";")})
	}

	return nil
}
