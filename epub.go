package epub

import (
	"archive/zip"
	"fmt"
	"os"

	"github.com/vidman22/epub-parser/parser"
)

func ParseEpub(path string) (*parser.ParsedBookResult, error) {
	fmt.Printf("Attempting to open: %s\n", path)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("Error: File does not exist.")
		return nil, err
	}

	fmt.Println("Success: File found!")

	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read epub with zip")
	}
	defer r.Close()

	res, err := parser.OpenBook(r)

	if err != nil {
		return nil, err
	}

	return res, nil
}
