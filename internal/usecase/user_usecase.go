package usecase

import (
	"context"
	"net/http"
	"runtime/debug"

	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/entity"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/repository"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/apperror"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/constant"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/database"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/utils"
)

type UserUseCase interface {
	Login(ctx context.Context, body *entity.User) (int, error)
	RegisterUser(ctx context.Context, body *entity.User) (*entity.UserDetail, error)
	GenerateToken(ctx context.Context, body *entity.User) (*entity.PasswordToken, error)
	ResetPassword(ctx context.Context, body *entity.PasswordToken, newPassword string) error
}

type userUseCaseImpl struct {
	userRepository          repository.UserRepository
	walletRepository        repository.WalletRepository
	passwordTokenRepository repository.PasswordTokenRepository
	transactor              database.Transactor
}

func NewUserUseCaseImpl(
	userRepository repository.UserRepository,
	walletRespository repository.WalletRepository,
	passwordTokenRepository repository.PasswordTokenRepository,
	transactor database.Transactor,
) *userUseCaseImpl {
	return &userUseCaseImpl{
		userRepository:          userRepository,
		walletRepository:        walletRespository,
		passwordTokenRepository: passwordTokenRepository,
		transactor:              transactor,
	}
}
func (u *userUseCaseImpl) RegisterUser(ctx context.Context, body *entity.User) (*entity.UserDetail, error) {
	var user *entity.User
	var wallet *entity.Wallet

	err := u.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		checkUserExist, err := u.userRepository.FindUserByEmail(txCtx, body.Email)
		if err != nil {
			return err
		}
		if checkUserExist.Email != "" && checkUserExist.Email == body.Email {
			return apperror.NewUserErrorType(http.StatusBadRequest, constant.ResponseMsgUserAlreadyExists, debug.Stack())
		}

		user, err = u.userRepository.CreateUser(txCtx, body)
		if err != nil {
			return err
		}

		wallet, err = u.walletRepository.CreateWallet(txCtx, user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &entity.UserDetail{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Wallet:    wallet,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}, nil
}

func (u *userUseCaseImpl) Login(ctx context.Context, body *entity.User) (int, error) {
	user, err := u.userRepository.FindUserByEmail(ctx, body.Email)
	if user.Email == "" {
		return 0, apperror.NewInputErrorType(http.StatusBadRequest, constant.ResponseMsgUserDoesNotExist, debug.Stack())
	}
	if err != nil {
		return 0, err
	}
	password, err := u.userRepository.FindUserPassword(ctx, body)
	if err != nil {
		return 0, err
	}
	plainPassword, err := utils.CheckPassword(body.Password, []byte(password))
	if err != nil {
		return 0, apperror.NewCredentialsErrorType(http.StatusUnauthorized, constant.ResponseMsgErrorCredentials)
	}
	if !plainPassword {
		return 0, apperror.NewCredentialsErrorType(http.StatusUnauthorized, constant.ResponseMsgErrorCredentials)
	}
	return int(user.ID), err
}

func (u *userUseCaseImpl) GenerateToken(ctx context.Context, body *entity.User) (*entity.PasswordToken, error) {
	var token *entity.PasswordToken

	err := u.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		checkUserExist, err := u.userRepository.FindUserByEmail(txCtx, body.Email)
		if err != nil {
			return err
		}
		if checkUserExist.Email == "" {
			return apperror.NewUserErrorType(http.StatusBadRequest, constant.ResponseMsgUserDoesNotExist, debug.Stack())
		}
		oldToken, err := u.passwordTokenRepository.CheckToken(txCtx, checkUserExist)
		if err != nil {
			return err
		}
		if !oldToken.DeletedAt.Valid {
			_, err = u.passwordTokenRepository.UpdateDeleteToken(txCtx, checkUserExist, oldToken.Token)
			if err != nil {
				return err
			}
		}
		token, err = u.passwordTokenRepository.CreateToken(txCtx, checkUserExist)
		if err != nil {
			return err
		}
		return nil

	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (u *userUseCaseImpl) ResetPassword(ctx context.Context, body *entity.PasswordToken, newPassword string) error {

	err := u.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		token, err := u.passwordTokenRepository.GetValidToken(txCtx, body.Token)
		if token.UserID == 0 {
			return apperror.NewCredentialsErrorType(http.StatusBadRequest, constant.ResponseMsgInvalidToken)
		}
		if err != nil {
			return err
		}
		user, err := u.userRepository.FindUserById(txCtx, token.UserID)
		if err != nil {
			return err
		}
		_, err = u.userRepository.UpdateUserPassword(txCtx, user, newPassword)
		if err != nil {
			return err
		}
		_, err = u.passwordTokenRepository.UpdateDeleteToken(ctx, user, body.Token)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
