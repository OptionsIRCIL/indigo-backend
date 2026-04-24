package controller

import (
	"net/http"

	"myoptions.info/indigo/backend/internal/util"
)

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		util.ReturnSerialized(w, http.StatusNotFound, util.HttpError{
			Message: "The requested resource was not found.",
		}, nil)
		return
	}

	token := util.FetchTokenFromContext(r)

	response := map[string]interface{}{
		"status":  "online",
		"version": "1.0.0",
		"user":    token.FirstName + " " + token.LastName,
		"context": token.Subject,
		"message": "Indigo CIL Backend API is operational.",
	}

	util.ReturnSerialized(w, http.StatusOK, response, nil)
}
