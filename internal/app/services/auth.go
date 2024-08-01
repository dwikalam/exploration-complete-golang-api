package services

import (
	"context"
	"errors"

	"github.com/dwikalam/ecommerce-service/internal/app/helpers"
	"github.com/dwikalam/ecommerce-service/internal/app/repositories"
	"github.com/dwikalam/ecommerce-service/internal/app/types/dto/repodto"
	"github.com/dwikalam/ecommerce-service/internal/app/types/dto/svcdto"
	"github.com/dwikalam/ecommerce-service/internal/app/types/interfaces"
)

type Auth struct {
	txManager interfaces.TransactionManager
	authRepo  *repositories.User
}

func NewAuth(
	txManager interfaces.TransactionManager,
	authRepo *repositories.User,
) (Auth, error) {
	if txManager == nil {
		return Auth{}, errors.New("nil txManager")
	}

	if authRepo == nil {
		return Auth{}, errors.New("nil authRepo")
	}

	return Auth{
		txManager: txManager,
		authRepo:  authRepo,
	}, nil
}

func (s *Auth) RegisterUser(
	ctx context.Context,
	argDto *svcdto.RegisterUserArg,
) (svcdto.RegisteredUserRet, error) {
	var (
		getUserByEmailArgDto repodto.GetUserByEmailArg = repodto.GetUserByEmailArg{
			Email: argDto.Email,
		}

		hashedPassword []byte

		createUserArgDto repodto.CreateUserArg
		createUserRetDto repodto.UserRet

		svcRetDto svcdto.RegisteredUserRet

		err error
	)

	if _, err = s.authRepo.GetByEmail(ctx, getUserByEmailArgDto); err != nil {
		return svcdto.RegisteredUserRet{}, err
	}

	if hashedPassword, err = helpers.BcryptHashedPassword(argDto.Password); err != nil {
		return svcdto.RegisteredUserRet{}, err
	}

	createUserArgDto = repodto.CreateUserArg{
		FullName: argDto.FullName,
		Email:    argDto.Email,
		Password: string(hashedPassword),
	}
	createUserRetDto, err = s.authRepo.Create(ctx, &createUserArgDto)

	svcRetDto = svcdto.RegisteredUserRet{
		ID:       createUserRetDto.ID,
		FullName: createUserRetDto.FullName,
		Email:    createUserRetDto.Email,
		Password: createUserRetDto.Password,
	}

	return svcRetDto, err
}
