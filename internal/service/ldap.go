package service

import (
	"crypto/tls"
	"fmt"
	"log"
	"regexp"

	"github.com/go-ldap/ldap/v3"
)

// An LdapConnection holds a variety of parameters for an [ldap.Conn], as well as holding a pointer to
// an [ldap.Conn] created from these parameters.
type LdapConnection struct {
	// The active LDAP connection
	Connection *ldap.Conn

	// The search base for user queries.
	Base string

	// The domain prefix used when authenticating users with bind auth.
	Domain string

	url    string
	secure bool
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

	// TODO: Recover if err

	return err
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

// FetchUser fetches basic user details from LDAP given a username.
func (c *LdapConnection) FetchUser(username string) (*LdapUser, error) {
	err := c.ensureConnection()
	if err != nil {
		return nil, err
	}

	// "%s@*" is used rather than just "%s" since the current userPrincipalNames are in the format "username@email.com".
	searchRequest := ldap.NewSearchRequest(
		c.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(userPrincipalName=%s@*))", ldap.EscapeDN(username)),
		[]string{"givenName", "sn", "memberOf", "cn"},
		nil,
	)

	searchResult, searchErr := c.Connection.Search(searchRequest)
	if searchErr != nil {
		return nil, searchErr
	} else if len(searchResult.Entries) != 1 {
		return nil, &LdapError{"Multiple entries returned by user query"}
	}

	entry := searchResult.Entries[0]

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

// Initialize does what it says on the package. Given a username and password, the server specified in
// [LdapConnection] is dialed and subsequently bound to.
func (c *LdapConnection) Initialize(username string, password string) error {
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

	err = c.Connection.Bind(username, password)
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
