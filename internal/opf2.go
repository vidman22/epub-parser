package parser

import (
	"encoding/xml"
)

func getManifest(metaData Manifest) Manifest {
	refs := make([]Item, len(*metaData.Item))

	for i, m := range *metaData.Item {
		refs[i] = Item{
			Id:        m.Id,
			Href:      m.Href,
			MediaType: m.MediaType,
			// Fallback:     m.Fallback,
			// MediaOverlay: "",
			// Properties:   "",
		}
	}
	return Manifest{
		Item: &refs,
		Id:   metaData.Id,
	}
}

func getSpine(metaData Spine) Spine {
	refs := make([]Itemref, len(metaData.Itemrefs))

	for i, ir := range metaData.Itemrefs {
		refs[i] = Itemref{
			Idref: ir.Idref,
		}
	}
	return Spine{
		Toc:      metaData.Toc,
		Itemrefs: refs,
	}
}

func getCoverId(metaMap []Meta) string {
	for _, meta := range metaMap {
		if meta.Name == "cover" && meta.Content != "" {
			return meta.Content
		}
	}
	return ""
}

func getTitles(metaData []DefaultAttributes) *[]DefaultAttributes {
	titles := make([]DefaultAttributes, len(metaData))
	for i, title := range metaData {
		titles[i] = DefaultAttributes{
			Text: title.Text,
			Lang: title.Lang,
		}
	}
	return &titles
}

func getLanguages(metaData []ID) *[]ID {
	languages := make([]ID, len(metaData))
	for i, language := range metaData {
		languages[i] = ID{
			Text: language.Text,
			Id:   language.Id,
		}
		// language.Text
	}
	return &languages
}

func getCreators(metaData []Creator) *[]Creator {
	if metaData != nil {
		creators := make([]Creator, len(metaData))
		for i, creator := range metaData {
			role, ok := Relator[creator.Role]
			if !ok && creator.Role != "" {
				role = "unknown"
			}
			creators[i] = Creator{
				Name:     creator.Text,
				FileAs:   creator.FileAs,
				RawRole:  creator.Role,
				Language: creator.Lang,
				Role:     role,
			}
		}
		return &creators
	}
	return nil
}

func getDefaultAttributes(metaData []DefaultAttributes) *[]DefaultAttributes {
	if metaData != nil {
		defaultAttributes := make([]DefaultAttributes, len(metaData))
		for i, defaultAttribute := range metaData {
			defaultAttributes[i] = DefaultAttributes{
				Text: defaultAttribute.Text,
				Lang: defaultAttribute.Lang,
			}
		}
		return &defaultAttributes
	}
	return nil
}

func getDate(metaData []Date) *[]Date {
	if metaData != nil {
		dates := make([]Date, len(metaData))
		for i, date := range metaData {
			dates[i] = date
		}
		return &dates
	}
	return nil
}

