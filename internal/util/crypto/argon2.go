package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"regexp"
	"strconv"

	"golang.org/x/crypto/argon2"
)

const argon2Memory = 64 * 1024
const argon2Threads = 4
const argon2KeyLen = 32
const argon2Time = 3
const argon2SaltLen = 32

var argon2Header = "$argon2id$v=" + strconv.Itoa(argon2.Version) +
	"$m=" + strconv.Itoa(argon2Memory/1024) +
	",t=" + strconv.Itoa(argon2Time) +
	",p=" + strconv.Itoa(argon2Threads)

var expr = regexp.MustCompile("^\\$argon2id\\$v=(?P<version>\\d+)\\$m=(?P<memory>\\d+),t=(?P<iterations>\\d+),p=(?P<parallelism>\\d+)\\$(?P<salt>.+)\\$(?P<hash>.+)")

func genSalt() ([]byte, error) {
	s := make([]byte, argon2SaltLen)
	_, err := rand.Read(s)
	return s, err
}

// Hash generates an argon2id hash given a password.
func Hash(password string) (string, error) {
	salt, saltErr := genSalt()
	if saltErr != nil {
		return "", saltErr
	}

	return hashWithSalt(password, salt), nil
}

// TODO: This function seems to be spitting out nonstandard encoded strings. Fix that!
func hashWithSalt(password string, salt []byte) string {
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLen,
	)

	return argon2Header + "$" + base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(hash)
}

// Verify compares a stored hash to a provided password for password authentication.
func Verify(hash string, password string) bool {
	match := expr.FindStringSubmatch(hash)
	if match == nil {
		return false
	}

	if len(match) != 7 {
		return false
	}

	salt, decodeErr := base64.RawStdEncoding.DecodeString(match[5])
	if decodeErr != nil {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(hashWithSalt(password, salt)), []byte(hash)) == 1
}
