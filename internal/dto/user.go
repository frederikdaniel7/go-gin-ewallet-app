package dto

import (
	"time"
)

type User struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type UserDetail struct {
	ID        int64          `json:"id"`
	Email     string         `json:"email"`
	Name      string         `json:"name"`
	Wallet    *WalletPreview `json:"wallet"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *time.Time     `json:"deleted_at"`
}

type UserObj struct {
	User UserDetail `json:"user"`
}

type CreateUser struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}
