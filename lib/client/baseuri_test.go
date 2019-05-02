package client_test

import (
	"testing"

	"github.com/oa-pass/pass-tools/lib/client"
)

func TestPrivateWithPublic(t *testing.T) {

	cases := []struct {
		testName string
		baseuri  string
		path     string
		expected string
	}{{
		testName: "slash slash",
		baseuri:  "http://example.org/foo/",
		path:     "/bar",
		expected: "http://example.org/foo/bar",
	}, {
		testName: "slash no",
		baseuri:  "http://example.org/foo/",
		path:     "bar",
		expected: "http://example.org/foo/bar",
	}, {
		testName: "no slash",
		baseuri:  "http://example.org/foo",
		path:     "/bar",
		expected: "http://example.org/foo/bar",
	}, {
		testName: "no no",
		baseuri:  "http://example.org/foo",
		path:     "bar",
		expected: "http://example.org/foo/bar",
	}, {
		testName: "uri",
		baseuri:  "http://example.org/foo",
		path:     "http://foo.bar/baz",
		expected: "http://foo.bar/baz",
	}, {
		testName: "query",
		baseuri:  "http://example.org/foo",
		path:     "?bar=baz",
		expected: "http://example.org/foo?bar=baz",
	}}

	for _, c := range cases {
		c := c
		t.Run(c.testName, func(t *testing.T) {
			result := client.BaseURI(c.baseuri).Join(c.path)

			if result != c.expected {
				t.Fatalf("%s did not match expected %s", result, c.expected)
			}
		})
	}
}
