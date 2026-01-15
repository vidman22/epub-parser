package parser

type ResultMetadata struct {
	MainId      string
	Title       string
	Identifier  string
	Language    string
	Creator     string
	Contributor string
	Publisher   string
	Subject     string
	Description string
	Date        string
	CoverPath   string
}

type Cover struct {
	FileName  string
	CoverPath string
	Mimetype  string
}

type Content struct {
	Html  string
	Title string
}

type ParsedBookResult struct {
	Metadata *ResultMetadata
	Texts    []Content
}

type DatabaseBook struct {
	ID                *int   `json:"id"`
	Title             string `json:"title"`
	Subtitle          string `json:"subtitle"`
	Author            string `json:"author"`
	PublicationYear   *int   `json:"publicationYear,omitempty"`
	CoverImageKey     string `json:"coverImageKey"`
	Isbn              string `json:"isbn"`
	Genre             string `json:"genre"`
	Publisher         string `json:"publisher"`
	LanguageID        int    `json:"languageID"`
	LastActiveChatper *int   `json:"lastActiveChapter"`
	LearnerLanguageID int    `json:"learnerLanguageID"`
	CefrLevel         string `json:"cefrLevel"`
	PageCount         int    `json:"pageCount"`
	CreatedBy         int    `json:"createdBy"`
	Description       string `json:"description"`
	SubjectID         int    `json:"subjectID"`
}

type DatabaseBookChapter struct {
	ID             *int    `json:"id"`
	BookID         int     `json:"book_id"`
	Order          int     `json:"order"`
	Title          string  `json:"title"`
	ReadingPath    string  `json:"readingPath"`
	AudioKey       *string `json:"audioKey"`
	PageStart      int     `json:"pageStart"`
	PageEnd        int     `json:"pageEnd"`
	TranscriptKey  string  `json:"transcriptKey"`
	TranslationKey *string `json:"translationKey"`
	WordCount      *int    `json:"wordCount"`
	CreatedBy      int     `json:"createdBy"`
}

type DatabaseBookWithChapters struct {
	ID                *int                  `json:"id"`
	Title             string                `json:"title"`
	Subtitle          string                `json:"subtitle"`
	Author            string                `json:"author"`
	PublicationYear   *int                  `json:"publicationYear,omitempty"`
	LastActiveChapter *int                  `json:"lastActiveChapter,omitempty"`
	CoverImageKey     string                `json:"coverImageKey"`
	Isbn              string                `json:"isbn"`
	Genre             string                `json:"genre"`
	Publisher         string                `json:"publisher"`
	LanguageID        int                   `json:"languageId"`
	LearnerLanguageID int                   `json:"learnerLanguageId"`
	CerfLevel         string                `json:"cerfLevel"`
	PageCount         int                   `json:"pageCount"`
	CreatedBy         int                   `json:"createdBy"`
	Description       string                `json:"description"`
	SubjectID         int                   `json:"subjectId"`
	BookChapters      []DatabaseBookChapter `json:"bookChapters"`
}
