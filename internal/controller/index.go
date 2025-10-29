package controller

import (
	"backend/internal/service"
	"backend/internal/util"
	"net/http"
)

func IndexHelloWorld(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*service.LdapUser)
	util.ReturnSerialized(w, 200, map[string]string{
		"Hello": user.FirstName,
	})
}
