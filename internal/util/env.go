package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// A Config carries all relevant environment variables that are used across the application.
// By fetching these variables ahead of time, we can avoid fetching environment variables
// at the time of a request.
type Config struct {
	// The target environment type. Should be either "prod" or "dev". This variable is intended to activate
	// or deactivate various application features that require additional configuration to run, thus making
	// development environment configuration easier while retaining these features in production.
	IndigoEnv string

	// The HMAC256 secret to use for signing JWTs. Should be at least 32 characters.
	IndigoSecret string

	// MariaDB username
	DbUser string

	// MariaDB password
	DbPassword string

	// MariaDB host (e.g., 127.0.0.1)
	DbHost string

	// MariaDB port (e.g., 3306)
	DbPort string

	// The search base to use when fetching user details from LDAP.
	LdapSearchBase string

	// The Active Directory domain users of the application belong to. I.E., for a user ORG\john.doe,
	// the domain is ORG.
	LdapDomain string

	// The URL of your Active Directory's LDAP. Should be of format "protocol://host:port".
	LdapUrl string

	// The username of an LDAP service account with the ability to read domain user details.
	LdapUsername string

	// The password of an LDAP service account with the ability to read domain user details.
	LdapPassword string

	// template database var
	GormDb string
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

// LoadConfig loads all relevant environment variables into a [Config]. If any variables are found to be missing or
// invalid, [log.Fatal] is called to terminate the application.
func LoadConfig() *Config {
	// Ignore .env load fail - User may want to specify only via typical environment vars
	_ = godotenv.Load()

	return &Config{
		IndigoEnv:    envOrDefault("INDIGO_ENV", "dev"),
		IndigoSecret: requireEnv("INDIGO_SECRET"),

		DbUser:     requireEnv("DB_USER"),
		DbPassword: envOrDefault("DB_PASSWORD", ""),      // Password blank on local setup
		DbHost:     envOrDefault("DB_HOST", "127.0.0.1"), // Default to localhost
		DbPort:     envOrDefault("DB_PORT", "3306"),      // Default to MariaDB standard port
		GormDb:     requireEnv("GORM_DB"),

		LdapSearchBase: requireEnv("LDAP_SEARCH_BASE"),
		LdapDomain:     requireEnv("LDAP_DOMAIN"),
		LdapUrl:        requireEnv("LDAP_URL"),
		LdapUsername:   requireEnv("LDAP_USERNAME"),
		LdapPassword:   requireEnv("LDAP_PASSWORD"),
	}
}
