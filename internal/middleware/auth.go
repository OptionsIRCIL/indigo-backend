package middleware

import (
	"backend/internal/service"
	"backend/internal/util"
	"context"
	"net/http"
)

func RequireAuth(jwtTransformer *service.JwtTransformer, next http.HandlerFunc) http.HandlerFunc {
	// Return a function wrapping the next handler
	return func(w http.ResponseWriter, r *http.Request) {
		// Grab IndigoAuth cookie
		cookies := r.CookiesNamed("IndigoAuth")
		if len(cookies) != 1 {
			util.ThrowHttpStatus(w, 401)
			return
		}

		// Parse cookie
		user, err := jwtTransformer.ValidateToken(cookies[0].Value)
		if err != nil {
			util.ThrowHttpStatus(w, 403)
			return
		}

		// Add user to context and continue
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
