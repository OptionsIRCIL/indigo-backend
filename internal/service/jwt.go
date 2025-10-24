package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtTransformer struct {
	secret []byte
}

type JwtFailure struct {
	Msg   string
	Fatal bool
}

type claimSet struct {
	FirstName string   `json:"given_name"`
	LastName  string   `json:"family_name"`
	Email     string   `json:"email"`
	Groups    []string `json:"groups"`
	jwt.RegisteredClaims
}

func (j *JwtFailure) Error() string {
	return j.Msg
}

func (j *JwtTransformer) SetSecret(secret []byte) error {
	if len(secret) < 32 {
		return &JwtFailure{"Secret too short!", true}
	}

	j.secret = secret
	return nil
}

func (j *JwtTransformer) VendToken(user LdapUser) (string, error) {
	if len(j.secret) < 32 {
		return "", &JwtFailure{"Secret too short!", true}
	}

	// https://www.iana.org/assignments/jwt/jwt.xhtml#claims
	claims := &claimSet{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Groups:    user.Groups,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "indigo-backend",
			Subject:   user.Username,
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
}

func (j *JwtTransformer) ValidateToken(encodedToken string) (*LdapUser, error) {
	token, tokenErr := jwt.ParseWithClaims(
		encodedToken,
		&claimSet{},
		func(token *jwt.Token) (any, error) {
			return j.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuer("indigo-backend"),
		jwt.WithIssuedAt(),
	)
	if tokenErr != nil {
		return nil, tokenErr
	}

	claims, ok := token.Claims.(*claimSet)
	if !ok {
		return nil, &JwtFailure{
			Msg:   "Claim extraction failed",
			Fatal: false,
		}
	}

	// TODO: Compare jwt iat to the last user password modification date

	return &LdapUser{
		FirstName: claims.FirstName,
		LastName:  claims.LastName,
		Username:  claims.Subject,
		Email:     claims.Email,
		Groups:    claims.Groups,
	}, nil
}
