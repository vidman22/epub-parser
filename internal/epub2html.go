package parser

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"net/url"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

type Params struct {
	rootDir       string
	manifestItems []Item
	spineItemRefs []Itemref
	tocMap        map[string]string
	title         string
	r             *zip.ReadCloser
}

// contentMap is map[fullContentPath]Title
func processEpubContent(params Params) ([]Content, Cover, error) {
	manifestItems := params.manifestItems
	rootDir := params.rootDir
	spineItemRefs := params.spineItemRefs
	tocMap := params.tocMap
	r := params.r

	//var cover *Cover
	manifestIDMap := make(map[string]string)
	manifestHrefMap := make(map[string]Item)
	likelyCoverHref := ""
	for _, item := range manifestItems {
		fullHref := filepath.Join(rootDir, item.Href)
		manifestIDMap[item.Id] = fullHref
		manifestHrefMap[fullHref] = item
		if strings.Contains(item.Id, "cover") {
			likelyCoverHref = fullHref
		}
	}

	var texts []Content

	for _, itemRef := range spineItemRefs {
		contentFilePath, ok := manifestIDMap[itemRef.Idref]
		if !ok {
			continue
		}

		// the toc map isn't guaranteed to have the titles for all the spine items unfortunately
		Title := tocMap[contentFilePath]

		if strings.Contains(itemRef.Idref, "cover") {
			continue
		}
		var combinedHTML strings.Builder

		fileData, err := readZipFile(r, contentFilePath)
		if err != nil {
			continue
		}

		doc, err := html.Parse(bytes.NewReader(fileData))
		if err != nil {
			continue
		}

		possibleTitle := extractRawHTML(doc, &combinedHTML, r, contentFilePath, manifestHrefMap)
		combinedHTML.WriteString("\n<hr />\n")
		stringHtml := combinedHTML.String()
		if Title == "" {
			Title = possibleTitle[0:int(math.Min(float64(len(possibleTitle)), 50))]
		}

		texts = append(texts, Content{Html: stringHtml, Title: Title})
	}
	var cover Cover
	if likelyCoverHref != "" {
		coverData, err := readZipFile(r, likelyCoverHref)
		if err == nil {
			filename := filepath.Base(likelyCoverHref)
			ext := filepath.Ext(filename)

			cover = Cover{
				FileName: filename,
				Ext:      ext,
				File:     coverData,
			}
		}
	}
	return texts, cover, nil
}

func readZipFile(r *zip.ReadCloser, filePath string) ([]byte, error) {
	cleanPath := filepath.Clean(filePath)
	if strings.HasPrefix(cleanPath, "..") {
		return nil, fmt.Errorf("invalid path trying to access parent directory: %s", filePath)
	}

	for _, f := range r.File {
		if f.Name == cleanPath {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open %s: %w", cleanPath, err)
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf("file %s not found in archive", cleanPath)
}

func extractRawHTML(n *html.Node, w io.StringWriter, r *zip.ReadCloser, contentFilePath string, manifestHrefMap map[string]Item) string {
	var findBodyAndExtract func(*html.Node)
	foundBody := false
	isFirstChild := false
	firstText := ""
	findBodyAndExtract = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			foundBody = true
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				isFirstChild = true
				titleString := renderNodeRaw(isFirstChild, c, w, r, contentFilePath, manifestHrefMap)
				if titleString != "" {
					firstText = titleString
				}
				isFirstChild = false
			}
			return
		}

		if !foundBody {
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				findBodyAndExtract(c)
				if foundBody {
					break
				}
			}
		}
	}

	findBodyAndExtract(n)

	return firstText
}

func renderNodeRaw(isFirstChild bool, n *html.Node, w io.StringWriter, r *zip.ReadCloser, contentFilePath string, manifestHrefMap map[string]Item) string {
	switch n.Type {
	case html.TextNode:
		w.WriteString(n.Data)
		if isFirstChild {
			return n.Data
		}
	case html.ElementNode:
		tag := n.Data
		switch tag {

		case "script", "style", "link", "meta", "head", "title", "svg":
			return ""
		}

		if tag == "img" {
			var src string
			for i, attr := range n.Attr {
				if attr.Key == "src" {
					src = attr.Val
					// Remove the original src attribute to replace it
					n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
					break
				}
			}

			if src != "" {
				// Resolve the image path relative to the current content file
				imagePath, err := url.JoinPath(filepath.Dir(contentFilePath), src)
				if err != nil {
					return ""
				}

				imageData, err := readZipFile(r, imagePath)
				if err != nil {
					return ""
				}

				item, ok := manifestHrefMap[imagePath]
				if !ok {
					return ""
				}
				mediaType := item.MediaType

				encodedData := base64.StdEncoding.EncodeToString(imageData)
				dataURI := fmt.Sprintf("data:%s;base64,%s", mediaType, encodedData)

				// Add the new src attribute with the data URI
				n.Attr = append(n.Attr, html.Attribute{Key: "src", Val: dataURI})
			}
		}

		var openTag strings.Builder
		openTag.WriteString("<")
		openTag.WriteString(tag)

		for _, attr := range n.Attr {
			if attr.Key == "class" {
				continue
			}
			openTag.WriteString(" ")
			openTag.WriteString(attr.Key)
			openTag.WriteString(`="`)
			openTag.WriteString(html.EscapeString(attr.Val))
			openTag.WriteString(`"`)
		}
		openTag.WriteString(">")
		w.WriteString(openTag.String())

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			renderNodeRaw(isFirstChild, c, w, r, contentFilePath, manifestHrefMap)
		}
		if n.FirstChild != nil || tag != "img" { // Self-closing for img if no children
			w.WriteString("</" + tag + ">")
		}

	case html.CommentNode:
		return ""
	case html.DoctypeNode:
		return ""
	}
	return ""
}
