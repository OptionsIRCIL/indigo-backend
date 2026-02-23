//go:debug x509negativeserial=1
package main

import (
	"flag"
	"log"
	"os"
	"slices"

	"github.com/go-playground/validator/v10"
	"myoptions.info/indigo/backend/internal/routine"
	"myoptions.info/indigo/backend/internal/util"
)

func main() {
	// Use a logger with no prefix for program startup
	l := log.New(os.Stderr, "", 0)
	l.Println("Indigo CIL, v0.0.0")

	// Check if valid subcommand exists
	if len(os.Args) < 2 || !slices.Contains([]string{"serve", "generate_openapi_spec", "create_user"}, os.Args[1]) {
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
  generate_openapi_spec:
        Generates an OpenAPI 3.1 spec in JSON format and dumps it to STDOUT.
  create_user:
        Create a new user for the local handler.`,
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

	case "generate_openapi_spec":
		os.Exit(routine.RunGenerateOpenApiSpec())

	case "create_user":
		flags := util.CreateUserRuntimeFlags{}
		set := flag.NewFlagSet("create_user", flag.ExitOnError)
		set.StringVar(&flags.Username, "username", "", "desired username")
		set.StringVar(&flags.Password, "password", "", "desired password")
		set.StringVar(&flags.FirstName, "first", "", "first name")
		set.StringVar(&flags.LastName, "last", "", "last name")

		if err := set.Parse(os.Args[2:]); err != nil {
			l.Fatalf("Could not parse arguments: %s\n", err)
		}

		v := validator.New()
		if err := v.Struct(flags); err != nil {
			l.Fatalf("Argument validation failed, see required syntax in -help\n%s", err)
		}

		os.Exit(routine.RunCreateUser(flags))
	}
}
