package schema

import "github.com/golang-jwt/jwt/v5"

// A Token contains any relevant user data that should be stored in the cookie.
// This gets encoded and signed as a JWT.
type Token struct {
	FirstName string           `json:"given_name"`
	LastName  string           `json:"family_name"`
	Groups    []string         `json:"groups"`
	Issuer    string           `json:"iss"`
	Subject   string           `json:"sub"`
	Audience  []string         `json:"aud"`
	ExpiresAt *jwt.NumericDate `json:"exp"`
	NotBefore *jwt.NumericDate `json:"nbf"`
	IssuedAt  *jwt.NumericDate `json:"iat"`
}

func (c Token) GetExpirationTime() (*jwt.NumericDate, error) {
	return c.ExpiresAt, nil
}

func (c Token) GetIssuedAt() (*jwt.NumericDate, error) {
	return c.IssuedAt, nil
}

func (c Token) GetNotBefore() (*jwt.NumericDate, error) {
	return c.NotBefore, nil
}

func (c Token) GetIssuer() (string, error) {
	return c.Issuer, nil
}

func (c Token) GetSubject() (string, error) {
	return c.Subject, nil
}

func (c Token) GetAudience() (jwt.ClaimStrings, error) {
	return c.Audience, nil
}
