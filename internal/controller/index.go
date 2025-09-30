package controller

import (
	"encoding/json"
	"net/http"
)

func IndexHelloWorld(w http.ResponseWriter, r *http.Request) {
	responseBody, err := json.Marshal(map[string]string{"Hello": "World"})
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(responseBody)
	w.Header().Set("Content-Type", "application/json")
	return

}
