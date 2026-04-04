package controller

import (
	"net/http"

	"myoptions.info/indigo/backend/internal/util"
)

func IndexHelloWorld(w http.ResponseWriter, r *http.Request) {
	token := util.FetchTokenFromContext(r)
	util.ReturnSerialized(w, 200, map[string]string{
		"Hello": token.FirstName + " " + token.LastName + " (" + token.Subject + ")",
	}, nil)
}
