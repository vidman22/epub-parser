package parser

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Container struct {
	Rootfile Rootfile `xml:"rootfiles>rootfile"`
}

type Rootfile struct {
	Path string `xml:"full-path,attr"`
	Type string `xml:"media-type,attr"`
}

type Book struct {
	Metadata  Metadata
	Manifest  Manifest
	Container Container
	Spine     Spine
	ZipReader *zip.ReadCloser
}

func (book *Book) ReadXML(fileName string, targetStruct interface{}) error {
	reader, err := book.open(fileName)
	if err != nil {
		return err
	}
	defer reader.Close()
	dec := xml.NewDecoder(reader)
	return dec.Decode(targetStruct)
}

func (book *Book) open(fileName string) (io.ReadCloser, error) {
	for _, file := range book.ZipReader.File {
		if file.Name == fileName {
			return file.Open()
		}
	}
	return nil, fmt.Errorf("file %s not exist", fileName)
}

func getTitles3(metaData []DefaultAttributes, metaMap map[string]map[string]Meta) *[]DefaultAttributes {
	titles := make([]DefaultAttributes, len(metaData))
	for i, title := range metaData {
		// fileAs := getMetadata(metaMap, title.Id, "file-as")
		// titleType := getMetadata(metaMap, title.Id, "title-type")
		titles[i] = DefaultAttributes{
			Text: title.Text,
			Id:   title.Id,
			Lang: title.Lang,
			// Language: title.Lang,
			// Type:     titleType,
			// FileAs:   fileAs,
		}
	}
	return &titles
}

func getCreators3(metaData []DefaultAttributes, metaMap map[string]map[string]Meta) *[]Creator {
	if metaData != nil {
		creators := make([]Creator, len(metaData))
		for i, creator := range metaData {
			fileAs := getMetadata(metaMap, creator.Id, "file-as")
			rawRole := getMetadata(metaMap, creator.Id, "role")
			role := getRole(metaMap, creator, rawRole)
			creators[i] = Creator{
				Name:     creator.Text,
				FileAs:   fileAs,
				RawRole:  rawRole,
				Language: creator.Lang,
				Role:     role,
			}
		}
		return &creators
	}
	return nil
}

func getMetadata(metaData map[string]map[string]Meta, id string, metaDataKey string) string {
	if id != "" {
		idData, idOk := metaData[id]
		if idOk {
			keyData, keyOk := idData[metaDataKey]
			if keyOk {
				return keyData.Text
			}
		}
	}
	return ""
}

func getMetadataSchema(metaData map[string]map[string]Meta, id string, metaDataKey string) string {
	if id != "" {
		idData, idOk := metaData[id]
		if idOk {
			keyData, keyOk := idData[metaDataKey]
			if keyOk {
				return keyData.Scheme
			}
		}
	}
	return ""
}

func getMetaMap(metaData []Meta) *map[string]map[string]Meta {
	metaMap := make(map[string]map[string]Meta)
	for _, meta := range metaData {
		if meta.Refines != "" && meta.Property != "" {
			id := strings.Replace(meta.Refines, "#", "", 1)
			innerMap, ok := metaMap[id]
			if !ok {
				innerMap = make(map[string]Meta)
				metaMap[id] = innerMap
			}
			innerMap[meta.Property] = meta
		}
	}
	return &metaMap
}

func getRole(metaMap map[string]map[string]Meta, creator DefaultAttributes, rawRole string) string {
	scheme := getMetadataSchema(metaMap, creator.Id, "role")
	role := "unknown"
	if scheme == "marc:relators" {
		role, _ = Relator[rawRole]
	}
	return role
}

func ParseOpf3(opfPath string, book *Book) error {
	opf := OPFPackage{}
	err := book.ReadXML(opfPath, &opf)
	if err != nil {
		return err
	}
	if opf.Metadata.Meta == nil || opf.Metadata.Identifier == nil {
		return errors.New("no metadata")
	}
	metaMap := getMetaMap(*opf.Metadata.Meta)

	identifiers := make([]ID, len(*opf.Metadata.Identifier))
	for i, identifier := range *opf.Metadata.Identifier {
		parts := strings.Split(identifier.Text, ":")
		if len(parts) == 2 {
			scheme := parts[0]
			id := parts[1]
			identifiers[i] = ID{
				Id:     id,
				Scheme: &scheme,
			}
		}
		// if identifier.Id == opf.UniqueIdentifier {
		// 	book.Metadata. = Identifier{
		// 		Id:     id,
		// 		Scheme: &scheme,
		// 	}
		// }
	}
	book.Metadata.Identifier = &identifiers

	if opf.Manifest != nil {
		book.Manifest = getManifest(*opf.Manifest)
	}

	if opf.Spine != nil {
		book.Spine = getSpine(*opf.Spine)
	}
	if opf.Metadata.Title != nil {
		book.Metadata.Title = getTitles3(*opf.Metadata.Title, *metaMap)
	}
	if opf.Metadata.Language != nil {
		book.Metadata.Language = getLanguages(*opf.Metadata.Language)
	}
	if opf.Metadata.Creator != nil {
		book.Metadata.Creator = getCreators(*opf.Metadata.Creator)
	}
	if opf.Metadata.Contributor != nil {
		book.Metadata.Contributor = getCreators(*opf.Metadata.Contributor)
	}
	if opf.Metadata.Publisher != nil {
		book.Metadata.Publisher = getDefaultAttributes(*opf.Metadata.Publisher)
	}
	if opf.Metadata.Subject != nil {
		book.Metadata.Subject = getDefaultAttributes(*opf.Metadata.Subject)
	}
	if opf.Metadata.Description != nil {
		book.Metadata.Description = getDefaultAttributes(*opf.Metadata.Description)
	}
	if opf.Metadata.Date != nil {
		book.Metadata.Date = getDate(*opf.Metadata.Date)
	}
	if metaMap != nil {
		book.Metadata.CoverId = getCoverId(*opf.Metadata.Meta)
	}
	return err
}

type Link struct {
	Href       string `xml:"href,attr"`
	Rel        string `xml:"rel,attr"`
	Id         string `xml:"id,attr,omitempty"`
	MediaType  string `xml:"media-type,attr,omitempty"`
	Properties string `xml:"properties,attr,omitempty"`
	Refines    string `xml:"refines,attr,omitempty"`
}
