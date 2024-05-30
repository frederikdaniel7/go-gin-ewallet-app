package repository

import (
	"context"
	"database/sql"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/frederikdaniel7/go-gin-ewallet-app/internal/entity"
	"github.com/frederikdaniel7/go-gin-ewallet-app/pkg/apperror"
	"github.com/frederikdaniel7/go-gin-ewallet-app/pkg/constant"
	"github.com/frederikdaniel7/go-gin-ewallet-app/pkg/database"
	"github.com/frederikdaniel7/go-gin-ewallet-app/pkg/utils"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, body *entity.Transaction) (*entity.Transaction, error)
	GetAllTransactions(ctx context.Context, userID int64, filter entity.TransactionFilter) ([]entity.Transaction, error)
	CountAllTransactions(ctx context.Context, userID int64, filter entity.TransactionFilter) (int, error)
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
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal, debug.Stack())
	}
	return &transaction, nil

}

func (r *transactionRepository) GetAllTransactions(ctx context.Context, userID int64, filter entity.TransactionFilter) ([]entity.Transaction, error) {
	var sb strings.Builder
	transactions := []entity.Transaction{}
	runner := database.PickQuerier(ctx, r.db)
	var data []interface{}
	sb.WriteString(`SELECT t.id, ws.wallet_number, us.name as sender_name, wr.wallet_number, ur.name as recipient_name,t.amount, t.source_of_funds, t.descriptions,t.created_at,t.updated_at,t.deleted_at from transactions t 
	JOIN wallets w ON w.id = t.sender_wallet_id OR w.id = t.recipient_wallet_id
	LEFT JOIN wallets ws ON t.sender_wallet_id = ws.id 
	LEFT JOIN wallets wr ON t.recipient_wallet_id = wr.id
	LEFT JOIN users us ON us.id = t.sender_wallet_id  
	LEFT JOIN users ur ON ur.id = t.recipient_wallet_id 
	where w.user_id = $1`)
	data = append(data, userID)
	queryString, dataParams := utils.ConvertQueryParamstoSql(filter)
	sb.WriteString(queryString)
	data = append(data, dataParams...)
	rows, err := runner.QueryContext(ctx, sb.String(), data...)
	if err != nil {
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal, debug.Stack())
	}
	defer rows.Close()
	for rows.Next() {
		transaction := entity.Transaction{}
		err := rows.Scan(&transaction.ID, &transaction.SenderWalletID, &transaction.SenderName, &transaction.RecipientWalletID, &transaction.RecipientName, &transaction.Amount, &transaction.SourceOfFunds, &transaction.Descriptions, &transaction.CreatedAt, &transaction.UpdatedAt, &transaction.DeletedAt)
		if err != nil {
			return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal, debug.Stack())
		}
		transactions = append(transactions, transaction)
		// log.Printf("transaction %#v",transaction.RecipientName )
	}

	err = rows.Err()
	if err != nil {
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal, debug.Stack())
	}

	return transactions, nil
}

func (r *transactionRepository) CountAllTransactions(ctx context.Context, userID int64, filter entity.TransactionFilter) (int, error) {
	var sb strings.Builder
	var dataCount int
	runner := database.PickQuerier(ctx, r.db)
	var data []interface{}
	sb.WriteString(`SELECT COUNT (*) from transactions t JOIN wallets w ON w.id = t.sender_wallet_id OR w.id = t.recipient_wallet_id where w.user_id = $1`)
	data = append(data, userID)
	queryString, dataParams := utils.ConvertQueryParamstoSqlforCount(filter)
	sb.WriteString(queryString)
	data = append(data, dataParams...)
	err := runner.QueryRowContext(ctx, sb.String(), data...).Scan(&dataCount)
	if err != nil {
		return 0, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal, debug.Stack())
	}

	return dataCount, nil
}
