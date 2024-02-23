package entity

import (
	"database/sql"
	"time"

	d "github.com/shopspring/decimal"
)

type Transaction struct {
	ID                int64
	SenderWalletID    int64
	RecipientWalletID int64
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
	Descriptions          string
}

type TransactionFilter struct {
}
