package epub

import (
	"archive/zip"
	"fmt"
	"os"

	"github.com/vidman22/epub-parser/internal"
)

func ParseEpub(path string) (*parser.ParsedBookResult, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, err
	}

	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read epub with zip")
	}
	defer r.Close()

	res, err := parser.OpenBook(r)

	if err != nil {
		return nil, err
	}

	return res, nil
}
