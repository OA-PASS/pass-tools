package client_test

import (
	"testing"

	"github.com/oa-pass/pass-tools/lib/client"
)

func TestJoin(t *testing.T) {

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

func TestRebase(t *testing.T) {
	cases := []struct {
		testName string
		baseuri  string
		path     string
		expected string
	}{{
		testName: "basic",
		baseuri:  "http://example.org/foo/bar/",
		path:     "http://example.org/foo/bar/baz/foo",
		expected: "http://example.org/foo/bar/baz/foo",
	}, {
		testName: "different hosts",
		baseuri:  "http://example.org/foo/bar/",
		path:     "http://foo.local:8080/foo/bar/baz/foo",
		expected: "http://example.org/foo/bar/baz/foo",
	}, {
		testName: "query and hash",
		baseuri:  "http://example.org/foo/bar/",
		path:     "http://foo.local:8080/foo/bar/baz/foo?foo=bar#baz",
		expected: "http://example.org/foo/bar/baz/foo?foo=bar#baz",
	}, {
		testName: "different base paths",
		baseuri:  "http://example.org/foo/bar/",
		path:     "http://foo.local:8080/baz/foo",
		expected: "http://foo.local:8080/baz/foo",
	}, {
		testName: "not a uri",
		baseuri:  "http://example.org/foo/bar/",
		path:     "foo/bar/baz/foo",
		expected: "foo/bar/baz/foo",
	}}

	for _, c := range cases {
		c := c
		t.Run(c.testName, func(t *testing.T) {
			result := client.BaseURI(c.baseuri).Rebase(c.path)

			if result != c.expected {
				t.Fatalf("%s did not match expected %s", result, c.expected)
			}
		})
	}
}
