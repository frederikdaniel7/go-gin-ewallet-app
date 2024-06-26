package entity

import (
	"database/sql"
	"time"

	d "github.com/shopspring/decimal"
)

type Transaction struct {
	ID                int64
	SenderWalletID    *int64
	SenderName		  *string
	RecipientWalletID int64
	RecipientName     string
	Amount            d.Decimal
	SourceOfFunds     string
	Descriptions      string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         sql.NullTime
}

type Transfer struct {
	RecipientWalletNumber string
	Amount                d.Decimal
	SourceOfFunds         string
	Descriptions          string
}

type TransactionPage struct {
	Transactions []Transaction
	ItemCount    int
	PageCount    int
	CurrentPage  int
}

type TransactionFilter struct {
	Search    string
	SortBy    string
	Order     string
	Transactiontype string
	Page      *int
	Limit     *int
	StartDate string
	EndDate   string
}
