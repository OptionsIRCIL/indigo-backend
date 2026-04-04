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

const defaultConfigLocation = "./config.json"
const defaultMaxFileSize = 1e8
const defaultAttachmentDirectory = "/srv/indigo"
const defaultFallbackAttachmentDirectory = "./attachments"

var defaultPermissibleMimeTypes = []string{
	"application/msword",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"application/vnd.oasis.opendocument.text",
	"application/vnd.ms-excel",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"application/vnd.oasis.opendocument.spreadsheet",
	"application/vnd.ms-powerpoint",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"application/vnd.oasis.opendocument.presentation",
	"application/pdf",
	"image/png",
	"image/jpeg",
}

// ApplicationConfig is used to store the application config in a global accessible to
// all parts of the application.
type ApplicationConfig struct {
	Authentication *AuthenticationConfigNode `json:"authentication,omitempty" validate:"required"`
	Database       *DatabaseConfigNode       `json:"database,omitempty" validate:"required"`
	Attachments    *AttachmentConfigNode     `json:"attachments,omitempty"`
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

type AttachmentConfigNode struct {
	Directory            string   `json:"string,omitempty" validate:"dir"`
	PermissibleMimeTypes []string `json:"permissibleMimeTypes,omitempty"`
	MaxFileSize          uint     `json:"maxFileSize,omitempty"`
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
			Attachments: &AttachmentConfigNode{
				MaxFileSize:          defaultMaxFileSize,
				PermissibleMimeTypes: defaultPermissibleMimeTypes,
			},
		}
	}

	// Check if a config location has been provided via an environment variable
	configLocation := os.Getenv("INDIGO_CONFIG_LOCATION")
	if configLocation == "" {
		configLocation = defaultConfigLocation
	}

	f, err := os.Open(configLocation)
	if err != nil {
		wdir, _ := os.Getwd()
		log.Fatalf("Failed to open config.\nConfig Location: \"%s\"\nWorking Dir: \"%s\"", configLocation, wdir)
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

	// Add any defaults
	if config.Attachments == nil {
		config.Attachments = &AttachmentConfigNode{}
	}

	if config.Attachments.Directory == "" {
		dir := defaultAttachmentDirectory
		dirStat, dirErr := os.Stat(dir)

		if os.IsNotExist(dirErr) {
			dir = defaultFallbackAttachmentDirectory
			dirStat, dirErr = os.Stat(dir)
			if os.IsNotExist(dirErr) {
				createErr := os.Mkdir(defaultFallbackAttachmentDirectory, 0770)
				if createErr != nil {
					log.Fatalln("Failed to create an attachments directory at " + defaultAttachmentDirectory + " and at fallback " + defaultFallbackAttachmentDirectory)
				}
				dirStat, dirErr = os.Stat(defaultFallbackAttachmentDirectory)
			}
		}
		if dirErr != nil {
			log.Fatalln("Failed to open attachments directory: ", dirErr)
		}
		if !dirStat.IsDir() {
			log.Fatalln("Default directory location " + dir + " is not a directory")
		}

		config.Attachments.Directory = dir
	}

	if config.Attachments.MaxFileSize == 0 {
		// Default to 10MB
		config.Attachments.MaxFileSize = defaultMaxFileSize
	}

	if config.Attachments.PermissibleMimeTypes == nil {
		config.Attachments.PermissibleMimeTypes = defaultPermissibleMimeTypes
	}

	return config
}

// Config provides global access to application configuration.
// I personally will kill you if you modify this at runtime.
var Config = readConfig()
