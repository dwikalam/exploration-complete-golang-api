package repodto

type GetUserByEmailArg struct {
	Email string
}

type CreateUserArg struct {
	FullName string
	Email    string
	Password string
}

type UserRet struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
