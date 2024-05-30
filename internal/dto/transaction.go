package dto

import (
	"time"

	d "github.com/shopspring/decimal"
)

type Transaction struct {
	ID                int64     `json:"id"`
	SenderWalletID    *int64    `json:"sender_wallet_number"`
	SenderName        *string   `json:"sender_name"`
	RecipientWalletID int64     `json:"recipient_wallet_number"`
	RecipientName     string    `json:"recipient_name"`
	Amount            d.Decimal `json:"amount"`
	SourceOfFunds     string    `json:"source_of_funds"`
	Descriptions      string    `json:"descriptions"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Transactions struct {
	Transactions []Transaction `json:"transactions"`
}

type TransactionObj struct {
	Transaction Transaction `json:"transaction"`
}

type Transfer struct {
	RecipientWalletNumber string  `json:"to" binding:"required,number,len=13,startswith=420"`
	Amount                float64 `json:"amount" binding:"required,min=1000,max=50000000"`
	Descriptions          string  `json:"description" binding:"required,max=35"`
}

type TopUpBody struct {
	Amount        float64 `json:"amount" binding:"required,min=50000,max=50000000"`
	SourceOfFunds int     `json:"source_of_funds" binding:"required,min=1,max=4"`
}

type TransactionPage struct {
	Transactions Transactions `json:"transactions"`
	ItemCount    int          `json:"item_count"`
	PageCount    int          `json:"page_count"`
	CurrentPage  int          `json:"current_page"`
}

type TransactionFilter struct {
	Search          string `form:"s"`
	SortBy          string `form:"sort_by" binding:"omitempty,oneof='amount' 'date' 'to'"`
	Order           string `form:"order" binding:"omitempty,oneof='asc' 'desc'"`
	Transactiontype string `form:"txtype" binding:"omitempty,oneof='transfer' 'topup' 'all'"`
	Page            *int   `form:"page" binding:"omitempty,min=1"`
	Limit           *int   `form:"limit" binding:"omitempty,min=1"`
	StartDate       string `form:"start" binding:"omitempty,datetime=2006-01-02"`
	EndDate         string `form:"end" binding:"omitempty,datetime=2006-01-02"`
}
