package util

type ServeRuntimeFlags struct {
	Port              int
	Socket            string
	SocketUid         int
	SocketGid         int
	AllowInsecureLdap bool
	AuthSameSite      string
}

type CreateUserRuntimeFlags struct {
	Username  string `validate:"required"`
	Password  string `validate:"required,min=8"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
}
