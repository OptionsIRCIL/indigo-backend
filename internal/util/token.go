package util

import (
	"context"
	"net/http"

	"gorm.io/gorm"
	"myoptions.info/indigo/backend/model/entity"
	"myoptions.info/indigo/backend/model/schema"
)

const tokenKey = "token"

// StoreTokenToContext pushes a [schema.Token] into the request [context.Context].
// This is literally so basic, but I immensely dislike how there's no type-safety on the
// context and how it uses string keys, so this function exists to keep types and keys at
// least a little consistent.
func StoreTokenToContext(token *schema.Token, r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), tokenKey, token))
}

// FetchTokenFromContext fetches a schema.Token from the request context so long
// as it has been previously stored.
func FetchTokenFromContext(r *http.Request) *schema.Token {
	if r.Context().Value(tokenKey) == nil {
		// This only happens if you're a bad programmer, git gud
		panic("Failed to fetch *schema.Token from request context - token was not set")
	}

	return r.Context().Value(tokenKey).(*schema.Token)
}

// TokenToEmployee queries a database connection to get the full entity.Employee given
// a schema.Token.
func TokenToEmployee(token *schema.Token, database *gorm.DB, ctx context.Context) *entity.Employee {
	employee, err := gorm.G[*entity.Employee](database).Where("username = ?", token.Subject).First(ctx)
	if err != nil {
		panic("bruh")
	}
	return employee
}
