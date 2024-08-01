package repositories

import (
	"context"
	"errors"

	"github.com/dwikalam/ecommerce-service/internal/app/types/dto/repodto"
	"github.com/dwikalam/ecommerce-service/internal/app/types/interfaces"
)

type User struct {
	db interfaces.DbQuerier
}

func NewUser(db interfaces.DbQuerier) (User, error) {
	if db == nil {
		return User{}, errors.New("nil db")
	}

	return User{
		db: db,
	}, nil
}

func (r *User) GetByEmail(ctx context.Context, argDto repodto.GetUserByEmailArg) (repodto.UserRet, error) {
	const (
		query string = `
			SELECT 
				*
			FROM
				user
			WHERE
				email = ?
		`
	)

	var (
		user repodto.UserRet

		err error
	)

	err = r.db.QueryRowContext(
		ctx,
		query,
		[]any{
			argDto.Email,
		},
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Password,
	)

	return user, err
}

func (r *User) Create(ctx context.Context, argDto *repodto.CreateUserArg) (repodto.UserRet, error) {
	if argDto == nil {
		return repodto.UserRet{}, errors.New("nil reqDto")
	}

	const (
		query = `
			INSERT INTO user (
				fullname, 
				email, 
				password
			) 
			VALUES (
				?, 
				?, 
				?
			)
			RETURNING (
				id,
				fullname,
				email,
				password
			)
		`
	)

	var (
		user repodto.UserRet

		err error
	)

	err = r.db.QueryRowContext(
		ctx,
		query,
		[]any{
			argDto.FullName,
			argDto.Email,
			argDto.Password,
		},
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Password,
	)

	return user, err
}
