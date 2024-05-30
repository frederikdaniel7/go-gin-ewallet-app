package usecase

import (
	"context"
	"fmt"
	"log"

	"math"
	"net/http"
	"runtime/debug"

	"github.com/frederikdaniel7/go-gin-ewallet-app/internal/entity"
	"github.com/frederikdaniel7/go-gin-ewallet-app/internal/repository"
	"github.com/frederikdaniel7/go-gin-ewallet-app/pkg/apperror"
	"github.com/frederikdaniel7/go-gin-ewallet-app/pkg/constant"
	"github.com/frederikdaniel7/go-gin-ewallet-app/pkg/database"
)

type TransactionUseCase interface {
	MakeTransfer(ctx context.Context, body *entity.Transfer, userID int64) (*entity.Transaction, error)
	TopUpBalance(ctx context.Context, body *entity.Transfer, userID int64) (*entity.Transaction, error)
	GetTransaction(ctx context.Context, params entity.TransactionFilter, userID int64) (*entity.TransactionPage, error)
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
			return apperror.NewInputErrorType(http.StatusBadRequest, constant.ResponseMsgInsufficientFunds, debug.Stack())
		}
		_, err = u.walletRepository.UpdateDecreaseWalletBalance(txCtx, senderWallet, body.Amount)
		if err != nil {
			return err
		}
		recipientWallet, err = u.walletRepository.FindWalletByWalletNumber(txCtx, body.RecipientWalletNumber)
		if recipientWallet == nil {
			return err
		}
		if err != nil {
			return err
		}
		if senderWallet.ID == recipientWallet.ID {
			return apperror.NewInputErrorType(http.StatusBadRequest, constant.ResponseMsgCannotTransferToSelf, debug.Stack())
		}
		_, err = u.walletRepository.UpdateAddWalletBalance(txCtx, recipientWallet, body.Amount)
		if err != nil {
			return err
		}
		transaction, err = u.transactionRepository.CreateTransaction(txCtx, &entity.Transaction{
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
		transaction, err = u.transactionRepository.CreateTransaction(txCtx, &entity.Transaction{
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

func (u *transactionUseCaseImpl) GetTransaction(ctx context.Context, params entity.TransactionFilter, userID int64) (*entity.TransactionPage, error) {
	var transactions []entity.Transaction
	var countData int
	pageCount := 1.0
	var defaultPage = 1
	countData, err := u.transactionRepository.CountAllTransactions(ctx, userID, params)
	if err != nil {
		return nil, err
	}

	if params.Limit != nil {
		pageCount = float64(countData) / float64(*params.Limit)
	}

	pageCountInt := 1
	if math.RoundToEven(pageCount) < pageCount {
		pageCountInt = int(math.RoundToEven(pageCount)) + 1
		log.Printf(fmt.Sprintf("page count: %v", pageCountInt))
	} else if math.RoundToEven(pageCount) > pageCount {
		pageCountInt = int(math.RoundToEven(pageCount))
	} else {
		pageCountInt = int((pageCount))
	}

	if params.Page == nil || *params.Page > pageCountInt {
		params.Page = &defaultPage
	}
	log.Printf(fmt.Sprintf("page count: %v", pageCountInt))
	transactions, err = u.transactionRepository.GetAllTransactions(ctx, userID, entity.TransactionFilter{
		Search:          params.Search,
		SortBy:          params.SortBy,
		Order:           params.Order,
		Transactiontype: params.Transactiontype,
		Page:            params.Page,
		Limit:           params.Limit,
		StartDate:       params.StartDate,
		EndDate:         params.EndDate,
	})
	if err != nil {
		return nil, err
	}

	return &entity.TransactionPage{
		Transactions: transactions,
		ItemCount:    countData,
		PageCount:    pageCountInt,
		CurrentPage:  *params.Page,
	}, nil
}
