// Package assign contains utilities for assigning "ownership" of PASS entities to specified users
package assign

import (
	"fmt"
	"log"
	"net/http"

	"github.com/oa-pass/pass-tools/lib/client"
	"github.com/oa-pass/pass-tools/lib/es"
	"github.com/pkg/errors"
)

type Submission struct {
	Submitter string // URI of a submitter
	To        string // URI or local key of User to assign the grant to
	Fedora    client.Performer
	Elastic   client.Performer
}

func (s Submission) Perform() error {

	submissions, err := s.withSubmitter(s.Submitter)
	if err != nil {
		return errors.Wrapf(err, "could not find submissions")
	}

	for _, submission := range submissions {
		err = s.Fedora.Perform(http.MethodPatch, submission, &client.Body{
			Content: fmt.Sprintf(`{
				"@context" : "%s",
				"@id" : "",
				"submitter": "%s",
				"@type" : "Submission"
			}`, client.Context, s.To),
			Type: client.ContentTypeJSONMerge,
		}, nil)
		if err != nil {
			return errors.Wrap(err, "could not update submission")
		}

		log.Printf("Assigned submission %s to %s", submission, s.To)
	}

	return nil
}

func (s Submission) withSubmitter(submitter string) ([]string, error) {

	var results es.IDResults
	err := s.Elastic.Perform(http.MethodPost, "", &client.Body{
		Content: es.QueryMatch(map[string]string{
			"@type":     "Submission",
			"submitter": submitter,
		}, 100),
		Type: client.ContentTypeJSON}, &results)

	if err != nil {
		return nil, errors.Wrap(err, "submission query failed")
	}

	submissions := make([]string, 0, results.Hits.Total)
	for _, hit := range results.Hits.Hit {
		submissions = append(submissions, es.RelativeURI(hit.Source.ID))
	}

	return submissions, nil
}
