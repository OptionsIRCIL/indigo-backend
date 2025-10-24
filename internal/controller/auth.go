package controller

import (
	"backend/internal/service"
	u "backend/internal/util"
	"log"
	"net/http"
	"regexp"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AuthEntry(conn *service.LdapConnection, transformer *service.JwtTransformer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Better status code handling
		var payload UserCredentials

		// Decode payload
		decodeErr := u.DecodeJSONBody(w, r, &payload)
		if decodeErr != nil {
			log.Printf("Encountered error %d: %s", decodeErr.Status, decodeErr.Msg)
			w.WriteHeader(decodeErr.Status)
			if decodeErr.Status > 499 {
				w.Write([]byte("Internal server error"))
			} else {
				w.Write([]byte(decodeErr.Msg))
			}
			return
		}

		// 422 if the username has any invalid chars
		usernameCharsOk, regexErr := regexp.MatchString(`^[A-Za-z0-9-_.]+$`, payload.Username)
		if regexErr != nil {
			w.WriteHeader(500)
			w.Write([]byte("Internal server error"))
			return
		} else if !usernameCharsOk {
			w.WriteHeader(422)
			w.Write([]byte("Invalid credentials"))
			return
		}

		// Parse out username and attempt bind auth
		fqUsername := "CN=" + payload.Username + "," + conn.Base
		authErr := conn.AttemptAuth(fqUsername, payload.Password)
		if authErr != nil {
			log.Print(authErr)
			w.WriteHeader(422)
			w.Write([]byte("Invalid credentials"))
			return
		}

		// Grab user details
		user, userFetchErr := conn.FetchUser(payload.Username)
		if userFetchErr != nil {
			log.Print(authErr)
			w.WriteHeader(422)
			w.Write([]byte("Invalid credentials"))
			return
		}

		token, tokenErr := transformer.VendToken(*user)
		if tokenErr != nil {
			// TODO: This is probably a pretty bad failure, handle with more urgency?
			w.WriteHeader(500)
			w.Write([]byte("Internal server error"))
			return
		}

		// TODO: Secure flag should be on by default but it'll mess with dev environments that don't have SSL
		http.SetCookie(w, &http.Cookie{
			Name:     "IndigoAuth",
			Value:    token,
			Path:     "/",
			MaxAge:   3600 * 10,
			Secure:   false,
			HttpOnly: true,
			//SameSite: http.SameSiteStrictMode,
		})
		w.WriteHeader(204)
	}
}
