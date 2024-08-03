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
	userStore iuserstore.Storer
	crypter   icrypto.Crypter
}

func New(
	txManager itransaction.Manager,
	userStore iuserstore.Storer,
	crypter icrypto.Crypter,
) (Auth, error) {
	if txManager == nil {
		return Auth{}, errors.New("nil txManager")
	}

	if userStore == nil {
		return Auth{}, errors.New("nil userStore")
	}

	if crypter == nil {
		return Auth{}, errors.New("nil crypter")
	}

	return Auth{
		txManager: txManager,
		userStore: userStore,
		crypter:   crypter,
	}, nil
}

func (s *Auth) RegisterUser(
	ctx context.Context,
	fullName string,
	email string,
	password string,
) (authsvcdto.RegisteredUser, error) {
	var (
		isEmailExist        bool
		hashedPassword      string
		createdUserStoreDto userstoredto.User
		registerUserSvcDto  authsvcdto.RegisteredUser

		isEmailExistErr error
		hashErr         error
		createUserErr   error
	)

	isEmailExist, isEmailExistErr = s.userStore.IsEmailExist(ctx, email)

	// Time consuming operation. Security measure
	hashedPassword, hashErr = s.crypter.Hash(password)

	if isEmailExistErr != nil {
		return authsvcdto.RegisteredUser{}, fmt.Errorf("IsEmailExist failed: %v", isEmailExistErr)
	}

	if hashErr != nil {
		return authsvcdto.RegisteredUser{}, fmt.Errorf("hash password failed: %v", hashErr)
	}

	if isEmailExist {
		return authsvcdto.RegisteredUser{}, errors.New("email already exist")
	}

	createdUserStoreDto, createUserErr = s.userStore.Create(ctx, fullName, email, hashedPassword)
	if createUserErr != nil {
		return authsvcdto.RegisteredUser{}, fmt.Errorf("creating user failed: %v", createUserErr)
	}

	registerUserSvcDto = authsvcdto.RegisteredUser{
		ID:       createdUserStoreDto.ID,
		FullName: createdUserStoreDto.FullName,
		Email:    createdUserStoreDto.Email,
		Password: createdUserStoreDto.Password,
	}

	return registerUserSvcDto, nil
}

func (s *Auth) ValidateLoginAttempt(
	ctx context.Context,
	email string,
	plainPassword string,
) error {
	const (
		placeholderHash string = "$2a$10$WIXnEK.SmlrME91uuiybY.aHwsxzwdM3FMVGHij4ztYoyL8pX5iBu"
	)

	var (
		userStoreDto   userstoredto.User
		hashedPassword string

		getByEmailErr error
		compareErr    error
	)

	userStoreDto, getByEmailErr = s.userStore.GetByEmail(ctx, email)

	switch userStoreDto.Password {
	case "":
		hashedPassword = placeholderHash
	default:
		hashedPassword = userStoreDto.Password
	}

	// Time consuming operation. Security measure
	compareErr = s.crypter.Compare(hashedPassword, plainPassword)

	if getByEmailErr != nil {
		return fmt.Errorf("get user by email failed: %v", getByEmailErr)
	}

	if compareErr != nil {
		return fmt.Errorf("compare passwords failed: %v", compareErr)
	}

	return nil
}
