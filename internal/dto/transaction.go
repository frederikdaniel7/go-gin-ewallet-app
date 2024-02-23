package dto

import (
	"database/sql"
	"time"

	d "github.com/shopspring/decimal"
)

type Transaction struct {
	ID                int64        `json:"id"`
	SenderWalletID    int64        `json:"sender_wallet_id"`
	RecipientWalletID int64        `json:"recipient_wallet_id"`
	Amount            d.Decimal    `json:"amount"`
	SourceOfFunds     string       `json:"source_of_funds"`
	Descriptions      string       `json:"descriptions"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
	DeletedAt         sql.NullTime `json:"deleted_at"`
}

type Transfer struct {
	RecipientWalletNumber string  `json:"to" binding:"required,number,len=13,startswith=420"`
	Amount                float64 `json:"amount" binding:"required,min=1000,max=50000000"`
	Descriptions          string  `json:"description" binding:"required"`
}
