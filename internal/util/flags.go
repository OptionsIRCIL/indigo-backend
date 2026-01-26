package util

type ServeRuntimeFlags struct {
	Port              int
	Socket            string
	SocketUid         int
	SocketGid         int
	AllowInsecureLdap bool
	AuthSameSite      string
}
