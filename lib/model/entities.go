package model

type Grant struct {
	ID string `json:"@id"`
	PI string `json:"pi,omitempty"`
}

type Submission struct {
	ID        string `json:"@id,omitempty"`
	Submitter string `json:"submitter,omitempty"`
	Metadata  string `json:"metadata,omitempty"`
}
