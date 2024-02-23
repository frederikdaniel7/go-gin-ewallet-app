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

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, body *entity.Transaction) (*entity.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *transactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, body *entity.Transaction) (*entity.Transaction, error) {
	transaction := entity.Transaction{}
	runner := database.PickQuerier(ctx, r.db)

	q := `INSERT INTO transactions (sender_wallet_id, recipient_wallet_id, amount, source_of_funds, descriptions)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, sender_wallet_id, recipient_wallet_id, amount, source_of_funds, descriptions, created_at
	`

	err := runner.QueryRowContext(ctx, q, body.SenderWalletID, body.RecipientWalletID, body.Amount, body.SourceOfFunds, body.Descriptions).
		Scan(&transaction.ID, &transaction.SenderWalletID, &transaction.RecipientWalletID, &transaction.Amount, &transaction.SourceOfFunds, &transaction.Descriptions, &transaction.CreatedAt)
	if err != nil {
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal)
	}
	return &transaction, nil

}
