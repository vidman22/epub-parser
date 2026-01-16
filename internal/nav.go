package parser

import (
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

// ParseNavDoc parses the EPUB 3 Navigation Document (NAV.xhtml) or toc.xhtml
// Returns a map of {Cleaned href (without hash): Link text content}
func ParseNavDoc(fBytes []byte, rootDir string) (map[string]string, error) {
	tocMap := make(map[string]string)
	doc, err := html.Parse(strings.NewReader(string(fBytes)))
	if err != nil {
		return nil, err
	}

	// Find all <a> tags in the document
	var collectLinks func(*html.Node)
	collectLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var href string

			// Get href attribute
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href = attr.Val
					break
				}
			}

			// Extract text content from the <a> tag
			var textBuilder strings.Builder
			var extractText func(*html.Node)
			extractText = func(tn *html.Node) {
				if tn.Type == html.TextNode {
					textBuilder.WriteString(tn.Data)
				}
				for c := tn.FirstChild; c != nil; c = c.NextSibling {
					extractText(c)
				}
			}
			extractText(n)
			contentText := strings.TrimSpace(textBuilder.String())

			if href != "" && contentText != "" {
				// Remove hash/anchor from href (e.g., "file.xhtml#section" -> "file.xhtml")
				cleanHref := href
				if idx := strings.Index(cleanHref, "#"); idx >= 0 {
					cleanHref = cleanHref[:idx]
				}

				// Use cleaned href as key, content text as value
				// Only add if this href hasn't been seen before
				if cleanHref != "" {
					if _, exists := tocMap[cleanHref]; !exists {
						fullPath := filepath.Join(rootDir, cleanHref)
						tocMap[fullPath] = contentText
					}
				}
			}
		}

		// Recursively traverse all children
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			collectLinks(c)
		}
	}
	collectLinks(doc)

	return tocMap, nil
}
