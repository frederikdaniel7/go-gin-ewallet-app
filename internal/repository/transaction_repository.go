package repository

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/entity"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/apperror"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/constant"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/database"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/utils"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, body *entity.Transaction) (*entity.Transaction, error)
	GetAllTransactions(ctx context.Context, userID int64, filter entity.TransactionFilter) ([]entity.Transaction, error)
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

func (r *transactionRepository) GetAllTransactions(ctx context.Context, userID int64, filter entity.TransactionFilter) ([]entity.Transaction, error) {
	//page := 1
	var sb strings.Builder
	transactions := []entity.Transaction{}
	runner := database.PickQuerier(ctx, r.db)
	var data []interface{}
	sb.WriteString(`SELECT t.id, t.sender_wallet_id, t.recipient_wallet_id,t.amount, t.source_of_funds, t.descriptions,t.created_at,t.updated_at,t.deleted_at from transactions t JOIN wallets w ON w.id = t.sender_wallet_id OR w.id = t.recipient_wallet_id where w.user_id = $1`)
	data = append(data, userID)
	queryString, dataParams := utils.ConvertQueryParamstoSql(filter)
	sb.WriteString(queryString)
	data = append(data, dataParams...)
	rows, err := runner.QueryContext(ctx, sb.String(), data...)
	if err != nil {
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal)
	}
	defer rows.Close()
	for rows.Next() {
		transaction := entity.Transaction{}
		err := rows.Scan(&transaction.ID, &transaction.SenderWalletID, &transaction.RecipientWalletID, &transaction.Amount, &transaction.SourceOfFunds, &transaction.Descriptions, &transaction.CreatedAt, &transaction.UpdatedAt, &transaction.DeletedAt)
		if err != nil {
			return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal)
		}
		transactions = append(transactions, transaction)
	}
	err = rows.Err()
	if err != nil {
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal)
	}

	return transactions, nil
}
