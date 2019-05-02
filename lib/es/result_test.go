package es_test

import (
	"testing"

	"github.com/oa-pass/pass-tools/lib/es"
)

func TestRelativeURI(t *testing.T) {
	cases := map[string]string{
		"http://foo.local/fcrepo/rest/submissions/123": "/submissions/123",
		"http://foo.local/some/other/path":             "http://foo.local/some/other/path",
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			if es.RelativeURI(k) != v {
				t.Fatalf("Got %s instead of %s", es.RelativeURI(k), v)
			}
		})
	}
}
