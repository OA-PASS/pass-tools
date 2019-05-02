package es

import (
	"fmt"
	"strings"
)

// QueryMatch creates an elasticsearch match query from the given terms
func QueryMatch(terms map[string]string, max int) string {
	conditions := make([]string, 0, len(terms))
	for k, v := range terms {
		conditions = append(conditions, fmt.Sprintf(`{
			"match": {
				"%s": "%s"
			}
		}`, k, v))
	}

	return fmt.Sprintf(`{
		"size":%d,
		"query" : {
		  "bool" : {
			"must": [ 
				%s
			]
		  }
		}
	  }`, max, strings.Join(conditions, ",\n"))
}
