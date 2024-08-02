package repodto

import "time"

type GetUserByEmailArg struct {
	Email string
}

type CreateUserArg struct {
	FullName string
	Email    string
	Password string
}

type UserRet struct {
	ID        string
	FullName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
