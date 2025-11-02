package service

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/go-ldap/ldap/v3"
)

// unixTimeJan11601 is the reason Microsoft makes me want to eat a slug.
// Jan 1st, 1601 in Unix Millis.
const unixTimeJan11601 = -11644473600000

// An LdapConnection holds a variety of parameters for an [ldap.Conn], as well as holding a pointer to
// an [ldap.Conn] created from these parameters.
type LdapConnection struct {
	// The active LDAP connection
	Connection *ldap.Conn

	// The search base for user queries.
	Base string

	// The domain prefix used when authenticating users with bind auth.
	Domain string

	url      string
	secure   bool
	username string
	password string
}

// An LdapUser is a struct carrying the core information needed to describe a user that exists in AD.
type LdapUser struct {
	FirstName string
	LastName  string
	Username  string
	Groups    []string
}

// An LdapError is used for any custom logic outside of [ldap].
type LdapError struct {
	Msg string
}

// Error implements the required method to be an error.
func (e *LdapError) Error() string {
	return e.Msg
}

func (c *LdapConnection) ensureConnection() error {
	// Run whoami to see if we're connected
	_, err := c.Connection.WhoAmI(nil)

	if err != nil {
		// Cast as ldap.Error
		var ldapErr *ldap.Error
		if errors.As(err, &ldapErr) {
			// What's the current error code?
			if ldapErr.ResultCode == 200 {
				// Connection was closed, probably due to inactivity. Try re-initializing.
				log.Println("LDAP connection closed! Attempting reconnect...")
				_ = c.Close()
				bindErr := c.Initialize()
				return bindErr
			}
		}

		// Failure :(
		return err
	}

	return nil
}

func (c *LdapConnection) userByUsername(username string, attributes []string) (*ldap.Entry, error) {
	// "%s@*" is used rather than just "%s" since the current userPrincipalNames are in the format "username@email.com".
	searchRequest := ldap.NewSearchRequest(
		c.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(userPrincipalName=%s@*))", ldap.EscapeDN(username)),
		attributes,
		nil,
	)

	searchResult, searchErr := c.Connection.Search(searchRequest)
	if searchErr != nil {
		return nil, searchErr
	} else if len(searchResult.Entries) != 1 {
		return nil, &LdapError{"Multiple entries returned by user query"}
	}

	return searchResult.Entries[0], nil
}

// AttemptAuth attempts to bind to an LDAP server with the supplied username and password.
// If connection fails, an error is returned. If connection succeeds, the connection is immediately
// closed and nil is returned.
func (c *LdapConnection) AttemptAuth(username string, password string) error {
	conn, err := ldap.DialURL(
		c.url,
		ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: !c.secure}),
	)
	if err != nil {
		return err
	}

	err = conn.Bind(username, password)
	if err != nil {
		return err
	}

	err = conn.Close()
	if err != nil {
		return err
	}

	return nil
}

// FetchPwdLastSet fetches the time a given user's password was changed in Unix millis.
func (c *LdapConnection) FetchPwdLastSet(username string) (int64, error) {
	err := c.ensureConnection()
	if err != nil {
		return -1, err
	}

	entry, searchErr := c.userByUsername(username, []string{"pwdLastSet"})
	if searchErr != nil {
		return -1, searchErr
	}

	// Microsoft refuses to do anything standard, so now we've gotta convert the time to unix millis.
	// https://learn.microsoft.com/en-us/windows/win32/adschema/a-pwdlastset
	passwordTime, convErr := strconv.ParseInt(entry.GetAttributeValue("pwdLastSet"), 10, 64)
	if convErr != nil {
		return -1, convErr
	}

	return passwordTime*10000 + unixTimeJan11601, nil
}

// FetchUser fetches basic user details from LDAP given a username.
func (c *LdapConnection) FetchUser(username string) (*LdapUser, error) {
	err := c.ensureConnection()
	if err != nil {
		return nil, err
	}

	entry, searchErr := c.userByUsername(username, []string{"givenName", "sn", "memberOf", "cn"})
	if searchErr != nil {
		return nil, searchErr
	}

	// Parse groups
	// TODO: Better group handling
	groupExpr := regexp.MustCompile(`^[Cc][Nn]=(.+?),.*$`)

	groups := make([]string, 0)
	for _, v := range entry.GetAttributeValues("memberOf") {
		groupMatch := groupExpr.FindStringSubmatch(v)
		if len(groupMatch) > 0 {
			groups = append(groups, groupMatch[1])
		}
	}

	// Assemble user
	return &LdapUser{
		Username:  entry.GetAttributeValue("cn"),
		FirstName: entry.GetAttributeValue("givenName"),
		LastName:  entry.GetAttributeValue("sn"),
		Groups:    groups,
	}, nil
}

// Initialize does what it says on the package. The server specified in [LdapConnection] is dialed and bound to.
func (c *LdapConnection) Initialize() error {
	var err error
	// TODO: this should fail cert check in prod if the cert is bad, but it's not doing that here for some reason.
	// Works fine in AttemptAuth, confusingly
	c.Connection, err = ldap.DialURL(
		c.url,
		ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: !c.secure}),
	)
	if err != nil {
		return err
	}

	log.Print("Contacted LDAP server. Attempting bind.")

	err = c.Connection.Bind(c.username, c.password)
	if err != nil {
		return err
	}

	log.Print("Bound to LDAP server.")

	return nil
}

// SetUrl is a setter for the LDAP server's URL.
func (c *LdapConnection) SetUrl(url string) {
	c.url = url
}

// SetSecure is used to enable or disable SSL certificate verification.
func (c *LdapConnection) SetSecure(v bool) {
	c.secure = v
}

// Close wraps [ldap.Conn.Close].
func (c *LdapConnection) Close() error {
	return c.Connection.Close()
}

// SetCredentials sets the username and password used for the LDAP service account.
func (c *LdapConnection) SetCredentials(username string, password string) {
	c.username = username
	c.password = password
}
