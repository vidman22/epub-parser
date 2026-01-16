package epub

import (
	"testing"

	parser "github.com/vidman22/epub-parser/internal"
)

func Test_parse_epub_2_0_opf(t *testing.T) {

	book, err := ParseEpub("./fixtures/drjekyllmrhyde_v2.epub")
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	assertMetadata(t, book.Metadata)

	var titles []string
	for _, c := range book.Texts {
		titles = append(titles, c.Title)
	}

	assertv2Titles(t, titles)
}

func Test_parse_epub_3_0_opf(t *testing.T) {

	book, err := ParseEpub("./fixtures/drjekyllmrhyde_v3.epub")

	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	assertMetadata(t, book.Metadata)

	var titles []string
	for _, c := range book.Texts {
		titles = append(titles, c.Title)
	}

	assertv3Titles(t, titles)
}

func assertv2Titles(t *testing.T, titles []string) {
	expectedTitles := []string{
		"The Strange Case Of Dr. Jekyll And Mr. Hyde",
		"STORY OF THE DOOR",
		"SEARCH FOR MR. HYDE",
		"DR. JEKYLL WAS QUITE AT EASE",
		"THE CAREW MURDER CASE",
		"INCIDENT OF THE LETTER",
		"INCIDENT OF DR. LANYON",
		"INCIDENT AT THE WINDOW",
		"THE LAST NIGHT",
		"DR. LANYON’S NARRATIVE",
		"HENRY JEKYLL’S FULL STATEMENT OF THE CASE",
		"THE FULL PROJECT GUTENBERG LICENSE",
	}

	if len(titles) != len(expectedTitles) {
		t.Logf("titles length expected %d but is %d", len(expectedTitles), len(titles))
		t.Fail()
		return
	}

	for i, expected := range expectedTitles {
		assertEquals("title["+string(rune(i+'0'))+"]", t, titles[i], expected)
	}
}

func assertv3Titles(t *testing.T, titles []string) {
	expectedTitles := []string{
		"Contents",
		"STORY OF THE DOOR",
		"SEARCH FOR MR. HYDE",
		"DR. JEKYLL WAS QUITE AT EASE",
		"THE CAREW MURDER CASE",
		"INCIDENT OF THE LETTER",
		"INCIDENT OF DR. LANYON",
		"INCIDENT AT THE WINDOW",
		"THE LAST NIGHT",
		"DR. LANYON’S NARRATIVE",
		"HENRY JEKYLL’S FULL STATEMENT OF THE CASE",
		"THE FULL PROJECT GUTENBERG LICENSE",
	}

	if len(titles) != len(expectedTitles) {
		t.Logf("titles length expected %d but is %d", len(expectedTitles), len(titles))
		t.Fail()
		return
	}

	for i, expected := range expectedTitles {
		assertEquals("title["+string(rune(i+'0'))+"]", t, titles[i], expected)
	}
}

func assertMetadata(t *testing.T, metaData *parser.ResultMetadata) {
	assertEquals("mainId", t, metaData.MainId, "//www.gutenberg.org/43", "http://www.gutenberg.org/43")
	assertEquals("title", t, metaData.Title, "The Strange Case of Dr. Jekyll and Mr. Hyde")
	assertEquals("identifier", t, metaData.Identifier, "//www.gutenberg.org/43", "http://www.gutenberg.org/43")
	assertEquals("language", t, metaData.Language, "en")
	assertEquals("creator", t, metaData.Creator, "Robert Louis Stevenson")
	assertEquals("contributor", t, metaData.Contributor, "")
	assertEquals("publisher", t, metaData.Publisher, "")
	assertEquals("subject", t, metaData.Subject, "Science fiction")
	assertEquals("description", t, metaData.Description, "")
	assertEquals("date", t, metaData.Date, "2008-06-27")
}

func assertEquals(fieldName string, t *testing.T, actuallyValue string, expectedValue string, expectedValue2 ...string) {
	if actuallyValue != expectedValue {
		if len(expectedValue2) == 0 || actuallyValue != expectedValue2[0] {
			if len(expectedValue2) == 0 {
				t.Logf("'%s' expected '%s' but is '%s'", fieldName, expectedValue, actuallyValue)
			} else {
				t.Logf("'%s' expected '%s' or '%s' but is '%s'", fieldName, expectedValue, expectedValue2[0], actuallyValue)
			}
			t.Fail()
		}
	}
}
