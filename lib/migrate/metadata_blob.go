package migrate

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/OA-PASS/metadata-schemas/lib/jsonschema"
	"github.com/OA-PASS/metadata-schemas/lib/schemas"
	"github.com/oa-pass/pass-tools/lib/client"
	"github.com/oa-pass/pass-tools/lib/es"
	"github.com/oa-pass/pass-tools/lib/model"
	"github.com/oa-pass/pass-tools/lib/model/obsolete"
	"github.com/pkg/errors"
)

type MetadataV0toV1 struct {
	DryRun  bool
	BaseURI client.BaseURI
	Fedora  client.Performer
	Elastic client.Performer
}

func (m MetadataV0toV1) Perform() error {
	submissions, err := m.find()
	if err != nil {
		return errors.Wrap(err, "could not find submissions")
	}

	for _, sub := range submissions {
		var origMetadata []map[string]interface{}
		err := m.extract(sub, &origMetadata)
		if err != nil {
			log.Printf("ERROR: Could not read metadata blob of %s: %s", sub.ID, err)
		}

		transformedMetadata, err := m.Transform(origMetadata)
		if err != nil {
			log.Printf("ERROR: Could not transform metadata blob of %s: %s", sub.ID, err)
			continue
		}

		serializedMetadata, err := json.Marshal(transformedMetadata)
		if err != nil {
			log.Printf("ERROR: Could not encode translated metadata blob: %s", err)
		}
		sub.Metadata = string(serializedMetadata)

		global, err := schemas.Load("jhu/global.json")
		if err != nil {
			return errors.Wrap(err, "Could not load schema")
		}

		err = jsonschema.NewValidator(global).Validate(serializedMetadata)
		if err != nil {
			log.Printf("ERROR schema invalid %+s", sub.Metadata)
		}

		if m.DryRun {
			log.Printf("Would have written %s", sub.Metadata)
			continue
		}

		err = m.Fedora.Perform(http.MethodPatch, m.BaseURI.Rebase(sub.ID), &client.Body{
			Content: sub, Type: client.ContentTypeJSONMerge}, nil)
		if err != nil {
			log.Printf("ERROR: could not update submission %s: %s", sub.ID, err)
		}
	}

	return nil
}

func (m MetadataV0toV1) extract(s model.Submission, metadataPtr interface{}) error {
	decoder := json.NewDecoder(strings.NewReader(s.Metadata))
	decoder.DisallowUnknownFields()
	return errors.Wrap(decoder.Decode(metadataPtr), "could not decode submission metadata")
}

func (m MetadataV0toV1) find() ([]model.Submission, error) {
	var results es.SubmissionResults
	err := m.Elastic.Perform(http.MethodPost, "", &client.Body{Content: `{
		"size": 500,
		"query": {
			"bool": {
				"must": {
					"match": {
						"@type": "Submission"
					}
				},
				"filter": {
					"exists": {
						"field": "metadata"
					}
				}
			}
		}
	}`, Type: client.ContentTypeJSON}, &results)
	if err != nil {
		return nil, errors.Wrap(err, "elestic query failed")
	}

	submissions := make([]model.Submission, 0, results.Hits.Total)
	for _, hit := range results.Hits.Hit {
		submissions = append(submissions, hit.Source)
	}

	return submissions, nil
}

func (m MetadataV0toV1) Transform(metadata []map[string]interface{}) (*model.Metadata, error) {
	var translated model.Metadata
	translated.Schema = model.MetadataSchemaID

	for _, member := range metadata {
		var err error
		label := member["id"]
		data := member["data"]

		id, _ := label.(string)

		switch id {
		case "agent_information":
			var old obsolete.MetadataAgentInformation
			err = m.parseBlock(id, data, &old)

			translated.AgentInformation.Name = old.Information.Name
			translated.AgentInformation.Version = old.Information.Version

		case "pmc":
			var old obsolete.MetadataPMC
			err = m.parseBlock(id, data, &old)

			translated.Nlmta = old.Nlmta

		case "crossref":
			var old obsolete.MetadataCrossref
			err = m.parseBlock(id, data, &old)

			translated.DOI = old.DOI
			translated.Publisher = old.Publisher
			translated.JournalShortTitle = old.JournalShortTitle

		case "common":
			var old obsolete.MetadataCommon
			err = m.parseBlock(id, data, &old)

			translated.ArticleTitle = old.ArticleTitle
			translated.JournalTitle = old.JournalTitle
			translated.Volume = old.Volume
			translated.Publisher = old.Publisher
			translated.PublicationDate = old.PublicationDate
			translated.Abstract = old.Abstract
			for _, author := range old.Authors {
				translated.Authors = append(translated.Authors, model.Author{
					Name:  author.Name,
					Orcid: author.Orcid,
				})
			}
			translated.UnderEmbargo = old.UnderEmbargo
			translated.EmbargoEndDate = old.EmbargoEndDate
			if len(old.ISSN) > 0 {
				translated.ISSNs = append(translated.ISSNs, model.ISSN{ISSN: old.ISSN})
			}
			for issn, types := range old.IssnMap {
				pubType, err := findPubType(types.PubType)
				if err != nil {
					return nil, errors.Wrap(err, "error translating ISSNs")
				}
				translated.ISSNs = append(translated.ISSNs, model.ISSN{
					ISSN:    issn,
					PubType: pubType,
				})
			}
		case "JScholarship":
			var old obsolete.MetadataJ10p
			err = m.parseBlock(id, data, &old)

			if translated.Agreements == nil {
				translated.Agreements = make(map[string]string, 4)
			}
			translated.Agreements["JScholarship"] = old.AgreementText
			// Not copying Authors and Agreed, as they are redundant or unnecessary
		default:
			err = errors.Errorf("Unknown metadata block %s", id)
		}

		if err != nil {
			return nil, err
		}
	}

	return &translated, nil
}

func findPubType(types []string) (string, error) {
	if len(types) > 1 {
		return "", errors.Errorf("more than one publication type %s", types)
	}

	switch types[0] {
	case "Print":
		return "Print", nil
	case "Electronic":
		return "Online", nil
	}

	return "", errors.Errorf("unknown publication type %s", types[0])
}

func (m MetadataV0toV1) parseBlock(block string, metadata interface{}, blockPtr interface{}) error {
	serialized, err := json.Marshal(metadata)
	if err != nil {
		return errors.Wrapf(err, "could not marshal metadata from block %s", block)
	}

	decoder := json.NewDecoder(bytes.NewReader(serialized))
	decoder.DisallowUnknownFields()
	return errors.Wrapf(decoder.Decode(blockPtr), "could not decode %s", block)
}
