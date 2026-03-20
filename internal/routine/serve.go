package routine

import (
	"fmt"
	"log"

	"myoptions.info/indigo/backend/internal/config"
	c "myoptions.info/indigo/backend/internal/config/routes"
	s "myoptions.info/indigo/backend/internal/service"
	"myoptions.info/indigo/backend/internal/util"
)

// RunServe serves the application proper. Initializes all services and listens on either a port or a socket.
func RunServe(flags util.ServeRuntimeFlags) int {
	// Populate LdapConnection (If applicable)
	l := s.LdapConnection{}
	if config.Config.Authentication.Ldap != nil {
		l.Base = config.Config.Authentication.Ldap.SearchBase
		l.Domain = config.Config.Authentication.Ldap.Domain
		l.SetUrl(config.Config.Authentication.Ldap.Url)
		l.SetCredentials(
			config.Config.Authentication.Ldap.Username,
			config.Config.Authentication.Ldap.Password,
		)

		// Initialize connection
		err := l.Initialize()
		if err != nil {
			log.Fatal(err)
		}
		//l.SetSecure((config.IndigoEnv == "prod") && !flags.AllowInsecureLdap)
		l.SetSecure(false)
		defer l.Connection.Close()
	}

	// Create routes using MuxWrapper
	mux := c.CreateMux(
		c.Services{
			Ldap:     &l,
			Flags:    flags,
			Database: util.ConnectToDatabase(),
		},
	)

	// Serve
	if flags.Socket == "" {
		log.Printf("Serving on :%d\n", flags.Port)
		log.Fatal(mux.ListenAndServe(fmt.Sprintf(":%d", flags.Port)))
	} else {
		log.Printf("Serving on socket %s\n", flags.Socket)
		log.Fatal(mux.ServeToSocket(flags.Socket, flags.SocketUid, flags.SocketGid))
	}

	return 0
}
