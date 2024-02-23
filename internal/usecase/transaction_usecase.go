package usecase

import (
	"context"
	"log"
	"net/http"

	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/entity"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/repository"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/apperror"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/constant"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/database"
)

type TransactionUseCase interface {
	MakeTransfer(ctx context.Context, body *entity.Transfer, userID int64) (*entity.Transaction, error)
	TopUpBalance(ctx context.Context, body *entity.Transfer, userID int64) (*entity.Transaction, error)
}

type transactionUseCaseImpl struct {
	userRepository        repository.UserRepository
	walletRepository      repository.WalletRepository
	transactionRepository repository.TransactionRepository
	transactor            database.Transactor
}

func NewTransactionUseCaseImpl(
	userRepository repository.UserRepository,
	walletRespository repository.WalletRepository,
	transactionRepository repository.TransactionRepository,
	transactor database.Transactor,
) *transactionUseCaseImpl {
	return &transactionUseCaseImpl{
		userRepository:        userRepository,
		walletRepository:      walletRespository,
		transactionRepository: transactionRepository,
		transactor:            transactor,
	}
}
func (u *transactionUseCaseImpl) MakeTransfer(ctx context.Context, body *entity.Transfer, userID int64) (*entity.Transaction, error) {
	var user *entity.User
	var senderWallet *entity.Wallet
	var recipientWallet *entity.Wallet
	var transaction *entity.Transaction
	user, err := u.userRepository.FindUserById(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, err
	}
	err = u.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		senderWallet, err = u.walletRepository.FindWalletByUserID(txCtx, user.ID)
		if senderWallet == nil {
			return err
		}
		if err != nil {
			return err
		}
		if senderWallet.Balance.LessThan(body.Amount) {
			return apperror.NewInputErrorType(http.StatusBadRequest, constant.ResponseMsgInsufficientFunds)
		}
		recipientWallet, err = u.walletRepository.FindWalletByWalletNumber(txCtx, body.RecipientWalletNumber)
		if recipientWallet == nil {
			return err
		}
		if err != nil {
			return err
		}
		if senderWallet.ID == recipientWallet.ID {
			return apperror.NewInputErrorType(http.StatusBadRequest, constant.ResponseMsgCannotTransferToSelf)
		}
		_, err = u.walletRepository.UpdateDecreaseWalletBalance(txCtx, senderWallet, body.Amount)
		if err != nil {
			return err
		}
		_, err = u.walletRepository.UpdateAddWalletBalance(txCtx, recipientWallet, body.Amount)
		if err != nil {
			return err
		}
		transaction, err = u.transactionRepository.CreateTransaction(ctx, &entity.Transaction{
			SenderWalletID:    &senderWallet.ID,
			RecipientWalletID: recipientWallet.ID,
			Amount:            body.Amount,
			SourceOfFunds:     constant.SourceWallet,
			Descriptions:      body.Descriptions,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (u *transactionUseCaseImpl) TopUpBalance(ctx context.Context, body *entity.Transfer, userID int64) (*entity.Transaction, error) {
	var user *entity.User
	var recipientWallet *entity.Wallet
	var transaction *entity.Transaction
	log.Println(userID)
	user, err := u.userRepository.FindUserById(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, err
	}
	err = u.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		recipientWallet, err = u.walletRepository.FindWalletByUserID(txCtx, user.ID)
		if recipientWallet == nil {
			return err
		}
		if err != nil {
			return err
		}
		_, err = u.walletRepository.UpdateAddWalletBalance(txCtx, recipientWallet, body.Amount)
		if err != nil {
			return err
		}
		transaction, err = u.transactionRepository.CreateTransaction(ctx, &entity.Transaction{
			RecipientWalletID: recipientWallet.ID,
			Amount:            body.Amount,
			SourceOfFunds:     body.SourceOfFunds,
			Descriptions:      body.Descriptions,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
