package hentities

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

type (
	// Repository is a base struct containing convenience methods for
	// repositories to perform queries and mutations.
	Repository struct {
		endpoint string
		client   *resty.Client
		headers  map[string]string
	}

	// RepositoryOption are the options that can be provided to during
	// construction of a Repository instance.
	RepositoryOption func(r *Repository)
)

type (
	// Call represents the data that will be sent to a GraphQL server
	// to perform either a query or mutation.
	call struct {
		query     query
		variables interface{}
		headers   map[string]string
	}

	// CallOption are the options that can modify the call parameters.
	callOption func(*call)
)

// WithClient is a RepositoryOption that provides a custom http.Client
// that will be used when calling the GraphQL server.
func WithClient(c *http.Client) RepositoryOption {
	return func(r *Repository) {
		r.client = resty.NewWithClient(c)
	}
}

// WithGlobalHeaders provides global headers that will be part of all
// the requests coming from the repository. Note that the header
// `Content-Type: application/json` is always present even without
// calling this method.
func WithGlobalHeaders(headers map[string]string) RepositoryOption {
	return func(r *Repository) {
		r.headers = headers
	}
}

// WithVariables adds the variables to the GraphQL request. This is made
// optional as not all queries will need variables. Use when needed.
func WithVariables(variables interface{}) callOption {
	return func(c *call) {
		c.variables = variables
	}
}

// WithRequestHeader adds additional header to the GraphQL request. In
// contrast to the global headers, this can be added on individual requests
// and will be appended to the global headers.
func WithRequestHeader(headers map[string]string) callOption {
	return func(c *call) {
		if c.headers == nil {
			c.headers = map[string]string{}
		}
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}

// NewRepository creates a new instance of Repository.
func NewRepository(endpoint string, opts ...RepositoryOption) *Repository {
	r := &Repository{endpoint: endpoint}
	for _, opt := range opts {
		opt(r)
	}
	if r.client == nil {
		r.client = resty.New()
	}
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	r.headers["Content-Type"] = "application/json"
	return r
}

// Query is a convenience method for performing GraphQL queries.
//
// The parameters are:
//   - ctx      - the context. This will return an error if nil.
//   - q        - the Query
//   - response - The struct or map that will contain the response.
//     Take note that the response does not need to
//     include the parent "data" field of a typical
//     GraphQL response. Instead, this will be the root
//     record instead.
//   - opts     - Additional options for the request.
func (r *Repository) Query(ctx context.Context, q Query, response interface{}, opts ...callOption) error {
	return r.request(ctx, query(q), response, opts...)
}

// Mutate is a convenience method for performing GraphQL mutations.
//
// The parameters are:
//   - ctx      - the context. This will return an error if nil.
//   - mu       - the Mutation
//   - response - The struct or map that will contain the response.
//     Take note that the response does not need to
//     include the parent "data" field of a typical
//     GraphQL response. Instead, this will be the root
//     record instead.
//   - opts     - Additional options for the request.
func (r *Repository) Mutate(ctx context.Context, mut Mutation, response interface{}, opts ...callOption) error {
	return r.request(ctx, query(mut), response, opts...)
}

// request does the actual HTTP request to the GraphQL server. This
// is currently called by both Query() and Mutate() methods. However,
// this may change if GraphQL specifications change or there are very
// different behaviors between Query and Mutations.
//
// This method does the following:
// - construct the request
// - parse the response depending if it's data, errors, or error
//   - data   - is a successful request
//   - errors - there's an error in performing the query
//   - error  - generic error that is not related to GraphQL (eg. 404)
func (r *Repository) request(ctx context.Context, q query, response interface{}, opts ...callOption) error {
	call := &call{query: q}
	for _, opt := range opts {
		opt(call)
	}
	req := call.toRequest(ctx, r.client, r.headers)
	res, err := req.Post(r.endpoint)
	if err != nil {
		// return ErrorConnection
	}

	defer res.RawBody().Close()

	body := res.Body()
	resdata := gjson.GetBytes(body, "data")
	reserrs := gjson.GetBytes(body, "errors")
	reserr := gjson.GetBytes(body, "error")

	switch {
	case resdata.Exists() && resdata.Type != gjson.Null:
		if response != nil {
			root := q.root()
			content := []byte(gjson.Get(resdata.Raw, root).Raw)
			if err := json.Unmarshal(content, response); err != nil {
				// return ErrorParseResponse
			}
		}

		return nil

	case reserrs.Exists():
		var actualError []map[string]interface{}
		if err := json.Unmarshal([]byte(reserrs.Raw), &actualError); err != nil {
			// return ErrorParseResponse
		}
		// attempt to return terror if found, first one found is returned.
		// if there is no terror, then we return the actualError as a metadata.
		// if an error_code is found, but no status code, it will default to 400.
		// if metadata is found, it will be returned as well.
		// for _, err := range actualError {
		// terror is usually under extensions
		// e, ok := err["extensions"]
		// if !ok {
		// 	continue
		// }
		// extensions, ok := e.(map[string]interface{})
		// if !ok {
		// 	continue
		// }
		// under error
		// er, ok := extensions["error"]
		// if !ok {
		// 	continue
		// }
		// errRes, ok := er.(map[string]interface{})
		// if !ok {
		// 	continue
		// }
		// with error code
		// ec, ok := errRes["error_code"]
		// if !ok {
		// 	continue
		// }
		// errorCode, ok := ec.(string)
		// if !ok {
		// 	continue
		// }
		// // status code
		// sc, ok := errRes["status_code"]
		// if !ok {
		// 	// return errors.New(errorCode, errors.WithStatusCode(http.StatusBadRequest))
		// 	return nil
		// }
		// statusCode, ok := sc.(int)
		// if !ok {
		// 	// return errors.New(errorCode, errors.WithStatusCode(http.StatusBadRequest))
		// 	return nil
		// }
		// and optional metadata
		// md, ok := errRes["metadata"]
		// if !ok {
		// 	// return errors.New(errorCode, errors.WithStatusCode(statusCode))
		// 	return nil
		// }
		// metadata, ok := md.(map[string]interface{})
		// if !ok {
		// 	// return errors.New(errorCode, errors.WithStatusCode(statusCode))
		// 	return nil
		// }
		// return errors.New(errorCode, errors.WithStatusCode(statusCode), errors.WithMetadata(metadata))
		return nil
		// }
		// return newRequestError(actualError)

	case reserr.Exists():
		return newRequestError(errors.New(reserr.String()))
	default:
		// return ErrorParseResponse
	}
	return err

}

// toRequest translates the call into a resty request.
func (c call) toRequest(ctx context.Context, client *resty.Client, globalHeaders map[string]string) *resty.Request {
	return client.R().
		SetContext(ctx).
		SetHeaders(globalHeaders).
		SetHeaders(c.headers).
		SetBody(graphQLRequest{
			Query:     c.query,
			Variables: c.variables,
		})
}
