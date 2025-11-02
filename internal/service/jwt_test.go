package service

import (
	"fmt"
	"testing"
)

const superDuperSecretKey = "coffeecoffeecoffeecoffeecoffeecoffeecoffeecoffeecoffeecoffeecoff"

func buildExampleToken() string {
	j := JwtTransformer{}
	j.SetSecret([]byte(superDuperSecretKey))

	token, _ := j.VendToken(LdapUser{
		FirstName: "John",
		LastName:  "Doe",
		Username:  "john.doe",
		Groups:    []string{},
	})

	return token
}

var signedExampleToken = buildExampleToken()

func ExampleJwtTransformer_VendToken() {
	j := JwtTransformer{}
	j.SetSecret([]byte(superDuperSecretKey))

	token, err := j.VendToken(LdapUser{
		FirstName: "John",
		LastName:  "Doe",
		Username:  "john.doe",
		Groups:    []string{},
	})

	if err == nil {
		fmt.Printf("Token vend succeeded! Token=%s...\n", token[:16])
	} else {
		fmt.Println("Token vend failed...")
	}

	// Output: Token vend succeeded! Token=eyJhbGciOiJIUzUx...
}

func ExampleJwtTransformer_ValidateToken() {
	j := JwtTransformer{}
	j.SetSecret([]byte(superDuperSecretKey))

	// Example has been placed in a global variable for brevity
	exampleToken := signedExampleToken // eyJhbGciOiJIUzUx...

	user, _, err := j.ValidateToken(exampleToken)

	if err == nil {
		fmt.Printf("Token verification succeeded! Username=%s\n", user.Username)
	} else {
		fmt.Println("Token verification failed...")
	}
	// Output: Token verification succeeded! Username=john.doe
}

func TestJwtTransformer_ValidateToken(t *testing.T) {
	jGoodSecret := JwtTransformer{}
	err := jGoodSecret.SetSecret([]byte(superDuperSecretKey))
	if err != nil {
		t.Error("jGoodSecret.SetSecret failed unexpectedly")
	}

	jShortSecret := JwtTransformer{}
	err = jShortSecret.SetSecret([]byte("too_short"))
	if err == nil {
		t.Error("jGoodSecret.ShortSecret succeeded unexpectedly")
	}

	jWrongSecret := JwtTransformer{}
	err = jWrongSecret.SetSecret([]byte("different_512_bit_secret_in_string_formdifferent_512_bit_secret_in_string_form"))
	if err != nil {
		t.Error("jWrongSecret.SetSecret failed unexpectedly")
	}

	_, _, err = jGoodSecret.ValidateToken(signedExampleToken)
	if err != nil {
		t.Error("jGoodSecret.ValidateToken failed unexpectedly")
	}

	_, _, err = jShortSecret.ValidateToken(signedExampleToken)
	if err == nil {
		t.Error("jShortSecret.ValidateToken succeeded unexpectedly")
	}

	_, _, err = jWrongSecret.ValidateToken(signedExampleToken)
	if err == nil {
		t.Error("jWrongSecret.ValidateToken succeeded unexpectedly")
	}
}
