package crypto

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"myoptions.info/indigo/backend/internal/config"
	"myoptions.info/indigo/backend/model/entity"
	"myoptions.info/indigo/backend/model/schema"
)

// A JwtFailure is returned when a method of [JwtTransformer] fails in either a fatal or non-fatal fashion.
type JwtFailure struct {
	// A description of the encountered error
	Msg string

	// If the error should be treated as fatal
	Fatal bool
}

// Error is implemented for conformity as an error.
func (j *JwtFailure) Error() string {
	return j.Msg
}

// VendToken vends a new JWT given an [LdapUser]. The returned string is a signed JWT that is
// safe to return to the user.
func VendToken(employee *entity.Employee) (string, error) {
	// https://www.iana.org/assignments/jwt/jwt.xhtml#claims
	claims := &schema.Token{
		Subject:   employee.Username,
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
		Groups:    []string{},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "indigo-backend",
		Audience:  []string{"indigo-backend", "indigo-frontend"},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(config.Config.Authentication.HmacKey))
}

// StringToToken takes signed token string and attempts to verify its signature and scan
// its claims into a [schema.Token]. [jwt.WithValidMethods] is employed to ensure that the only
// accepted method is HS512. iat and exp claims are also required.
func StringToToken(encodedToken string) (*schema.Token, error) {
	token, tokenErr := jwt.ParseWithClaims(
		encodedToken,
		&schema.Token{},
		func(token *jwt.Token) (any, error) {
			return []byte(config.Config.Authentication.HmacKey), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Name}),
		jwt.WithIssuer("indigo-backend"),
		jwt.WithAudience("indigo-backend"),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	)
	if tokenErr != nil {
		return nil, tokenErr
	}

	claims, ok := token.Claims.(*schema.Token)
	if !ok {
		return nil, &JwtFailure{
			Msg:   "Claim extraction failed",
			Fatal: false,
		}
	}
	return claims, nil
}
