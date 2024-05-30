package dto

import (
	"database/sql"
	"time"

	d "github.com/shopspring/decimal"
)

type Wallet struct {
	ID           int64        `json:"id"`
	UserID       int64        `json:"user_id"`
	WalletNumber string       `json:"wallet_number"`
	Balance      d.Decimal    `json:"balance"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    sql.NullTime `json:"deleted_at"`
}

type WalletPreview struct {
	WalletNumber string    `json:"wallet_number"`
	Balance      d.Decimal `json:"balance"`
}