func ParseOpf(opfFilePath string, book *Book) error {
	opf := OPFPackage{}
	err := book.ReadXML(opfFilePath, &opf)
	if err != nil {
		return err
	}

	identifiers := make([]ID, len(*opf.Metadata.Identifier))
	for i, identifier := range *opf.Metadata.Identifier {
		identifiers[i] = ID{
			Id:     identifier.Text,
			Scheme: identifier.Scheme,
		}
		// if identifier.Id == opf.ID {
		// 	book.Metadata.MainId = Identifier{
		// 		Id:     identifier.Text,
		// 		Scheme: identifier.Scheme,
		// 	}
		// }
	}

	book.Metadata.Identifier = &identifiers

	if opf.Spine != nil {
		book.Spine = getSpine(*opf.Spine)
	}

	if opf.Manifest != nil {
		book.Manifest = getManifest(*opf.Manifest)
	}

	if opf.Metadata.Title != nil {
		book.Metadata.Title = getTitles(*opf.Metadata.Title)
	}
	if opf.Metadata != nil {
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
	if opf.Metadata != nil {
		book.Metadata.CoverId = getCoverId(*opf.Metadata.Meta)
	}

	return err
}

type OPFPackage struct {
	XMLName          xml.Name  `xml:"package"`
	Metadata         *Metadata `xml:"metadata"`
	Manifest         *Manifest `xml:"manifest"`
	Spine            *Spine    `xml:"spine"`
	Version          string    `xml:"version,attr"`
	UniqueIdentifier string    `xml:"unique-identifier,attr"`
	ID               string    `xml:"id,attr,omitempty"`
	Prefix           string    `xml:"prefix,attr,omitempty"`
	Lang             string    `xml:"lang,attr,omitempty"`
	Dir              string    `xml:"dir,attr,omitempty"`
}

type Identifier struct {
	Id     string
	Scheme *string
}

type Metadata struct {
	Title       *[]DefaultAttributes `xml:"title,dc:title"`
	CoverId     string               `xml:"coverId,omitempty"`
	Identifier  *[]ID                `xml:"identifier,dc:identifier"`
	Language    *[]ID                `xml:"language,dc:language"`
	Creator     *[]Creator           `xml:"creator,dc:creator,omitempty"`
	Contributor *[]Creator           `xml:"contributor,dc:contributor,omitempty"`
	Publisher   *[]DefaultAttributes `xml:"publisher,omitempty"`
	Subject     *[]DefaultAttributes `xml:"subject,dc:subject,omitempty"`
	Description *[]DefaultAttributes `xml:"description,dc:description,omitempty"`
	Date        *[]Date              `xml:"date,dc:date,omitempty"`
	Type        *[]ID                `xml:"type,omitempty"`
	Format      *[]ID                `xml:"format,omitempty"`
	Source      *[]DefaultAttributes `xml:"source,omitempty"`
	Relation    *[]DefaultAttributes `xml:"relation,omitempty"`
	Coverage    *[]DefaultAttributes `xml:"coverage,omitempty"`
	Rights      *[]DefaultAttributes `xml:"rights,omitempty"`
	Meta        *[]Meta              `xml:"meta,omitempty"`
}

type Spine struct {
	Toc      string    `xml:"toc,attr"`
	Itemrefs []Itemref `xml:"itemref"`
}

type Itemref struct {
	Idref string `xml:"idref,attr"`
}
type Creator struct {
	Text     string `xml:",chardata"`
	FileAs   string `xml:"file-as,attr,omitempty"`
	Id       string `xml:"id,attr,omitempty"`
	Lang     string `xml:"lang,attr,omitempty"`
	Role     string `xml:"role,attr,omitempty"`
	Name     string
	Language string
	RawRole  string
}

type DefaultAttributes struct {
	Text string `xml:",chardata"`
	Id   string `xml:"id,attr,omitempty"`
	Lang string `xml:"lang,attr,omitempty"`
}

type Date struct {
	Text  string `xml:",chardata"`
	Event string `xml:"event,attr,omitempty"`
	Id    string `xml:"id,attr,omitempty"`
}

type ID struct {
	Text   string  `xml:",chardata"`
	Id     string  `xml:"id,attr,omitempty"`
	Scheme *string `xml:"scheme,attr,omitempty"`
}

type Meta struct {
	Text     string `xml:",chardata"`
	Lang     string `xml:"lang,attr,omitempty"`
	Scheme   string `xml:"scheme,attr,omitempty"`
	Id       string `xml:"id,attr,omitempty"`
	Dir      string `xml:"dir,attr,omitempty"`
	Property string `xml:"property,attr,omitempty"` //omitempty because the deprecated meta has no property
	Refines  string `xml:"refines,attr,omitempty"`
	Name     string `xml:"name,attr,omitempty"`    //deprecated in 3
	Content  string `xml:"content,attr,omitempty"` //deprecated in 3
}

type Manifest struct {
	Id   string  `xml:"id,attr,omitempty"`
	Item *[]Item `xml:"item"`
}

type Item struct {
	Id                string  `xml:"id,attr"`
	Href              string  `xml:"href,attr"`
	MediaType         string  `xml:"media-type,attr"`
	Properties        string  `xml:"properties,attr"`
	Fallback          *string `xml:"fallback,attr,omitempty"`
	FallbackStyle     *string `xml:"fallback-style,attr,omitempty"`
	RequiredModules   *string `xml:"required-modules,attr,omitempty"`
	RequiredNamespace *string `xml:"required-namespace,attr,omitempty"`
}

type Title struct {
	Title    string
	Language string
	Type     string
	FileAs   string
}

var Relator = map[string]string{
	"abr":  "abridger",
	"acp":  "art copyist",
	"act":  "actor",
	"adi":  "art director",
	"adp":  "adapter",
	"aft":  "author of afterword, colophon, etc.",
	"anc":  "announcer",
	"anl":  "analyst",
	"anm":  "animator",
	"ann":  "annotator",
	"ant":  "bibliographic antecedent",
	"ape":  "appellee",
	"apl":  "appellant",
	"app":  "applicant",
	"aqt":  "author in quotations or text abstracts",
	"arc":  "architect",
	"ard":  "artistic director",
	"arr":  "arranger",
	"art":  "artist",
	"asg":  "assignee",
	"asn":  "associated name",
	"ato":  "autographer",
	"att":  "attributed name",
	"auc":  "auctioneer",
	"aud":  "author of dialog",
	"aue":  "audio engineer",
	"aui":  "author of introduction, etc.",
	"aup":  "audio producer",
	"aus":  "screenwriter",
	"aut":  "author",
	"bdd":  "binding designer",
	"bjd":  "bookjacket designer",
	"bka":  "book artist",
	"bkd":  "book designer",
	"bkp":  "book producer",
	"blw":  "blurb writer",
	"bnd":  "binder",
	"bpd":  "bookplate designer",
	"brd":  "broadcaster",
	"brl":  "braille embosser",
	"bsl":  "bookseller",
	"cad":  "casting director",
	"cas":  "caster",
	"ccp":  "conceptor",
	"chrc": "choreographer",
	"-clb": "collaborator",
	"cli":  "client",
	"cll":  "calligrapher",
	"clr":  "colorist",
	"clt":  "collotyper",
	"cmm":  "commentator",
	"cmp":  "composer",
	"cmt":  "compositor",
	"cnd":  "conductor",
	"cng":  "cinematographer",
	"cns":  "censor",
	"coe":  "contestant-appellee",
	"col":  "collector",
	"com":  "compiler",
	"con":  "conservator",
	"cop":  "camera operator",
	"cor":  "collection registrar",
	"cos":  "contestant",
	"cot":  "contestant-appellant",
	"cou":  "court governed",
	"cov":  "cover designer",
	"cpc":  "copyright claimant",
	"cpe":  "complainant-appellee",
	"cph":  "copyright holder",
	"cpl":  "complainant",
	"cpt":  "complainant-appellant",
	"cre":  "creator",
	"crp":  "correspondent",
	"crr":  "corrector",
	"crt":  "court reporter",
	"csl":  "consultant",
	"csp":  "consultant to a project",
	"cst":  "costume designer",
	"ctb":  "contributor",
	"cte":  "contestee-appellee",
	"ctg":  "cartographer",
	"ctr":  "contractor",
	"cts":  "contestee",
	"ctt":  "contestee-appellant",
	"cur":  "curator",
	"cwt":  "commentator for written text",
	"dbd":  "dubbing director",
	"dbp":  "distribution place",
	"dfd":  "defendant",
	"dfe":  "defendant-appellee",
	"dft":  "defendant-appellant",
	"dgc":  "degree committee member",
	"dgg":  "degree granting institution",
	"dgs":  "degree supervisor",
	"dis":  "dissertant",
	"djo":  "dj",
	"dln":  "delineator",
	"dnc":  "dancer",
	"dnr":  "donor",
	"dpc":  "depicted",
	"dpt":  "depositor",
	"drm":  "draftsman",
	"drt":  "director",
	"dsr":  "designer",
	"dst":  "distributor",
	"dtc":  "data contributor",
	"dte":  "dedicatee",
	"dtm":  "data manager",
	"dto":  "dedicator",
	"dub":  "dubious author",
	"edc":  "editor of compilation",
	"edd":  "editorial director",
	"edm":  "editor of moving image work",
	"edt":  "editor",
	"egr":  "engraver",
	"elg":  "electrician",
	"elt":  "electrotyper",
	"eng":  "engineer",
	"enj":  "enacting jurisdiction",
	"etr":  "etcher",
	"evp":  "event place",
	"exp":  "expert",
	"fac":  "facsimilist",
	"fds":  "film distributor",
	"fld":  "field director",
	"flm":  "film editor",
	"fmd":  "film director",
	"fmk":  "filmmaker",
	"fmo":  "former owner",
	"fmp":  "film producer",
	"fnd":  "funder",
	"fon":  "founder",
	"fpy":  "first party",
	"frg":  "forger",
	"gdv":  "game developer",
	"gis":  "geographic information specialist",
	"-grt": "graphic technician",
	"his":  "host institution",
	"hnr":  "honoree",
	"hst":  "host",
	"ill":  "illustrator",
	"ilu":  "illuminator",
	"ins":  "inscriber",
	"inv":  "inventor",
	"isb":  "issuing body",
	"itr":  "instrumentalist",
	"ive":  "interviewee",
	"ivr":  "interviewer",
	"jud":  "judge",
	"jug":  "jurisdiction governed",
	"lbr":  "laboratory",
	"lbt":  "librettist",
	"ldr":  "laboratory director",
	"led":  "lead",
	"lee":  "libelee-appellee",
	"lel":  "libelee",
	"len":  "lender",
	"let":  "libelee-appellant",
	"lgd":  "lighting designer",
	"lie":  "libelant-appellee",
	"lil":  "libelant",
	"lit":  "libelant-appellant",
	"lsa":  "landscape architect",
	"lse":  "licensee",
	"lso":  "licensor",
	"ltg":  "lithographer",
	"ltr":  "letterer",
	"lyr":  "lyricist",
	"mcp":  "music copyist",
	"mdc":  "metadata contact",
	"med":  "medium",
	"mfp":  "manufacture place",
	"mfr":  "manufacturer",
	"mka":  "makeup artist",
	"mod":  "moderator",
	"mon":  "monitor",
	"mrb":  "marbler",
	"mrk":  "markup editor",
	"msd":  "musical director",
	"mte":  "metal-engraver",
	"mtk":  "minute taker",
	"mup":  "music programmer",
	"mus":  "musician",
	"mxe":  "mixing engineer",
	"nan":  "news anchor",
	"nrt":  "narrator",
	"onp":  "onscreen participant",
	"opn":  "opponent",
	"org":  "originator",
	"orm":  "organizer",
	"osp":  "onscreen presenter",
	"oth":  "other",
	"own":  "owner",
	"pad":  "place of address",
	"pan":  "panelist",
	"pat":  "patron",
	"pbd":  "publishing director",
	"pbl":  "publisher",
	"pdr":  "project director",
	"pfr":  "proofreader",
	"pht":  "photographer",
	"plt":  "platemaker",
	"pma":  "permitting agency",
	"pmn":  "production manager",
	"pop":  "printer of plates",
	"ppm":  "papermaker",
	"ppt":  "puppeteer",
	"pra":  "praeses",
	"prc":  "process contact",
	"prd":  "production personnel",
	"pre":  "presenter",
	"prf":  "performer",
	"prg":  "programmer",
	"prm":  "printmaker",
	"prn":  "production company",
	"pro":  "producer",
	"prp":  "production place",
	"prs":  "production designer",
	"prt":  "printer",
	"prv":  "provider",
	"pta":  "patent applicant",
	"pte":  "plaintiff-appellee",
	"ptf":  "plaintiff",
	"pth":  "patent holder",
	"ptt":  "plaintiff-appellant",
	"pup":  "publication place",
	"rap":  "rapporteur",
	"rbr":  "rubricator",
	"rcd":  "recordist",
	"rce":  "recording engineer",
	"rcp":  "addressee",
	"rdd":  "radio director",
	"red":  "redaktor",
	"ren":  "renderer",
	"res":  "researcher",
	"rev":  "reviewer",
	"rpc":  "radio producer",
	"rps":  "repository",
	"rpt":  "reporter",
	"rpy":  "responsible party",
	"rse":  "respondent-appellee",
	"rsg":  "restager",
	"rsp":  "respondent",
	"rsr":  "restorationist",
	"rst":  "respondent-appellant",
	"rth":  "research team head",
	"rtm":  "research team member",
	"rxa":  "remix artist",
	"sad":  "scientific advisor",
	"sce":  "scenarist",
	"scl":  "sculptor",
	"scr":  "scribe",
	"sde":  "sound engineer",
	"sds":  "sound designer",
	"sec":  "secretary",
	"sfx":  "special effects provider",
	"sgd":  "stage director",
	"sgn":  "signer",
	"sht":  "supporting host",
	"sll":  "seller",
	"sng":  "singer",
	"spk":  "speaker",
	"spn":  "sponsor",
	"spy":  "second party",
	"srv":  "surveyor",
	"std":  "set designer",
	"stg":  "setting",
	"stl":  "storyteller",
	"stm":  "stage manager",
	"stn":  "standards body",
	"str":  "stereotyper",
	"swd":  "software developer",
	"tad":  "technical advisor",
	"tau":  "television writer",
	"tcd":  "technical director",
	"tch":  "teacher",
	"ths":  "thesis advisor",
	"tld":  "television director",
	"tlg":  "television guest",
	"tlh":  "television host",
	"tlp":  "television producer",
	"trc":  "transcriber",
	"trl":  "translator",
	"tyd":  "type designer",
	"tyg":  "typographer",
	"uvp":  "university place",
	"vac":  "voice actor",
	"vdg":  "videographer",
	"vfx":  "visual effects provider",
	"-voc": "vocalist",
	"wac":  "writer of added commentary",
	"wal":  "writer of added lyrics",
	"wam":  "writer of accompanying material",
	"wat":  "writer of added text",
	"wdc":  "woodcutter",
	"wde":  "wood engraver",
	"wfs":  "writer of film story",
	"wft":  "writer of intertitles",
	"win":  "writer of introduction",
	"wit":  "witness",
	"wpr":  "writer of preface",
	"wst":  "writer of supplementary textual content",
	"wts":  "writer of television story",
}
