package controller

import (
	"fmt"
	"net/http"
	"strings"
)

// https://developer.mozilla.org/en-US/docs/Glossary/CORS-safelisted_request_header
var corsSafelistedHeaders = []string{
	"Accept",
	"Accept-Language",
	"Content-Language",
	"Content-Type",
	"Range",
}

// ProvideOptions responds to an arbitrary OPTIONS request with a list of CORS headers
// detailing the methods and CORS requirements for the endpoint.
func ProvideOptions(methods []string) http.HandlerFunc {
	// TODO: Add Access-Control-Allow-Origin, config for allowed origins
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(
			"Access-Control-Allow-Methods",
			fmt.Sprintf("OPTIONS, %s", strings.Join(methods, ", ")),
		)
		w.Header().Add(
			"Access-Control-Allow-Headers",
			strings.Join(corsSafelistedHeaders, ", "),
		)
		w.WriteHeader(204)
	}
}
