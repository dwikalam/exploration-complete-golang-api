package authsvc

import (
	"context"
	"errors"
	"fmt"
	"sync"

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
		wg sync.WaitGroup

		createUserStoreDto userstoredto.User
		registerUserSvcDto authsvcdto.RegisteredUser
		err                error

		isEmailExistChan  chan bool  = make(chan bool, 1)
		emailCheckErrChan chan error = make(chan error, 1)

		hashedPasswordChan chan string = make(chan string, 1)
		hashErrChan        chan error  = make(chan error, 1)
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		var (
			hashedPassword string
			err            error
		)

		hashedPassword, err = s.crypter.Hash(password)

		hashedPasswordChan <- hashedPassword
		hashErrChan <- err
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		var (
			isEmailExist bool
			err          error
		)

		isEmailExist, err = s.userStore.IsEmailExist(ctx, email)

		isEmailExistChan <- isEmailExist
		emailCheckErrChan <- err
	}()

	wg.Wait()

	if err = <-emailCheckErrChan; err != nil {
		return authsvcdto.RegisteredUser{}, fmt.Errorf("IsEmailExist failed: %v", err)
	}

	if err = <-hashErrChan; err != nil {
		return authsvcdto.RegisteredUser{}, fmt.Errorf("hash password failed: %v", err)
	}

	if <-isEmailExistChan {
		return authsvcdto.RegisteredUser{}, errors.New("email already exist")
	}

	createUserStoreDto, err = s.userStore.Create(ctx, fullName, email, <-hashedPasswordChan)
	if err != nil {
		return authsvcdto.RegisteredUser{}, fmt.Errorf("creating user failed: %v", err)
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
	var (
		wg sync.WaitGroup

		err error

		userStoreDtoChan  chan userstoredto.User = make(chan userstoredto.User, 1)
		getByEmailErrChan chan error             = make(chan error, 1)

		compareErrChan chan error = make(chan error, 1)
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		const (
			placeholderHash string = "$2a$10$WIXnEK.SmlrME91uuiybY.aHwsxzwdM3FMVGHij4ztYoyL8pX5iBu"
		)

		var (
			userStoreDto userstoredto.User
			err          error
		)

		userStoreDto, err = s.userStore.GetByEmail(ctx, email)

		getByEmailErrChan <- err

		switch err != nil {
		case true:
			userStoreDtoChan <- userstoredto.User{Password: placeholderHash}
		case false:
			userStoreDtoChan <- userStoreDto
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		compareErrChan <- s.crypter.Compare((<-userStoreDtoChan).Password, plainPassword)
	}()

	wg.Wait()

	if err = <-getByEmailErrChan; err != nil {
		return fmt.Errorf("get user by email failed: %v", err)
	}

	if err = <-compareErrChan; err != nil {
		return fmt.Errorf("compare passwords failed: %v", err)
	}

	return nil
}
