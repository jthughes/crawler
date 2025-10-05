package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getH1FromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	header := doc.Find("h1").First().Text()
	return header
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	paragraphs := doc.Find("main").Find("p").First()
	if len(paragraphs.Nodes) != 0 {
		return paragraphs.Text()
	}
	return doc.Find("p").First().Text()
}
