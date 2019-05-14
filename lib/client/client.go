package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/oa-pass/pass-tools/lib/log"
	"github.com/pkg/errors"
)

const (
	headerAccept      = "Accept"
	headerContentType = "Content-Type"
	headerUserAgent   = "User-Agent"
)

const (
	ContentTypeJSON      = "application/json"
	ContentTypeJSONLD    = "application/ld+json"
	ContentTypeJSONMerge = "application/merge-patch+json"
)

const (
	mediaJSONTypes = "application/json, application/ld+json"
)

const Context = "https://oa-pass.github.io/pass-data-model/src/main/resources/context-3.4.jsonld"

// Simple PASS client for hitting a server (Fedora, Elasticsearch, etc)
// and unmarshalling JSON results.
type Simple struct {
	Requester
	BaseURI     BaseURI
	Credentials *Credentials
	Log         log.Instance
}

// Credentials for Basic Auth
type Credentials struct {
	Username string
	Password string
}

// Requester performs http requests
type Requester interface {
	Do(req *http.Request) (*http.Response, error)
}

// Reader reads an the content at the given URL, and unmarshals it into
// the given resultPointer, which is a pointer to a struct or map
type Reader interface {
	Read(url string, resultPointer interface{}) error
}

type Body struct {
	Content interface{} // bytes(string, array, reader), or object to marshal into JSON
	Type    string      // http content type
}

// Performer sends the given body to the given url (e.g. a POST or PATCH).
// If resultPointer is a non-nil pointer to a struct or map, and the result
// is JSON, it unmarshals it to the given resultPointer,
//
// The body may be bytes ([]byte, or io.Reader), or a struct that can be
// unmarshaled to JSON.
type Performer interface {
	Perform(method, url string, body *Body, resultPointer interface{}) error
}

func (s *Simple) Read(url string, resultPointer interface{}) error {
	return s.Perform(http.MethodGet, url, nil, resultPointer)
}

func (s *Simple) Perform(method, url string, body *Body, resultPointer interface{}) (err error) {

	url = s.BaseURI.Join(url)
	var sniff bytes.Buffer

	reader, err := toReader(body)
	if err != nil {
		return errors.Wrapf(err, "could not form input body from %s", body)
	}

	if s.Log.Trace != nil {
		reader = io.TeeReader(reader, &sniff)
	}

	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return errors.Wrapf(err, "could not build http request to %s", url)
	}

	if s.Credentials != nil {
		request.SetBasicAuth(s.Credentials.Username, s.Credentials.Password)
	}
	request.Header.Set(headerUserAgent, "pass-tools")
	request.Header.Set(headerAccept, mediaJSONTypes)
	if body != nil && body.Type != "" {
		request.Header.Set(headerContentType, body.Type)
	}

	resp, err := s.Do(request)
	if err != nil {
		return errors.Wrapf(err, "error connecting to %s", url)
	}
	defer resp.Body.Close()

	if s.Log.Trace != nil {
		sent := sniff.String()

		var headers strings.Builder

		for k, v := range request.Header {
			fmt.Fprintf(&headers, "  %s: %s\n", k, strings.Join(v, ", "))
		}

		if len(sent) > 0 {
			s.Log.Tracef("Sent %s to %s with headers:\n%s\n...and body:\n%s", method, url, headers.String(), sent)
		} else {
			s.Log.Tracef("Sent empty %s to %s", method, url)
		}

		s.Log.Tracef("Got response code %d from %s to %s", resp.StatusCode, method, url)
	}

	if resp.StatusCode > 299 {
		errBody, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf("Request failed with code %d: %s", resp.StatusCode, string(errBody))
	}

	switch dest := resultPointer.(type) {
	case func(io.Reader) error:
		return dest(resp.Body)
	case nil:
		_, err = io.Copy(ioutil.Discard, resp.Body)
	default:
		return json.NewDecoder(resp.Body).Decode(resultPointer)
	}

	return err
}

func toReader(b *Body) (io.Reader, error) {

	if b == nil {
		return nil, nil
	}

	switch body := b.Content.(type) {
	case nil:
		return nil, nil
	case string:
		return strings.NewReader(body), nil
	case io.Reader:
		return body, nil
	case []byte:
		return bytes.NewReader(body), nil
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(b.Content)
	if err != nil {
		return nil, errors.Wrap(err, "could not encode JSON")
	}

	return &buf, nil
}
