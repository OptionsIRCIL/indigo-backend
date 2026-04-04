package middleware

import (
	"net/http"

	"myoptions.info/indigo/backend/internal/service"
	"myoptions.info/indigo/backend/internal/util"
	"myoptions.info/indigo/backend/internal/util/crypto"
)

func RequireAuth(l *service.LdapConnection, next http.HandlerFunc) http.HandlerFunc {
	// Return a function wrapping the next handler
	return func(w http.ResponseWriter, r *http.Request) {
		// Grab IndigoAuth cookie
		cookies := r.CookiesNamed("IndigoAuth")
		if len(cookies) != 1 {
			util.ThrowHttpStatus(w, 401)
			return
		}

		// Parse cookie
		token, err := crypto.StringToToken(cookies[0].Value)
		if err != nil {
			util.ThrowHttpStatus(w, 403)
			return
		}

		// Verify iat > last password modification time
		// TODO: Rework to support both LDAP and local expiry
		/*lastModifiedTime, ldapErr := l.FetchPwdLastSet(user.Username)
		if ldapErr != nil {
			util.ThrowHttpUnhandled(w, ldapErr)
			return
		}
		if iat < lastModifiedTime {
			// More descriptive/useful error could probably be returned but idk
			// For now, we'll just return a 401.
			util.ThrowHttpStatus(w, 401)
			return
		}*/

		// Add token to context and continue
		next.ServeHTTP(w, util.StoreTokenToContext(token, r))
	}
}
