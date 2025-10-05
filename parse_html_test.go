package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetH1FromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetH1FromHTMLMissing(t *testing.T) {
	inputBody := "<html><body><h2>Test Title</h2></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := ""

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetH1FromHTMLMultiple(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1><h1>Second title.</h1></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetH1FromHTMLNested(t *testing.T) {
	inputBody := "<html><body><h1>Test <b>Title</b></h1></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	inputBody := `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "Main paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLNoMain(t *testing.T) {
	inputBody := `<html><body>
		<p>First paragraph.</p>
		<p>Main paragraph.</p>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "First paragraph."

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLNoParagraph(t *testing.T) {
	inputBody := `<html><head><title>Title.</title></head><body>
	</body></html>`
	actual := getFirstParagraphFromHTML(inputBody)
	expected := ""

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body><a href="https://blog.boot.dev"><span>Boot.dev</span></a></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://blog.boot.dev"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetURLsFromHTMLRelative(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body><a href="/blog"><span>Boot.dev</span></a></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://blog.boot.dev/blog"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetImagesFromHTMLAbsolute(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body><img src="https://blog.boot.dev/logo.png" alt="Logo"></body></html>`

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	actual, err := getImagesFromHTML(inputBody, parsedURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://blog.boot.dev/logo.png"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetImagesFromHTMLRelative(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body><img src="/logo.png" alt="Logo"></body></html>`

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	actual, err := getImagesFromHTML(inputBody, parsedURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://blog.boot.dev/logo.png"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetImagesFromHTMLMultiple(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body>
		<img src="/logo.png" alt="Logo">
		<img src="https://cdn.boot.dev/banner.jpg">
	</body></html>`

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	actual, err := getImagesFromHTML(inputBody, parsedURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"https://blog.boot.dev/logo.png",
		"https://cdn.boot.dev/banner.jpg",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestExtractPageData(t *testing.T) {
	cases := []struct {
		name    string
		pageURL string
		html    string
		want    PageData
	}{
		{
			name:    "basic: h1, main paragraph, relative link and img",
			pageURL: "https://blog.boot.dev",
			html: `
<html>
  <body>
    <h1>Hello World</h1>
    <main><p>First paragraph inside main.</p></main>
    <a href="/about">About</a>
    <img src="/logo.png" alt="Logo">
  </body>
</html>`,
			want: PageData{
				URL:            "https://blog.boot.dev",
				H1:             "Hello World",
				FirstParagraph: "First paragraph inside main.",
				OutgoingLinks:  []string{"https://blog.boot.dev/about"},
				ImageURLs:      []string{"https://blog.boot.dev/logo.png"},
			},
		},
		{
			name:    "fallback paragraph when no <main>",
			pageURL: "https://blog.boot.dev",
			html: `
<html>
  <body>
    <h1>Title</h1>
    <p>Outside paragraph wins.</p>
    <a href="/x">x</a>
    <img src="/img.png">
  </body>
</html>`,
			want: PageData{
				URL:            "https://blog.boot.dev",
				H1:             "Title",
				FirstParagraph: "Outside paragraph wins.",
				OutgoingLinks:  []string{"https://blog.boot.dev/x"},
				ImageURLs:      []string{"https://blog.boot.dev/img.png"},
			},
		},
		{
			name:    "malformed HTML still parsed; absolute link and image",
			pageURL: "https://blog.boot.dev",
			html: `
<html body>
  <h1>Messy</h1>
  <a href="https://other.com/path">Other</a>
  <img src="https://cdn.boot.dev/banner.jpg">
</html body>`,
			want: PageData{
				URL:            "https://blog.boot.dev",
				H1:             "Messy",
				FirstParagraph: "", // no <p> present
				OutgoingLinks:  []string{"https://other.com/path"},
				ImageURLs:      []string{"https://cdn.boot.dev/banner.jpg"},
			},
		},
		{
			name:    "no h1 and no paragraph",
			pageURL: "https://blog.boot.dev",
			html: `
<html>
  <body>
    <a href="/only-link">Only link</a>
    <img src="/only.png">
  </body>
</html>`,
			want: PageData{
				URL:            "https://blog.boot.dev",
				H1:             "",
				FirstParagraph: "",
				OutgoingLinks:  []string{"https://blog.boot.dev/only-link"},
				ImageURLs:      []string{"https://blog.boot.dev/only.png"},
			},
		},
		{
			name:    "multiple links and images preserve order",
			pageURL: "https://blog.boot.dev",
			html: `
<html><body>
  <h1>t</h1>
  <main><p>p</p></main>
  <a href="/a1">a1</a>
  <a href="https://x.dev/a2">a2</a>
  <img src="/i1.png">
  <img src="https://x.dev/i2.png">
</body></html>`,
			want: PageData{
				URL:            "https://blog.boot.dev",
				H1:             "t",
				FirstParagraph: "p",
				OutgoingLinks: []string{
					"https://blog.boot.dev/a1",
					"https://x.dev/a2",
				},
				ImageURLs: []string{
					"https://blog.boot.dev/i1.png",
					"https://x.dev/i2.png",
				},
			},
		},
		{
			name:    "invalid base URL → empty link/image slices",
			pageURL: `:\\invalidBaseURL`,
			html: `
<html>
  <body>
    <h1>Title</h1>
    <p>Paragraph</p>
    <a href="/path">path</a>
    <img src="/logo.png">
  </body>
</html>`,
			want: PageData{
				URL:            `:\\invalidBaseURL`,
				H1:             "Title",
				FirstParagraph: "Paragraph",
				OutgoingLinks:  nil,
				ImageURLs:      nil,
			},
		},
	}

	for _, tc := range cases {
		tc := tc // shadow the loop variable.
		t.Run(tc.name, func(t *testing.T) {
			got := extractPageData(tc.html, tc.pageURL)

			if got.URL != tc.want.URL {
				t.Errorf("URL: want %q, got %q", tc.want.URL, got.URL)
			}
			if got.H1 != tc.want.H1 {
				t.Errorf("H1: want %q, got %q", tc.want.H1, got.H1)
			}
			if got.FirstParagraph != tc.want.FirstParagraph {
				t.Errorf("FirstParagraph: want %q, got %q", tc.want.FirstParagraph, got.FirstParagraph)
			}
			if !reflect.DeepEqual(got.OutgoingLinks, tc.want.OutgoingLinks) {
				t.Errorf("OutgoingLinks: want %v, got %v", tc.want.OutgoingLinks, got.OutgoingLinks)
			}
			if !reflect.DeepEqual(got.ImageURLs, tc.want.ImageURLs) {
				t.Errorf("ImageURLs: want %v, got %v", tc.want.ImageURLs, got.ImageURLs)
			}
		})
	}
}
