//go:debug x509negativeserial=1
package main

import (
	"flag"
	"log"
	"os"
	"slices"

	"myoptions.info/indigo/backend/internal/routine"
	"myoptions.info/indigo/backend/internal/util"
)

func main() {
	// Use a logger with no prefix for program startup
	l := log.New(os.Stderr, "", 0)
	l.Println("Indigo CIL, v0.0.0")

	// Check if valid subcommand exists
	if len(os.Args) < 2 || !slices.Contains([]string{"serve", "dump_routes"}, os.Args[1]) {
		// Drill down exact error to be a little more helpful
		var preciseError string
		if len(os.Args) < 2 {
			preciseError = "Error: No subcommand supplied."
		} else {
			preciseError = "Error: Unknown subcommand."
		}

		l.Fatalln(
			preciseError,
			`
Subcommands available:
  serve:
        Serves backend application on a port or Unix socket. View further options by passing the -help flag.
  dump_routes:
        Dumps configured routes to stdout.`,
		)
	}

	// Execute subcommand
	switch os.Args[1] {
	case "serve":
		flags := util.ServeRuntimeFlags{}
		set := flag.NewFlagSet("serve", flag.ExitOnError)
		set.IntVar(&flags.Port, "port", 8080, "specifies the port the http server runs on")
		set.StringVar(&flags.Socket, "socket", "", "specifies a socket to listen on, takes priority over -port")
		set.IntVar(&flags.SocketUid, "socket_uid", -1, "if desired, change the owning UID on the listening socket")
		set.IntVar(&flags.SocketGid, "socket_gid", -1, "if desired, change the owning GID on the listening socket")
		set.BoolVar(&flags.AllowInsecureLdap, "allow_insecure_ldap", false, "allow insecure connections to LDAP")
		set.StringVar(&flags.AuthSameSite, "auth_same_site", "", "configure the SameSite attribute of returned cookies")

		if err := set.Parse(os.Args[2:]); err != nil {
			l.Fatalf("Could not parse arguments: %s\n", err)
		}

		os.Exit(routine.RunServe(flags))

	case "dump_routes":
		os.Exit(routine.RunDumpRoutes())
	}
}
