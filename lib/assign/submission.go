// Package assign contains utilities for assigning "ownership" of PASS entities to specified users
package assign

import (
	"fmt"
	"net/http"

	"github.com/oa-pass/pass-tools/lib/client"
	"github.com/oa-pass/pass-tools/lib/es"
	"github.com/oa-pass/pass-tools/lib/log"
	"github.com/pkg/errors"
)

type Submission struct {
	Submitter string // URI of a submitter
	To        string // URI or local key of User to assign the grant to
	Grant     string
	BaseURI   client.BaseURI
	Fedora    client.Performer
	Elastic   client.Performer
	DryRun    bool
	Log       log.Instance
}

func (s Submission) Perform() (err error) {

	submissions, err := s.withSubmitterAndGrant(s.Submitter, s.Grant)
	if err != nil {
		return errors.Wrapf(err, "could not find submissions")
	}
	s.Log.Debugf("Found %d submissions to assign", len(submissions))

	for _, submission := range submissions {
		if !s.DryRun {
			err = s.Fedora.Perform(http.MethodPatch, s.BaseURI.Rebase(submission), &client.Body{
				Content: fmt.Sprintf(`{
				"@context" : "%s",
				"@id" : "",
				"submitter": "%s",
				"@type" : "Submission"
			}`, client.Context, s.BaseURI.Rebase(s.To)),
				Type: client.ContentTypeJSONMerge,
			}, nil)
			if err != nil {
				s.Log.Warnf("Could not assign submission %s to %s: %s", submission, s.To, err)
			}
			s.Log.Printf("Assigned submission %s to %s", submission, s.To)
		} else {
			s.Log.Printf("Would have assigned submission %s to %s", s.BaseURI.Rebase(submission), s.BaseURI.Rebase(s.To))
		}
	}

	return err
}

func (s Submission) withSubmitterAndGrant(submitter, grant string) ([]string, error) {

	var results es.IDResults
	err := s.Elastic.Perform(http.MethodPost, "", &client.Body{
		Content: es.QueryMatch(map[string]string{
			"@type":     "Submission",
			"submitter": submitter,
			"grants":    grant,
		}, 100),
		Type: client.ContentTypeJSON}, &results)

	if err != nil {
		return nil, errors.Wrap(err, "submission query failed")
	}

	submissions := make([]string, 0, results.Hits.Total)
	for _, hit := range results.Hits.Hit {
		submissions = append(submissions, hit.Source.ID)
	}

	return submissions, nil
}
