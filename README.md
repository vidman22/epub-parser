# EPUB Parser

A Go library for parsing EPUB files and converting them to HTML content with embedded images.

## Features

- Parse EPUB 2.0 and 3.0 files
- Extract content with proper chapter/section ordering
- Convert images to base64-encoded data URIs
- Preserve document structure and formatting
- Extract table of contents and metadata
- Handle nested content files and relative paths
- Filter out unnecessary elements (scripts, styles, SVG)

## Installation
Add the library to your project using go get:

go get github.com/vidman22/epub-parser

## Acknowledgments

This library was inspired by [mathieu-keller/epub-parser](https://github.com/mathieu-keller/epub-parser).

## tag

git tag -a v0.1.8 -m "version 0.1.8" \n
git push origin v0.1.8

## License
This project is licensed under the MIT License. See the LICENSE file for details.

## Usage 

```aiignore
package main

import (
	"fmt"
	"log"

	"github.com/vidman22/epub-parser"
)

func main() {
	//epub.ParseEpub("/Users/johnvidmar/Downloads/augustine-of-hippo_the-city-of-god_marcus-dods-george-wilson-j-j-smith.epub")
	book2, err := epub.UnzipAndParseEpub("./fixtures/drjekyllmrhyde_v2.epub")
	//book3, err := epub.ParseEpub("./fixtures/drjekyllmrhyde_v3.epub")

	if err != nil {
		log.Fatal(err)
	}

	var titles []string
	for _, c := range book2.Texts {
		titles = append(titles, c.Title)
	}
	fmt.Println(titles)
	fmt.Println(book2.Metadata.Title)
	//fmt.Println(book3.Metadata.Title)
	//epub.ParseEpub("../fixtures/drjekyllmrhyde_v2.epub")
}

```