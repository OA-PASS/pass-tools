package client_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/oa-pass/pass-tools/lib/client"
)

type fakeRequester struct {
	f func(req *http.Request) (*http.Response, error)
}

func (r *fakeRequester) Do(req *http.Request) (*http.Response, error) {
	if r.f != nil {
		return r.f(req)
	}

	return nil, nil
}

type fakeBody struct {
	io.Reader
	closeFunc func() error
}

func (b *fakeBody) Close() error {
	if b.closeFunc != nil {
		return b.closeFunc()
	}
	return nil
}

func TestSimpleFetchEntity(t *testing.T) {
	baseURI := client.BaseURI("https://example.org/foo/fcrepo/rest")

	resource := baseURI.Join("/foo/bar")

	username := "foo"
	password := "bar"

	passClient := client.Simple{
		Requester: &fakeRequester{
			f: func(req *http.Request) (*http.Response, error) {

				if req.URL.String() != resource {
					t.Fatalf("resource URI is incorrect")
				}

				user, pass, ok := req.BasicAuth()
				if !ok || user != username || pass != password {
					t.Fatalf("basic auth is wrong")
				}

				return &http.Response{
					Body: &fakeBody{
						Reader: strings.NewReader(`{
							"foo" : [
								"bar",
								"baz"
							]
						}`),
					},
				}, nil
			},
		},
		BaseURI: baseURI,
		Credentials: &client.Credentials{
			Username: username,
			Password: password,
		},
	}

	ref := make(map[string]interface{})
	err := passClient.Read(resource, &ref)
	if err != nil {
		t.Fatalf("Client fetch resulted in error %+v", err)
	}

	diffs := deep.Equal(ref["foo"].([]interface{}), []interface{}{"bar", "baz"})
	if len(diffs) > 0 {
		t.Fatalf("found difference in deserialized content %s", diffs)
	}
}

func TestSimpleEntityErrors(t *testing.T) {
	cases := []struct {
		name string
		url  string
		f    func(req *http.Request) (*http.Response, error)
	}{
		{
			name: "badURI",
			url:  "0http://bad",
		},
		{
			name: "httpError",
			url:  "http://example.org/foo",
			f: func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("this is an error")
			},
		},
		{
			name: "badJSON",
			url:  "http://example.org/foo",
			f: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					Body: &fakeBody{
						Reader: strings.NewReader(`{BAD JSON-,`),
					},
				}, nil
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			passClient := client.Simple{
				Requester: &fakeRequester{
					f: c.f,
				},
			}

			ref := make(map[string]interface{})

			err := passClient.Read(c.url, &ref)
			if err == nil {
				t.Fatalf("Should have terminated with an error")
			}
		})
	}
}
