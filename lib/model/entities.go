package model

type Grant struct {
	Context string `json:"@context,omitempty"`
	ID      string `json:"@id"`
	PI      string `json:"pi,omitempty"`
}

type Submission struct {
	Context   string `json:"@context,omitempty"`
	ID        string `json:"@id"`
	Submitter string `json:"submitter,omitempty"`
	Metadata  string `json:"metadata,omitempty"`
}
