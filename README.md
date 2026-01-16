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
git tag -a v0.1.3 -m "version 0.1.3"
git push origin v0.1.3

git tag -a v0.1.3 -m "version 0.1.3"
git push origin v0.1.3


## License
This project is licensed under the MIT License. See the LICENSE file for details.
