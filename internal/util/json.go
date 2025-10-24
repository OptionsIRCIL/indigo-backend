package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type JsonDecodeFailed struct {
	Status int
	Msg    string
}

func (err *JsonDecodeFailed) Error() string {
	return err.Msg
}

// DecodeJSONBody mostly adapted from https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body /**
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
