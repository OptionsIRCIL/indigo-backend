package util

import (
	"encoding/json"
	"log"
	"net/http"
)

// HttpError is a serializable struct to be returned on any 4xx or 5xx errors.
// By utilizing this struct, the client will be able to display more descriptive error messages.
type HttpError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var httpStatusDescriptions = map[int]string{
	401: "Unauthorized",
	403: "Forbidden",
	500: "Internal server error",
}

// ReturnSerialized marshals a struct of any type into a JSON string and adds
// it to an [http.ResponseWriter]. The status code and content type are also set.
func ReturnSerialized(w http.ResponseWriter, status int, payload any) {
	serialized, err := json.Marshal(payload)
	if err != nil {
		// We might get stuck in a recursion loop doing this but i don't give a swag atm
		ThrowHttpUnhandled(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(serialized)
}

// ThrowHttpError writes both an HTTP status and message in a serialized format,
// [HttpError], to an [http.ResponseWriter].
func ThrowHttpError(w http.ResponseWriter, status int, message string) {
	ReturnSerialized(w, status, &HttpError{
		Status:  status,
		Message: message,
	})
}

// ThrowHttpStatus will write a suitable generic response body given a known
// HTTP status code.
func ThrowHttpStatus(w http.ResponseWriter, status int) {
	msg, ok := httpStatusDescriptions[status]
	if ok {
		ThrowHttpError(w, status, msg)
	} else {
		log.Print(
			"Request to throw to non-cataloged HTTP status ",
			status,
			" rewritten to 500: Internal server error",
		)
		fiveHundred, _ := httpStatusDescriptions[500]
		ThrowHttpError(w, 500, fiveHundred)
	}
}

// ThrowHttpUnhandled wraps around ThrowHttpStatus with code 500.
// Additionally, it writes the encountered error to the server log.
func ThrowHttpUnhandled(w http.ResponseWriter, e error) {
	log.Print("Unhandled error: ", e.Error())
	ThrowHttpStatus(w, 500)
}
