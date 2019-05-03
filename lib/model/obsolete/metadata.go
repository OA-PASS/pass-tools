package obsolete

type MetadataAgentInformation struct {
	Information struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"information"`
}

type MetadataCrossref struct {
	DOI               string `json:"doi"`
	Publisher         string `json:"publisher"`
	JournalShortTitle string `json:"journal-title-short"`
}

type MetadataCommon struct {
	ArticleTitle string `json:"title"`
	JournalTitle string `json:"journal-title"`
	Volume       string `json:"volume"`
	ISSN         string `json:"ISSN"`
	IssnMap      map[string]struct {
		PubType []string `json:"pub-type"`
	} `json:"issn-map"`
	Issue           string `json:"issue"`
	Publisher       string `json:"publisher"`
	PublicationDate string `json:"publicationDate"`
	Abstract        string `json:"abstract"`
	Authors         []struct {
		Name  string `json:"author"`
		Orcid string `json:"orcid"`
	} `json:"authors"`
	UnderEmbargo   string `json:"under-embargo"`
	EmbargoEndDate string `json:"Embargo-end-date"`
	Subjects       string `json:"subjects"`
}

type MetadataJ10p struct {
	Authors []struct {
		Name string `json:"author"`
	} `json:"authors"`
	AgreementText string `json:"embargo"`
	Agreed        string `json:"agreement-to-deposit"`
}

type MetadataPMC struct {
	Nlmta string `json:"nlmta"`
}
