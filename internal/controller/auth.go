package controller

import (
	"log"
	"net/http"
	"regexp"

	"myoptions.info/indigo/backend/internal/service"
	u "myoptions.info/indigo/backend/internal/util"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func throwAuthenticationError(w http.ResponseWriter) {
	u.ThrowHttpError(w, 422, "Invalid credentials")
}

func AuthEntry(
	config *u.Config,
	conn *service.LdapConnection,
	transformer *service.JwtTransformer,
	sameSite string,
) http.HandlerFunc {
	// Determine SameSite value
	var sameSiteHeader http.SameSite
	if config.IndigoEnv == "dev" {
		sameSiteHeader = http.SameSiteNoneMode
	} else {
		switch sameSite {
		case "none":
			sameSiteHeader = http.SameSiteNoneMode
			break
		case "lax":
			sameSiteHeader = http.SameSiteLaxMode
			break
		default:
			sameSiteHeader = http.SameSiteStrictMode
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var payload UserCredentials

		// Decode payload
		decodeErr := u.DecodeJsonBody(w, r, &payload)
		if decodeErr != nil {
			if decodeErr.Status > 499 {
				u.ThrowHttpUnhandled(w, decodeErr)
			} else {
				// TODO: Anything that shouldn't be logged from DecodeJsonBody?
				u.ThrowHttpError(w, decodeErr.Status, decodeErr.Msg)
			}
			return
		}

		// 422 if the username has any invalid chars
		usernameCharsOk, regexErr := regexp.MatchString(`^[A-Za-z0-9-_. ]+$`, payload.Username)
		if regexErr != nil {
			u.ThrowHttpUnhandled(w, regexErr)
			return
		} else if !usernameCharsOk {
			throwAuthenticationError(w)
			return
		}

		// Parse out username and attempt bind auth using typical AD credential format
		fqUsername := config.LdapDomain + "\\" + payload.Username
		authErr := conn.AttemptAuth(fqUsername, payload.Password)
		if authErr != nil {
			log.Printf("Authentication error for user %s: %s", fqUsername, authErr)
			throwAuthenticationError(w)
			return
		}

		// Grab user details
		user, userFetchErr := conn.FetchUser(payload.Username)
		if userFetchErr != nil {
			log.Printf("Unable to fetch details for user %s: %s", payload.Username, userFetchErr)
			throwAuthenticationError(w)
			return
		}

		token, tokenErr := transformer.VendToken(*user)
		if tokenErr != nil {
			u.ThrowHttpUnhandled(w, tokenErr)
			return
		}

		// TODO: Set domain field in cookie when prod?
		http.SetCookie(w, &http.Cookie{
			Name:     "IndigoAuth",
			Value:    token,
			Path:     "/",
			MaxAge:   3600 * 10,
			Secure:   config.IndigoEnv != "dev",
			HttpOnly: false,
			SameSite: sameSiteHeader,
		})
		w.WriteHeader(204)

		log.Printf("Authentication attempt successful for user: %s", payload.Username)
	}
}

func DeleteCookie(c *u.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sameSite http.SameSite
		if c.IndigoEnv == "dev" {
			sameSite = http.SameSiteNoneMode
		} else {
			sameSite = http.SameSiteStrictMode
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "IndigoAuth",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Secure:   c.IndigoEnv != "dev",
			HttpOnly: false,
			SameSite: sameSite,
		})
		w.WriteHeader(204)
	}
}
