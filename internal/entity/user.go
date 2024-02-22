package entity

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64
	Email     string
	Name      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type UserDetail struct {
	ID        int64
	Email     string
	Name      string
	Wallet    *Wallet
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
