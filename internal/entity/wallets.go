package entity

import (
	"database/sql"
	"time"

	d "github.com/shopspring/decimal"
)

type Wallet struct {
	ID           int64
	UserID       int64
	WalletNumber string
	Balance      d.Decimal
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}
