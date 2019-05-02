// Package assign contains utilities for assigning "ownership" of PASS entities to specified users
package assign

import (
	"net/http"
	"strings"

	"fmt"

	"github.com/oa-pass/pass-tools/lib/client"
	"github.com/oa-pass/pass-tools/lib/es"
	"github.com/oa-pass/pass-tools/lib/model"
	"github.com/pkg/errors"
)

type Grant struct {
	ID          string // local key of a grant
	To          string // URI or local key of User to assign the grant to
	Submissions bool   // Assign submissions where submitter is the old PI
	Fedora      client.Performer
	Elastic     client.Performer
}

func (g Grant) Perform() error {
	grant, err := g.findGrant(g.ID)
	if err != nil {
		return errors.Wrap(err, "failed to find grant")
	}

	user, err := g.findUser(g.To)
	if err != nil {
		return errors.Wrapf(err, "failed to find user")
	}

	err = g.Fedora.Perform(http.MethodPatch, grant.ID, &client.Body{
		Content: fmt.Sprintf(`{
			"@context" : "%s",
			"@id" : "",
			"pi": "%s"
			"@type" : "Grant"
		}`, client.Context, user),
		Type: client.ContentTypeJSONMerge,
	}, nil)

	if !g.Submissions {
		return errors.Wrap(err, "could not update grant")
	}

	return Submission{
		Submitter: es.RelativeURI(grant.PI),
		To:        user,
		Fedora:    g.Fedora,
		Elastic:   g.Elastic,
	}.Perform()
}

func (g Grant) findUser(id string) (string, error) {
	if strings.HasPrefix(id, "http") {
		return id, nil
	}
	var results es.IDResults
	err := g.find("User", "locatorIds", id, &results)

	if results.Hits.Total != 1 {
		return id, errors.Errorf("Expected one user for %s, got %d", id, results.Hits.Total)
	}
	return results.Hits.Hit[0].Source.ID, err
}

func (g Grant) findGrant(id string) (*model.Grant, error) {

	if strings.HasPrefix(id, "http") {
		var grant model.Grant
		return &grant, g.Fedora.Perform(http.MethodGet, id, nil, &grant)
	}

	var results es.GrantResults
	err := g.find("Grant", "localKey", id, &results)

	if results.Hits.Total != 1 {
		return nil, errors.Errorf("Expected one grant for %s, got %d", id, results.Hits.Total)
	}

	return &results.Hits.Hit[0].Source, err
}

func (g Grant) find(t, field, key string, resultsPtr interface{}) error {
	return errors.Wrapf(g.Elastic.Perform(http.MethodPost, "", &client.Body{
		Content: es.QueryMatch(map[string]string{
			"@type": t,
			field:   key,
		}, 2),
		Type: client.ContentTypeJSON}, resultsPtr), "elasticsearch query for %s ($%s = %s) failed", t, key, field)
}
