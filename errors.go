package hentities

var (
// ErrorConnection is returned when there is an error connecting to the GraphQL server.
// ErrorConnection = errors.New("CONN_ERROR")
// ErrorParseResponse is returned when there is an error reading or parsing the
// response from the GraphQL server.
// ErrorParseResponse = errors.New("RESPONSE_ERROR")
)

// newRequeserrors are request errors and contains the internal error as a metadata.
func newRequestError(errs interface{}) error {
	return nil
	// return errors.New(
	// 	"REQUEST_ERROR",
	// 	errors.WithStatusCode(http.StatusBadRequest),
	// 	errors.WithMetadata(map[string]interface{}{
	// 		"errors": errs,
	// 	}))
}
