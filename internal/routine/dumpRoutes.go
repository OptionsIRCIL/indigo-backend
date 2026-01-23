package routine

import (
	"fmt"

	c "myoptions.info/indigo/backend/internal/config"
	s "myoptions.info/indigo/backend/internal/service"
	"myoptions.info/indigo/backend/internal/util"
)

// RunDumpRoutes initializes some dummy services and dumps all configured routes from config.CreateMux.
// Utilized during testing to verify that all configured routes are documented in the OpenAPI spec.
func RunDumpRoutes() int {
	// Create routes using MuxWrapper. Provide dummy services.
	mux := c.CreateMux(
		c.Services{
			Config: &util.Config{},
			Ldap:   &s.LdapConnection{},
			Jwt:    &s.JwtTransformer{},
		},
	)

	fmt.Println(mux.DumpRoutes())
	return 0
}
