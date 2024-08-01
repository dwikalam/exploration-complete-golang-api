package svcdto

type RegisterUserArg struct {
	FullName string
	Email    string
	Password string
}

type RegisteredUserRet struct {
	ID       string
	FullName string
	Email    string
	Password string
}
