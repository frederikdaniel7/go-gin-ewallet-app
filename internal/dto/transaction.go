package dto

import (
	"time"

	d "github.com/shopspring/decimal"
)

type Transaction struct {
	ID                int64     `json:"id"`
	SenderWalletID    *int64    `json:"sender_wallet_id"`
	RecipientWalletID int64     `json:"recipient_wallet_id"`
	Amount            d.Decimal `json:"amount"`
	SourceOfFunds     string    `json:"source_of_funds"`
	Descriptions      string    `json:"descriptions"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Transactions struct {
	Transactions []Transaction `json:"transactions"`
}

type Transfer struct {
	RecipientWalletNumber string  `json:"to" binding:"required,number,len=13,startswith=420"`
	Amount                float64 `json:"amount" binding:"required,min=1000,max=50000000"`
	Descriptions          string  `json:"description" binding:"required"`
}

type TopUpBody struct {
	Amount        float64 `json:"amount" binding:"required,min=50000,max=50000000"`
	SourceOfFunds int     `json:"source_of_funds" binding:"required,min=1,max=4"`
}

type TransactionFilter struct {
	Search    string `form:"s"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof='amount' 'date' 'to'"`
	Order     string `form:"order" binding:"omitempty,oneof='asc' 'desc'"`
	Page      *int   `form:"page" binding:"omitempty,min=1"`
	Limit     *int   `form:"limit" binding:"omitempty,min=1"`
	StartDate string `form:"start" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"end" binding:"omitempty,datetime=2006-01-02"`
}
