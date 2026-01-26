//go:build prod

package entry

import (
	"flag"
	"log"
	"os"

	"myoptions.info/indigo/backend/internal/routine"
	"myoptions.info/indigo/backend/internal/util"
)

func Entry(l *log.Logger) {
	flags := util.ServeRuntimeFlags{}
	set := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	set.IntVar(&flags.Port, "port", 8080, "specifies the port the http server runs on")
	set.StringVar(&flags.Socket, "socket", "", "specifies a socket to listen on, takes priority over -port")
	set.IntVar(&flags.SocketUid, "socket_uid", -1, "if desired, change the owning UID on the listening socket")
	set.IntVar(&flags.SocketGid, "socket_gid", -1, "if desired, change the owning GID on the listening socket")
	set.BoolVar(&flags.AllowInsecureLdap, "allow_insecure_ldap", false, "allow insecure connections to LDAP")
	set.StringVar(&flags.AuthSameSite, "auth_same_site", "", "configure the SameSite attribute of returned cookies")

	if err := set.Parse(os.Args[1:]); err != nil {
		l.Fatalf("Could not parse arguments: %s\n", err)
	}

	os.Exit(routine.RunServe("prod", flags))
}
