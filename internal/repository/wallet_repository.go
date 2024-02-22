package repository

import (
	"context"
	"database/sql"
	"net/http"

	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/entity"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/apperror"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/constant"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/database"
)

type WalletRepository interface {
	// FindAll(ctx context.Context) ([]entity.User, error)
	// FindSimilarUserByName(ctx context.Context, name string) ([]entity.User, error)
	// FindWalletByWalletNumber(ctx context.Context, walletNumber string) (*entity.User, error)
	// FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateWallet(ctx context.Context, body *entity.User) (*entity.Wallet, error)
	// FindUserPassword(ctx context.Context, body entity.User) (string, error)
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
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal)
	}
	return &wallet, nil
}
