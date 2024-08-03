package userstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dwikalam/ecommerce-service/internal/app/db/sqldb/isqldb"
	"github.com/dwikalam/ecommerce-service/internal/app/store/userstore/userstoredto"
)

type SQLStore struct {
	sqldb isqldb.Querier
}

func NewSQLStore(sqldb isqldb.Querier) (SQLStore, error) {
	if sqldb == nil {
		return SQLStore{}, errors.New("nil sqldb")
	}

	return SQLStore{
		sqldb: sqldb,
	}, nil
}

func (store *SQLStore) GetByEmail(ctx context.Context, email string) (userstoredto.User, error) {
	const (
		query string = `
			SELECT 
				*
			FROM
				user_
			WHERE
				email_ = ?
		`
	)

	var (
		user userstoredto.User

		err error
	)

	err = store.sqldb.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil && err != sql.ErrNoRows {
		return user, fmt.Errorf("QueryRowContext failed: %v", err)
	}

	return user, err
}

func (store *SQLStore) Create(
	ctx context.Context,
	fullName string,
	email string,
	password string,
) (userstoredto.User, error) {
	const (
		query = `
			INSERT INTO user_ (
				fullname_, 
				email_, 
				password_
			) 
			VALUES (
				?, 
				?, 
				?
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
		user userstoredto.User

		err error
	)

	err = store.sqldb.QueryRowContext(
		ctx,
		query,
		fullName,
		email,
		password,
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, fmt.Errorf("QueryRowContext failed: %v", err)
	}

	return user, nil
}
