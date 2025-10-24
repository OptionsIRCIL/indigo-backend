package controller

import (
	"backend/internal/service"
	"encoding/json"
	"net/http"
)

func IndexHelloWorld(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*service.LdapUser)

	responseBody, _ := json.Marshal(map[string]string{"Hello": user.FirstName})

	w.Write(responseBody)
	w.Header().Set("Content-Type", "application/json")
}
