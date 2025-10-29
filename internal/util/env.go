package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	IndigoEnv      string
	IndigoSecret   string
	LdapSearchBase string
	LdapDomain     string
	LdapUrl        string
	LdapUsername   string
	LdapPassword   string
}

func requireEnv(envKey string) string {
	envValue := os.Getenv(envKey)
	if envValue == "" {
		log.Fatal("Missing required environment variable " + envKey + "! Exiting...")
	}
	return envValue
}

func envOrDefault(envKey string, defaultValue string) string {
	envValue := os.Getenv(envKey)
	if envValue == "" {
		envValue = defaultValue
	}
	return envValue
}

func LoadConfig() *Config {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal(envErr)
	}

	return &Config{
		IndigoEnv:      envOrDefault("INDIGO_ENV", "dev"),
		IndigoSecret:   requireEnv("INDIGO_SECRET"),
		LdapSearchBase: requireEnv("LDAP_SEARCH_BASE"),
		LdapDomain:     requireEnv("LDAP_DOMAIN"),
		LdapUrl:        requireEnv("LDAP_URL"),
		LdapUsername:   requireEnv("LDAP_USERNAME"),
		LdapPassword:   requireEnv("LDAP_PASSWORD"),
	}
}
