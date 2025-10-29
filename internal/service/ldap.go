package service

import (
	"crypto/tls"
	"fmt"
	"log"
	"regexp"

	"github.com/go-ldap/ldap/v3"
)

type LdapConnection struct {
	Connection *ldap.Conn
	Base       string
	Domain     string
	url        string
	secure     bool
}

type LdapUser struct {
	FirstName string
	LastName  string
	Username  string
	Email     string
	Groups    []string
}

type LdapError struct {
	Msg string
}

func (e *LdapError) Error() string {
	return e.Msg
}

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

func (c *LdapConnection) FetchUser(username string) (*LdapUser, error) {
	searchRequest := ldap.NewSearchRequest(
		c.Base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(userPrincipalName=%s@*))", ldap.EscapeDN(username)),
		[]string{"givenName", "mail", "sn", "memberOf", "cn"},
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
		Email:     entry.GetAttributeValue("mail"),
		Groups:    groups,
	}, nil
}

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

func (c *LdapConnection) SetUrl(url string) {
	c.url = url
}

func (c *LdapConnection) SetSecure(v bool) {
	c.secure = v
}
