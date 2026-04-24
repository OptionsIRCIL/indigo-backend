package controller

import (
	"net/http"

	"myoptions.info/indigo/backend/internal/util"
)

func IndexHelloWorld(w http.ResponseWriter, r *http.Request) {
	util.ThrowHttpStatus(w, 404)
}
