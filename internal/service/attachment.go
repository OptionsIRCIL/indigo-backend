package service

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"myoptions.info/indigo/backend/internal/config"
)

// AttachmentError is a generic error for issues encountered when handling attachment storage or retrieval.
type AttachmentError struct {
	Msg string
}

func (a AttachmentError) Error() string {
	return a.Msg
}

// filenameValid runs a regex on a passed filename to determine if it contains
// any disallowed characters. It then picks out a file extension from the string.
// https://stackoverflow.com/a/31976060
func filenameValid(filename string) (string, bool) {
	if len(filename) > 256 {
		return "", false
	}

	expr := regexp.MustCompile(`^[^\n\x00/\\<>:"|?*.~][^\n\x00/\\<>:"|?*]*(?P<ext>\.[a-zA-Z0-9]+)$`)
	match := expr.FindStringSubmatch(filename)
	if len(match) != 2 {
		return "", false
	}

	return match[1], true
}

func getContentType(content []byte) string {
	return strings.Split(http.DetectContentType(content), ";")[0]
}

// StoreFile checks an uploaded file for safety (file name, mime type, extension) and
// writes it to the local file system. id is used as the file name and the extension
// is omitted from the destination file name.
func StoreFile(id string, name string, contents []byte) (string, error) {
	extension, valid := filenameValid(name)
	if !valid {
		return "", AttachmentError{"Filename invalid"}
	}

	mimetype := getContentType(contents)
	if !slices.Contains(config.Config.Attachments.PermissibleMimeTypes, mimetype) {
		return "", AttachmentError{"Invalid mimetype"}
	}

	validExtensions, mimeErr := mime.ExtensionsByType(mimetype)
	if mimeErr != nil {
		return "", mimeErr
	}

	if !slices.Contains(validExtensions, extension) {
		return "", AttachmentError{"Extension mismatch"}
	}

	// TODO: collision testing? extremely unlikely that the same UUID would be drawn twice, though
	writeErr := os.WriteFile(
		filepath.Join(config.Config.Attachments.Directory, id),
		contents,
		0660,
	)
	return mimetype, writeErr
}

// RetrieveFile fetches a file from the filesystem given its ID.
func RetrieveFile(id string) ([]byte, error) {
	// TODO: hook in with gorm?
	return os.ReadFile(filepath.Join(config.Config.Attachments.Directory, id))
}
