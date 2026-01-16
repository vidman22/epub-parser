package parser

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type OPFHeaderDetails struct {
	Version          string `xml:"version,attr"`
	UniqueIdentifier string `xml:"unique-identifier,attr"`
	ID               string `xml:"id,attr,omitempty"`
	Prefix           string `xml:"prefix,attr,omitempty"`
	Lang             string `xml:"lang,attr,omitempty"`
	Dir              string `xml:"dir,attr,omitempty"`
}

func getLikelyTOC(manifestItems *[]Item, navDir string) (likelyTocPathV2 string, likelyTocPathV3 string) {

	if manifestItems != nil {
		items := *manifestItems
		for _, v := range items {
			if filepath.Ext(v.Href) == ".ncx" {
				likelyTocPathV2 = filepath.Join(navDir, v.Href)
			}
			// this is to handle the variations of .xhtml like htm and html
			if (strings.Contains(v.Href, "toc") || strings.Contains(v.Properties, "nav")) && filepath.Ext(v.Href) != ".ncx" {
				likelyTocPathV3 = filepath.Join(navDir, v.Href)
			}
		}
	}
	// if for some reason version 3 parser doesn't have a v3 toc but has the old version .ncx file, use the .ncx file
	if likelyTocPathV3 == "" && likelyTocPathV2 != "" {
		likelyTocPathV3 = likelyTocPathV2
	}
	return likelyTocPathV2, likelyTocPathV3
}

// OpenBook will open epub2 and epub3 files toc.ncx is epub2 toc.xhtml is epub3
func OpenBook(reader *zip.ReadCloser) (*ParsedBookResult, error) {
	book := &Book{ZipReader: reader}
	err := book.ReadXML("META-INF/container.xml", &book.Container)
	if err != nil {
		return nil, err
	}
	header := OPFHeaderDetails{}
	err = book.ReadXML(book.Container.Rootfile.Path, &header)
	if err != nil {
		return nil, err
	}
	ebookVersion, err := strconv.ParseFloat(header.Version, 64)
	if err != nil {
		return nil, err
	}

	if book.Container.Rootfile.Path == "" {
		return nil, fmt.Errorf("parser poorly formatted")
	}

	rootDir := filepath.Dir(book.Container.Rootfile.Path)

	tocMap := make(map[string]string)

	switch {
	case ebookVersion >= 3.0 && ebookVersion < 4.0:
		err := ParseOpf3(book.Container.Rootfile.Path, book)
		if err != nil {
			return nil, err
		}
		_, likelyTocPathV3 := getLikelyTOC(book.Manifest.Item, rootDir)
		fBytes, err := readZipFile(reader, likelyTocPathV3)
		if err != nil {
			return nil, fmt.Errorf("failed to read zip file %s: %w", likelyTocPathV3, err)
		}
		tocMap, err = ParseNavDoc(fBytes, rootDir)
		if err != nil {
			return nil, fmt.Errorf("failed to parse nav doc: %w", err)
		}

	case ebookVersion >= 2.0 && ebookVersion < 3.0:
		err := ParseOpf(book.Container.Rootfile.Path, book)
		if err != nil {
			return nil, err
		}
		likelyTocPathV2, _ := getLikelyTOC(book.Manifest.Item, rootDir)
		fBytes, err := readZipFile(reader, likelyTocPathV2)
		if err != nil {
			return nil, fmt.Errorf("failed to read zip file %s: %w", likelyTocPathV2, err)
		}
		tocMap, err = ParseNcx(fBytes, rootDir)

		if err != nil {
			return nil, fmt.Errorf("failed to parse NCX: %w", err)
		}

	default:
		return nil, fmt.Errorf("%f is not a supported version", ebookVersion)
	}

	res, cover, err := processEpubContent(Params{
		rootDir:       rootDir,
		manifestItems: *book.Manifest.Item,
		spineItemRefs: book.Spine.Itemrefs,
		tocMap:        tocMap,
		r:             reader,
	})

	if err != nil {
		return nil, err
	}

	md := book.Metadata

	resMetadata := ResultMetadata{
		MainId: func() string {
			if md.Identifier != nil && len(*md.Identifier) > 0 {
				return (*md.Identifier)[0].Id
			}
			return ""
		}(),
		Title: func() string {
			if md.Title != nil && len(*md.Title) > 0 {
				return (*md.Title)[0].Text
			}
			return ""
		}(),
		Identifier: func() string {
			if md.Identifier != nil && len(*md.Identifier) > 0 {
				return (*md.Identifier)[0].Id
			}
			return ""
		}(),
		Language: func() string {
			if md.Language != nil && len(*md.Language) > 0 {
				return (*md.Language)[0].Text
			}
			return ""
		}(),
		Creator: func() string {
			if md.Creator != nil && len(*md.Creator) > 0 {
				return (*md.Creator)[0].Name
			}
			return ""
		}(),
		Contributor: func() string {
			if md.Contributor != nil && len(*md.Contributor) > 0 {
				return (*md.Contributor)[0].Name
			}
			return ""
		}(),
		Cover: cover,
		Publisher: func() string {
			if md.Publisher != nil && len(*md.Publisher) > 0 {
				return (*md.Publisher)[0].Text
			}
			return ""
		}(),
		Subject: func() string {
			if md.Subject != nil && len(*md.Subject) > 0 {
				return (*md.Subject)[0].Text
			}
			return ""
		}(),
		Description: func() string {
			if md.Description != nil && len(*md.Description) > 0 {
				return (*md.Description)[0].Text
			}
			return ""
		}(),
		Date: func() string {
			if md.Date != nil && len(*md.Date) > 0 {
				return (*md.Date)[0].Text
			}
			return ""
		}(),
	}

	return &ParsedBookResult{
		Metadata: &resMetadata,
		Texts:    res,
	}, nil
}
