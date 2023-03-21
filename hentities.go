package hentities

import (
	"regexp"
	"strings"
)

type (
	// Query represents a GraphQL query.
	Query string

	// Mutation represents a GraphQL mutation.
	Mutation string

	// query is used internally to represent both Query and Mutations.
	query string

	// graphqlRequest is the JSON request to a GraphQL server.
	graphQLRequest struct {
		Query     query       `json:"query"`
		Variables interface{} `json:"variables,omitempty"`
	}
)

// root computes the request root of the query. This is used internally to map
// the internal json structure inside `data` to the provided interface{} in
// the request.
func (q query) root() string {
	rp := regexp.MustCompile(`(query|mutation|subscription)\s*([a-zA-Z0-9]*\s*)?(\(.+\)\s*)?`)
	fq := string(rp.ReplaceAll([]byte(q), []byte("")))

	const zero, one, three = 0, 1, 3
	results := strings.SplitN(fq, "{", three)

	if len(results) == three {
		rp2 := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)`)
		matches := rp2.FindStringSubmatch(results[one])
		if len(matches) > one {
			return matches[zero]
		}
	}

	return ""
}
