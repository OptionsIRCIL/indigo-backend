package controller

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"time"

	"gorm.io/gorm"
	"myoptions.info/indigo/backend/internal/config"
	"myoptions.info/indigo/backend/internal/service"
	u "myoptions.info/indigo/backend/internal/util"
	"myoptions.info/indigo/backend/internal/util/crypto"
	"myoptions.info/indigo/backend/model/entity"
)

type AuthenticationProvider func(credentials UserCredentials) (*entity.Employee, error)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthenticationError struct {
	Msg string
}

func (a AuthenticationError) Error() string {
	return a.Msg
}

func throwAuthenticationError(w http.ResponseWriter) {
	u.ThrowHttpError(w, 422, "Invalid credentials")
}

func firstEmployeeOrNew(username string, database *gorm.DB) (*entity.Employee, error) {
	ctx := context.Background()
	count, err := gorm.G[entity.Employee](database).Where("username = ?", username).Count(ctx, "username")
	if err != nil {
		return nil, err
	}
	if count > 1 {
		return nil, AuthenticationError{Msg: "Multiple Employees found for username, what???"}
	}

	// Pull and update if an employee is available
	if count == 1 {
		employee, err := gorm.G[entity.Employee](database).Where("username = ?", username).First(ctx)
		if err != nil {
			return nil, err
		}

		// TODO: Update employee first/last name

		return &employee, nil
	}

	// No employee, make a new one
	employee := &entity.Employee{Username: username}
	err = gorm.G[entity.Employee](database).Create(ctx, employee)
	if err != nil {
		return nil, err
	}
	return employee, nil
}

func AuthEntry(
	conn *service.LdapConnection,
	database *gorm.DB,
	sameSite string,
) http.HandlerFunc {
	// Determine SameSite value
	var sameSiteHeader http.SameSite
	// config.IndigoEnv == "dev"
	if true {
		sameSiteHeader = http.SameSiteNoneMode
	} else {
		switch sameSite {
		case "none":
			sameSiteHeader = http.SameSiteNoneMode
			break
		case "lax":
			sameSiteHeader = http.SameSiteLaxMode
			break
		default:
			sameSiteHeader = http.SameSiteStrictMode
		}
	}

	ldapAuthenticator := func(credentials UserCredentials) (*entity.Employee, error) {
		// LDAP enabled?
		if conn.Connection == nil {
			return nil, AuthenticationError{Msg: "Authenticator not enabled"}
		}

		// Parse out username and attempt bind auth using typical AD credential format
		fqUsername := config.Config.Authentication.Ldap.Domain + "\\" + credentials.Username
		authErr := conn.AttemptAuth(fqUsername, credentials.Password)
		if authErr != nil {
			log.Printf("Authentication error for user %s: %s", fqUsername, authErr)
			return nil, authErr
		}

		// Grab user details
		_, userFetchErr := conn.FetchUser(credentials.Username)
		if userFetchErr != nil {
			log.Printf("Unable to fetch details for user %s: %s", credentials.Username, userFetchErr)
			return nil, userFetchErr
		}

		return firstEmployeeOrNew(credentials.Username, database)
	}

	localAuthenticator := func(credentials UserCredentials) (*entity.Employee, error) {
		// Local enabled?
		if config.Config.Authentication.Local == nil {
			return nil, AuthenticationError{Msg: "Authenticator not enabled"}
		}

		// Grab user from database
		ctx := context.Background()
		user, userErr := gorm.G[entity.LocalUser](database).Where("username = ?", credentials.Username).First(ctx)
		if userErr != nil {
			return nil, userErr
		}

		// Verify password
		if !crypto.Verify(user.PasswordHash, credentials.Password) {
			return nil, AuthenticationError{Msg: "Credentials do not match"}
		}

		// Check if expired
		// TODO: cascade expiry date up to vended JWT if possible
		if user.ExpiresAt != nil && user.ExpiresAt.UnixMilli() < time.Now().UnixMilli() {
			return nil, AuthenticationError{Msg: "Account expired"}
		}

		return firstEmployeeOrNew(credentials.Username, database)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var payload UserCredentials

		// Decode payload
		decodeErr := u.DecodeJsonBody(w, r, &payload)
		if decodeErr != nil {
			if decodeErr.Status > 499 {
				u.ThrowHttpUnhandled(w, decodeErr)
			} else {
				// TODO: Anything that shouldn't be logged from DecodeJsonBody?
				u.ThrowHttpError(w, decodeErr.Status, decodeErr.Msg)
			}
			return
		}

		// 422 if the username has any invalid chars
		usernameCharsOk, regexErr := regexp.MatchString(`^[A-Za-z0-9-_. ]+$`, payload.Username)
		if regexErr != nil {
			u.ThrowHttpUnhandled(w, regexErr)
			return
		} else if !usernameCharsOk {
			throwAuthenticationError(w)
			return
		}

		employee, localErr := localAuthenticator(payload)
		if localErr != nil {
			ldapEmployee, ldapError := ldapAuthenticator(payload)
			if ldapError != nil {
				throwAuthenticationError(w)
				return
			}
			employee = ldapEmployee
		}

		token, tokenErr := crypto.VendToken(employee)
		if tokenErr != nil {
			u.ThrowHttpUnhandled(w, tokenErr)
			return
		}

		// TODO: Set domain field in cookie when prod?
		http.SetCookie(w, &http.Cookie{
			Name:   "IndigoAuth",
			Value:  token,
			Path:   "/",
			MaxAge: 3600 * 10,
			Secure:/*config.IndigoEnv != "dev"*/ true,
			HttpOnly: false,
			SameSite: sameSiteHeader,
		})

		u.ReturnSerialized(w, 200, *employee, []string{"get"})
		log.Printf("Authentication attempt successful for user: %s", payload.Username)
	}
}

func DeleteCookie() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sameSite http.SameSite
		if /*c.IndigoEnv == "dev"*/ true {
			sameSite = http.SameSiteNoneMode
		} else {
			sameSite = http.SameSiteStrictMode
		}

		http.SetCookie(w, &http.Cookie{
			Name:   "IndigoAuth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
			Secure:/*c.IndigoEnv != "dev"*/ true,
			HttpOnly: false,
			SameSite: sameSite,
		})
		w.WriteHeader(204)
	}
}
