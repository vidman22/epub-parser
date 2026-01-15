package parser

import (
	"encoding/xml"
	"path/filepath"
	"strings"
)

type NCX struct {
	XMLName   xml.Name   `xml:"ncx"`
	NavPoints []NavPoint `xml:"navMap>navPoint"`
}

type NavPoint struct {
	NavLabel  NavLabel   `xml:"navLabel"`
	Content   XMLContent `xml:"content"`
	NavPoints []NavPoint `xml:"navPoint"` // Add this line for nested navPoints
}

type NavLabel struct {
	Text string `xml:"text"`
}

type XMLContent struct {
	Src string `xml:"src,attr"`
}

// ParseNcx parses the EPUB 2 NCX file (toc.ncx)
// Returns a map of {Full Content Path: Title}
func ParseNcx(fBytes []byte, rootDir string) (map[string]string, error) {
	var ncx NCX
	err := xml.Unmarshal(fBytes, &ncx)
	if err != nil {
		return nil, err
	}

	tocMap := make(map[string]string)

	// Process all navPoints recursively
	for _, navPoint := range ncx.NavPoints {
		processNavPoint(navPoint, rootDir, tocMap)
	}

	return tocMap, nil
}

// processNavPoint recursively processes a navPoint and its children
func processNavPoint(navPoint NavPoint, rootDir string, tocMap map[string]string) {
	title := strings.TrimSpace(navPoint.NavLabel.Text)
	href := navPoint.Content.Src

	if title != "" && href != "" {
		cleanHref := strings.SplitN(href, "#", 2)[0]
		fullPath := filepath.Join(rootDir, cleanHref)
		if _, exists := tocMap[fullPath]; !exists {
			tocMap[fullPath] = title
		}
	}

	// Recursively process nested navPoints
	for _, childNavPoint := range navPoint.NavPoints {
		processNavPoint(childNavPoint, rootDir, tocMap)
	}
}
