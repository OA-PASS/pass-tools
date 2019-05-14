// Package assign contains utilities for assigning "ownership" of PASS entities to specified users
package assign

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/oa-pass/pass-tools/lib/client"
	"github.com/oa-pass/pass-tools/lib/es"
	"github.com/oa-pass/pass-tools/lib/log"
	"github.com/oa-pass/pass-tools/lib/model"
	"github.com/pkg/errors"
)

type Grant struct {
	ID          string // local key of a grant
	To          string // URI or local key of User to assign the grant to
	BaseURI     client.BaseURI
	Submissions bool // Assign submissions where submitter is the old PI
	Fedora      client.Performer
	Elastic     client.Performer
	DryRun      bool
	Log         log.Instance
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

	if !g.DryRun {
		err := g.Fedora.Perform(http.MethodPatch, g.BaseURI.Rebase(grant.ID), &client.Body{
			Content: fmt.Sprintf(`{
			"@context" : "%s",
			"@id" : "",
			"pi": "%s",
			"@type" : "Grant"
		}`, client.Context, g.BaseURI.Rebase(user)),
			Type: client.ContentTypeJSONMerge,
		}, nil)

		if err != nil {
			return errors.Wrapf(err, "could not assign grant %s to user %s", g.ID, g.To)
		}
	} else {
		g.Log.Printf("Would have assigned grant %s (%s) to user %s (%s)",
			g.ID, g.BaseURI.Rebase(grant.ID), g.To, g.BaseURI.Rebase(user))
	}

	if !g.Submissions {
		g.Log.Printf("Not assigning submissions")
		return nil
	}

	return Submission{
		Submitter: grant.PI,
		To:        user,
		Grant:     grant.ID,
		BaseURI:   g.BaseURI,
		Fedora:    g.Fedora,
		Elastic:   g.Elastic,
		Log:       g.Log,
		DryRun:    g.DryRun,
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
		g.Log.Debugf("Got an http URI for the grant, Getting it")
		var grant model.Grant
		return &grant, g.Fedora.Perform(http.MethodGet, id, nil, &grant)
	}

	var results es.GrantResults
	g.Log.Debugf("Findng grant %s", id)
	err := g.find("Grant", "localKey", id, &results)

	if err != nil {
		return nil, errors.Wrapf(err, "search for grant %s failed", id)
	}

	if results.Hits.Total != 1 {
		return nil, errors.Errorf("Expected one grant for %s, got %d %+v", id, results.Hits.Total, results)
	}

	return &results.Hits.Hit[0].Source, err
}

func (g Grant) find(t, field, key string, resultsPtr interface{}) error {

	return errors.Wrapf(g.Elastic.Perform(http.MethodPost, "", &client.Body{
		Content: es.QueryMatch(map[string]string{
			"@type": t,
			field:   key,
		}, 2), Type: client.ContentTypeJSON,
	}, resultsPtr), "elasticsearch query for %s ($%s = %s) failed", t, key, field)
}
