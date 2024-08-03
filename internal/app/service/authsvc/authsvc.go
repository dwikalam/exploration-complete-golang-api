package authsvc

import (
	"context"
	"errors"
	"fmt"

	"github.com/dwikalam/ecommerce-service/internal/app/helperdependency/crypto/icrypto"
	"github.com/dwikalam/ecommerce-service/internal/app/service/authsvc/authsvcdto"
	"github.com/dwikalam/ecommerce-service/internal/app/store/userstore/iuserstore"
	"github.com/dwikalam/ecommerce-service/internal/app/store/userstore/userstoredto"
	"github.com/dwikalam/ecommerce-service/internal/app/transaction/itransaction"
)

type Auth struct {
	txManager itransaction.Manager
	authStore iuserstore.Storer
	crypter   icrypto.Crypter
}

func NewAuth(
	txManager itransaction.Manager,
	authStore iuserstore.Storer,
	crypter icrypto.Crypter,
) (Auth, error) {
	if txManager == nil {
		return Auth{}, errors.New("nil txManager")
	}

	if authStore == nil {
		return Auth{}, errors.New("nil authStore")
	}

	if crypter == nil {
		return Auth{}, errors.New("nil crypter")
	}

	return Auth{
		txManager: txManager,
		authStore: authStore,
		crypter:   crypter,
	}, nil
}

func (s *Auth) RegisterUser(
	ctx context.Context,
	fullName string,
	email string,
	password string,
) (authsvcdto.RegisteredUser, error) {
	const (
		svcErrorPlaceholder string = "register user failed"
	)

	var (
		hashedPassword string

		createUserStoreDto userstoredto.User

		registerUserSvcDto authsvcdto.RegisteredUser

		err error
	)

	if _, err = s.authStore.GetByEmail(ctx, email); err == nil {
		return authsvcdto.RegisteredUser{}, fmt.Errorf("%s: email already exist", svcErrorPlaceholder)
	}

	if hashedPassword, err = s.crypter.Hash(password); err != nil {
		return authsvcdto.RegisteredUser{}, fmt.Errorf("%s: %w", svcErrorPlaceholder, err)
	}

	createUserStoreDto, err = s.authStore.Create(ctx, fullName, email, hashedPassword)
	if err != nil {
		return authsvcdto.RegisteredUser{}, fmt.Errorf("%s: %w", svcErrorPlaceholder, err)
	}

	registerUserSvcDto = authsvcdto.RegisteredUser{
		ID:       createUserStoreDto.ID,
		FullName: createUserStoreDto.FullName,
		Email:    createUserStoreDto.Email,
		Password: createUserStoreDto.Password,
	}

	return registerUserSvcDto, nil
}

func (s *Auth) ValidateLoginAttempt(
	ctx context.Context,
	email string,
	plainPassword string,
) error {
	const (
		svcErrorPlaceholder string = "validate login attempt failed"
	)

	var (
		userRepoDto userstoredto.User

		err error
	)

	if userRepoDto, err = s.authStore.GetByEmail(ctx, email); err != nil {
		return fmt.Errorf("%s: %w", svcErrorPlaceholder, err)
	}

	if err = s.crypter.Compare(userRepoDto.Password, plainPassword); err != nil {
		return fmt.Errorf("%s: %w", svcErrorPlaceholder, err)
	}

	return nil
}
