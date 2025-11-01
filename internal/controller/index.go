package controller

import (
	"net/http"

	"myoptions.info/indigo/backend/internal/service"
	"myoptions.info/indigo/backend/internal/util"
)

func IndexHelloWorld(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*service.LdapUser)
	util.ReturnSerialized(w, 200, map[string]string{
		"Hello": user.FirstName,
	})
}
