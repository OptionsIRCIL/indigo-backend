package util

import (
	"encoding/json"
	"log"
	"net/http"
)

type HttpError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var httpStatusDescriptions = map[int]string{
	401: "Unauthorized",
	403: "Forbidden",
	500: "Internal server error",
}

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

func ThrowHttpError(w http.ResponseWriter, status int, message string) {
	ReturnSerialized(w, status, &HttpError{
		Status:  status,
		Message: message,
	})
}

func ThrowHttpStatus(w http.ResponseWriter, status int) {
	msg, ok := httpStatusDescriptions[status]
	if ok {
		ThrowHttpError(w, status, msg)
	} else {
		log.Print(
			"Request to throw to non-cataloged HTTP status",
			status,
			"rewritten to 500: Internal server error",
		)
		fiveHundred, _ := httpStatusDescriptions[500]
		ThrowHttpError(w, 500, fiveHundred)
	}
}

func ThrowHttpUnhandled(w http.ResponseWriter, e error) {
	log.Print("Unhandled error:", e.Error())
	ThrowHttpStatus(w, 500)
}
