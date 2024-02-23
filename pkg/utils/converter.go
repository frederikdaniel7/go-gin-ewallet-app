package utils

import (
	"fmt"
	"strings"

	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/dto"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/entity"
)

func ConvertUserDetailtoJson(user entity.UserDetail) dto.UserDetail {
	var walletNumSB strings.Builder
	walletNumSB.WriteString("420")
	walletNumSB.WriteString(user.Wallet.WalletNumber)
	converted := dto.UserDetail{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Wallet: &dto.WalletPreview{
			WalletNumber: walletNumSB.String(),
			Balance:      user.Wallet.Balance,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if user.DeletedAt.Valid {
		converted.DeletedAt = &user.DeletedAt.Time
	}
	return converted
}

func ConvertTokentoJson(token entity.PasswordToken) dto.PasswordToken {
	converted := dto.PasswordToken{
		Token:     token.Token,
		ExpiredAt: token.ExpiredAt,
	}
	return converted
}
func ConvertWalletNumber(walletNum string) string {
	return walletNum[3:13]
}

func ConvertTransactiontoJson(transaction entity.Transaction) dto.Transaction {
	senderWalletNumber := 4200000000000 + *transaction.SenderWalletID
	recipientWalletNumber := 4200000000000 + transaction.RecipientWalletID
	converted := dto.Transaction{
		ID:                transaction.ID,
		SenderWalletID:    &senderWalletNumber,
		RecipientWalletID: recipientWalletNumber,
		Amount:            transaction.Amount,
		SourceOfFunds:     transaction.SourceOfFunds,
		Descriptions:      transaction.Descriptions,
		CreatedAt:         transaction.CreatedAt,
		UpdatedAt:         transaction.UpdatedAt,
	}
	return converted
}

func ConvertTransactionstoJson(transactions []entity.Transaction) []dto.Transaction {
	transactionsJson := []dto.Transaction{}

	for _, values := range transactions {
		senderWalletNumber := 4200000000000 + *values.SenderWalletID
		RecipientWalletNumber := 4200000000000 + values.RecipientWalletID
		transactionsJson = append(transactionsJson, dto.Transaction{
			ID:                values.ID,
			SenderWalletID:    &senderWalletNumber,
			RecipientWalletID: RecipientWalletNumber,
			Amount:            values.Amount,
			SourceOfFunds:     values.SourceOfFunds,
			Descriptions:      values.Descriptions,
			CreatedAt:         values.CreatedAt,
			UpdatedAt:         values.UpdatedAt,
		})
	}
	return transactionsJson
}

func ConvertQueryJsonToObject(queryParams dto.TransactionFilter) entity.TransactionFilter {
	converted := entity.TransactionFilter{
		Search:    queryParams.Search,
		SortBy:    queryParams.SortBy,
		Order:     queryParams.Order,
		Page:      queryParams.Page,
		Limit:     queryParams.Limit,
		StartDate: queryParams.StartDate,
		EndDate:   queryParams.EndDate,
	}
	return converted
}

func ConvertQueryParamstoSql(params entity.TransactionFilter) (string, []interface{}) {
	var query strings.Builder

	var filters []interface{}
	var countParams = 2

	if params.Search != "" {
		query.WriteString(fmt.Sprintf(` AND t.descriptions ILIKE '%%' ||$%v|| '%%'`, countParams))
		filters = append(filters, params.Search)
		countParams++
	}
	if params.StartDate != "" {
		query.WriteString(fmt.Sprintf(` AND t.created_at >= $%v`, countParams))
		filters = append(filters, params.StartDate)
		countParams++
	}
	if params.EndDate != "" {
		query.WriteString(fmt.Sprintf(` AND t.created_at <= $%v`, countParams))
		filters = append(filters, params.EndDate)
		countParams++
	}
	if params.SortBy != "" {
		switch params.SortBy {
		case "to":
			query.WriteString(` ORDER BY t.recipient_wallet_id`)
		case "amount":
			query.WriteString(` ORDER BY t.amount`)
		case "date":
			query.WriteString(` ORDER BY t.created_at`)
		}
	} else if params.SortBy == "" {
		query.WriteString(` ORDER BY t.created_at`)
	}
	if params.Order != "" {
		if params.Order == "asc" {
			query.WriteString(` ASC`)
		}
		if params.Order == "desc" {
			query.WriteString(` DESC`)
		}
	}
	if params.Limit != nil {
		query.WriteString(fmt.Sprintf(` LIMIT $%v`, countParams))
		filters = append(filters, *params.Limit)
		countParams++
	}
	if params.Page != nil {
		setLimit := *params.Limit
		if params.Limit == nil {
			setLimit = 10
		}
		query.WriteString(fmt.Sprintf(` OFFSET $%v`, countParams))
		filters = append(filters, (*params.Page-1)*setLimit)
		countParams++
	}
	return query.String(), filters

}

func ConvertQueryParamstoSqlforCount(params entity.TransactionFilter) (string, []interface{}) {
	var query strings.Builder

	var filters []interface{}
	var countParams = 2

	if params.Search != "" {
		query.WriteString(fmt.Sprintf(` AND t.descriptions ILIKE '%%' ||$%v|| '%%'`, countParams))
		filters = append(filters, params.Search)
		countParams++
	}
	if params.StartDate != "" {
		query.WriteString(fmt.Sprintf(` AND t.created_at >= $%v`, countParams))
		filters = append(filters, params.StartDate)
		countParams++
	}
	if params.EndDate != "" {
		query.WriteString(fmt.Sprintf(` AND t.created_at <= $%v`, countParams))
		filters = append(filters, params.EndDate)
		countParams++
	}

	return query.String(), filters

}
