package controller

import (
	"fmt"
	"net/http"
	"strings"
)

// ProvideOptions responds to an arbitrary OPTIONS request with a list of CORS headers
// detailing the methods and CORS requirements for the endpoint.
func ProvideOptions(methods []string) http.HandlerFunc {
	// TODO: Add Access-Control-Allow-Origin, config for allowed origins
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(
			"Access-Control-Allow-Methods",
			fmt.Sprintf("OPTIONS, %s", strings.Join(methods, ", ")),
		)
		w.WriteHeader(204)
	}
}
