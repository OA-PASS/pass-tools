package client

import (
	"strings"
)

type BaseURI string

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
