package client

import (
	"net/url"
	"strings"
)

type BaseURI string

// Join a uri path to a baseURI.  If the given path is already a URI, it is simply
// returned.
func (b BaseURI) Join(path string) string {
	if b == "" {
		return path
	}

	if strings.HasPrefix(path, "http") {
		return path
	}

	if strings.HasPrefix(path, "?") {
		return strings.Trim(string(b), "/") + path
	}

	return strings.Join([]string{strings.Trim(string(b), "/"), strings.TrimLeft(path, "/")}, "/")
}

// Given a URI whose path begins with the same path as the baseURI, replace the
// given URI's host and protocol with the baseURI's host and protocol
//
// (e.g https://pass.jhu.edu/fcrepo/rest/foo/bar -> http://localhost:8080/fcrepo/rest/foo/bar)
func (b BaseURI) Rebase(uri string) string {

	if !strings.HasPrefix(uri, "http") {
		return uri
	}

	this, err := url.Parse(string(b))
	if err != nil {
		return uri
	}

	that, err := url.Parse(uri)
	if err != nil {
		return uri
	}

	if !strings.HasPrefix(that.Path, this.Path) {
		return uri
	}

	var result strings.Builder
	result.WriteString(strings.TrimPrefix(that.Path, this.Path))

	if len(that.Query()) > 0 {
		result.WriteRune('?')
		result.WriteString(that.RawQuery)
	}

	if len(that.Fragment) > 0 {
		result.WriteRune('#')
		result.WriteString(that.Fragment)
	}

	return b.Join(result.String())
}
