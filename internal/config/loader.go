package config

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"log"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
)

const configLocation = "./config.json"

// ApplicationConfig is used to store the application config in a global accessible to
// all parts of the application.
type ApplicationConfig struct {
	Authentication *AuthenticationConfigNode `json:"authentication,omitempty" validate:"required"`
	Database       *DatabaseConfigNode       `json:"database,omitempty" validate:"required"`
}

type AuthenticationConfigNode struct {
	HmacKey string                         `json:"hmacKey,omitempty" validate:"required,min=64"`
	Ldap    *LdapAuthenticationConfigNode  `json:"ldap,omitempty"`
	Local   *LocalAuthenticationConfigNode `json:"local,omitempty"`
}

type LdapAuthenticationConfigNode struct {
	Url        string `json:"url,omitempty" validate:"required"`
	Domain     string `json:"domain,omitempty" validate:"required"`
	SearchBase string `json:"searchBase,omitempty" validate:"required"`
	Username   string `json:"username,omitempty" validate:"required"`
	Password   string `json:"password,omitempty" validate:"required"`
}

type LocalAuthenticationConfigNode struct {
}

type DatabaseConfigNode struct {
	Dsn string `json:"dsn,omitempty" validate:"required"`
}

func readConfig() *ApplicationConfig {
	// Check if in a unit test
	// TODO: Optimize out in production build?
	if testing.Testing() {
		// Populate with dummy data
		secret := make([]byte, 64)
		_, err := rand.Read(secret)
		if err != nil {
			log.Fatalln(err)
		}

		return &ApplicationConfig{
			Authentication: &AuthenticationConfigNode{
				HmacKey: string(secret),
			},
		}
	}

	f, err := os.Open(configLocation)
	if err != nil {
		log.Fatalln("Failed to open config")
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		log.Fatalln("Failed to read config")
	}

	config := &ApplicationConfig{}
	err = json.Unmarshal(contents, config)
	if err != nil {
		log.Fatalln("Failed to parse config", err)
	}

	v := validator.New()
	err = v.Struct(config)
	if err != nil {
		log.Fatalln("Config contains errors", err)
	}

	return config
}

// Config provides global access to application configuration.
// I personally will kill you if you modify this at runtime.
var Config = readConfig()
