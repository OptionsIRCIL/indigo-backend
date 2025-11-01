package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// A JsonDecodeFailed abstracts away any error that may be encountered while
// scanning to a struct in [util.DecodeJsonBody].
type JsonDecodeFailed struct {
	// An HTTP status relevant to the reason that scanning failed.
	Status int

	// A brief description of the error. Any 4xx errors may be shown to the
	// user, but 5xx errors should be only shown in server logs.
	Msg string
}

// Error interface. Returns err.Msg.
func (err *JsonDecodeFailed) Error() string {
	return err.Msg
}

// DecodeJSONBody scans the contents of an [http.ResponseWriter]'s body into a given struct iff the Content-Type of the
// body is application/json and the payload conforms to the struct. Adapted from an [Alex Edwards article] pertaining
// to parsing JSON request bodies.
//
// [Alex Edwards article]: https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, target interface{}) *JsonDecodeFailed {
	// Require "application/json" content type
	contentType := r.Header.Get("Content-Type")
	contentTypeMatches, err := regexp.MatchString(`^application/json(?:$|;)`, strings.TrimSpace(strings.ToLower(contentType)))

	if err != nil {
		return &JsonDecodeFailed{500, "Invalid regexp in DecodeJSONBody"}
	}
	if !contentTypeMatches {
		return &JsonDecodeFailed{422, "Invalid content type"}
	}

	// Max payload of 2MB
	r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&target)
	if err != nil {
		var syntaxError *json.SyntaxError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
		case errors.Is(err, io.ErrUnexpectedEOF):
			return &JsonDecodeFailed{422, "Malformed request body"}
		case errors.As(err, &maxBytesError):
			return &JsonDecodeFailed{413, "Request entity too large"}
		default:
			return &JsonDecodeFailed{400, "Bad request"}
		}
	}

	// Ensure we hit EOF
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return &JsonDecodeFailed{422, "Unexpected trailing content"}
	}

	return nil
}
