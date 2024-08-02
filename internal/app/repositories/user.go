package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
				user_
			WHERE
				email_ = $1
		`
	)

	var (
		user repodto.UserRet

		err error
	)

	err = r.db.QueryRowContext(
		ctx,
		query,
		argDto.Email,
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil && err != sql.ErrNoRows {
		return user, fmt.Errorf("get user by email failed: %w", err)
	}

	return user, err
}

func (r *User) Create(ctx context.Context, argDto *repodto.CreateUserArg) (repodto.UserRet, error) {
	if argDto == nil {
		return repodto.UserRet{}, errors.New("nil reqDto")
	}

	const (
		query = `
			INSERT INTO user_ (
				fullname_, 
				email_, 
				password_
			) 
			VALUES (
				$1, 
				$2, 
				$3
			)
			RETURNING
				id_,
				fullname_,
				email_,
				password_,
				created_at_,
				updated_at_
		`
	)

	var (
		user repodto.UserRet

		err error
	)

	err = r.db.QueryRowContext(
		ctx,
		query,
		argDto.FullName,
		argDto.Email,
		argDto.Password,
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, fmt.Errorf("creating user failed: %w", err)
	}

	return user, nil
}
