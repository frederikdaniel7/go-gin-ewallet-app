package repository

import (
	"context"
	"database/sql"
	"net/http"
	"runtime/debug"

	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/entity"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/apperror"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/constant"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/database"
	d "github.com/shopspring/decimal"
)

type WalletRepository interface {
	FindWalletByWalletNumber(ctx context.Context, walletNumber string) (*entity.Wallet, error)
	FindWalletByUserID(ctx context.Context, userID int64) (*entity.Wallet, error)
	CreateWallet(ctx context.Context, body *entity.User) (*entity.Wallet, error)
	UpdateAddWalletBalance(ctx context.Context, w *entity.Wallet, amount d.Decimal) (*entity.Wallet, error)
	UpdateDecreaseWalletBalance(ctx context.Context, w *entity.Wallet, amount d.Decimal) (*entity.Wallet, error)
}

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *walletRepository {
	return &walletRepository{
		db: db,
	}
}
func (r *walletRepository) CreateWallet(ctx context.Context, body *entity.User) (*entity.Wallet, error) {
	wallet := entity.Wallet{}
	runner := database.PickQuerier(ctx, r.db)

	q := `INSERT INTO wallets (user_id) 
	VALUES ($1)
	returning id, user_id, wallet_number, balance, created_at, updated_at, deleted_at`

	err := runner.QueryRowContext(ctx, q, body.ID).Scan(&wallet.ID, &wallet.UserID, &wallet.WalletNumber, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
	if err != nil {
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal, debug.Stack())
	}
	return &wallet, nil
}

func (r *walletRepository) FindWalletByUserID(ctx context.Context, userID int64) (*entity.Wallet, error) {
	wallet := entity.Wallet{}
	runner := database.PickQuerier(ctx, r.db)

	q := `SELECT id, user_id, wallet_number, balance, created_at, updated_at, deleted_at from wallets where user_id = $1 FOR UPDATE`
	err := runner.QueryRowContext(ctx, q, userID).Scan(&wallet.ID, &wallet.UserID, &wallet.WalletNumber, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgUserDoesNotExist, debug.Stack())
		}
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal, debug.Stack())
	}
	return &wallet, nil
}

func (r *walletRepository) FindWalletByWalletNumber(ctx context.Context, walletNumber string) (*entity.Wallet, error) {
	wallet := entity.Wallet{}
	runner := database.PickQuerier(ctx, r.db)

	q := `SELECT id, user_id, wallet_number, balance, created_at, updated_at, deleted_at from wallets where wallet_number = $1 FOR UPDATE`

	err := runner.QueryRowContext(ctx, q, walletNumber).Scan(&wallet.ID, &wallet.UserID, &wallet.WalletNumber, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgUserDoesNotExist, debug.Stack())
		}
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal, debug.Stack())
	}
	return &wallet, nil
}

func (r *walletRepository) UpdateAddWalletBalance(ctx context.Context, w *entity.Wallet, amount d.Decimal) (*entity.Wallet, error) {
	wallet := entity.Wallet{}
	runner := database.PickQuerier(ctx, r.db)

	q := `UPDATE wallets set balance = balance + $1, updated_at = now() where id = $2 RETURNING id, user_id, wallet_number, balance, created_at, updated_at, deleted_at`

	err := runner.QueryRowContext(ctx, q, amount, w.ID).
		Scan(&wallet.ID, &wallet.UserID, &wallet.WalletNumber, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
	if err != nil {
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal, debug.Stack())
	}
	return &wallet, nil
}

func (r *walletRepository) UpdateDecreaseWalletBalance(ctx context.Context, w *entity.Wallet, amount d.Decimal) (*entity.Wallet, error) {
	wallet := entity.Wallet{}
	runner := database.PickQuerier(ctx, r.db)

	q := `UPDATE wallets set balance = balance - $1, updated_at = now() where id = $2 RETURNING id, user_id, wallet_number, balance, created_at, updated_at, deleted_at`

	err := runner.QueryRowContext(ctx, q, amount, w.ID).
		Scan(&wallet.ID, &wallet.UserID, &wallet.WalletNumber, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt)
	if err != nil {
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal, debug.Stack())
	}
	return &wallet, nil
}
