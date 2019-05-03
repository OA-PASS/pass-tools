package model

const MetadataSchemaID = "https://oa-pass.github.io/metadata-schemas/jhu/global.json"

// Metadata encapsulates submission metadata
type Metadata struct {
	Schema           string            `json:"$schema"`
	Agreements       map[string]string `json:"agreements,omitempty"`
	Abstract         string            `json:"abstract,omitempty"`
	AgentInformation struct {
		Name    string `json:"name,omitempty"`
		Version string `json:"version,omitempty"`
	} `json:"agent_information"`
	Authors           []Author `json:"authors,omitempty"`
	DOI               string   `json:"doi,omitempty"`
	EmbargoEndDate    string   `json:"Embargo-end-date,omitempty"`
	Nlmta             string   `json:"journal-NLMTA-ID,omitempty"`
	JournalTitle      string   `json:"journal-title,omitempty"`
	JournalShortTitle string   `json:"journal-title-short,omitempty"`
	Issue             string   `json:"issue,omitempty"`
	ISSNs             []ISSN   `json:"issns,omitempty"`
	Publisher         string   `json:"publisher,omitempty"`
	PublicationDate   string   `json:"publicationDate,omitempty"`
	ArticleTitle      string   `json:"title,omitempty"`
	UnderEmbargo      string   `json:"under-embargo,omitempty"`
	Volume            string   `json:"volume,omitempty"`
}

type Author struct {
	Name  string `json:"author,omitempty"`
	Orcid string `json:"orcid,omitempty"`
}

type ISSN struct {
	ISSN    string `json:"issn,omitempty"`
	PubType string `json:"pubType,omitempty"`
}
