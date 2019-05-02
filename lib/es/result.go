package es

import (
	"strings"

	"github.com/oa-pass/pass-tools/lib/model"
)

const fcrepoPrefix = "/fcrepo/rest"

// IDResults encapsulates an elasticsearch results containing matching entity IDs.
type IDResults struct {
	Hits struct {
		Total int `json:"total"`
		Hit   []struct {
			Source struct {
				ID string `json:"@id"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// GrantResults encapsulates an elasticsearch results containing matching grants
type GrantResults struct {
	Hits struct {
		Total int `json:"total"`
		Hit   []struct {
			Source model.Grant `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func RelativeURI(uri string) string {
	idx := strings.Index(uri, "/fcrepo/rest")
	if idx < 1 {
		return uri
	}

	return uri[idx+len(fcrepoPrefix):]
}
