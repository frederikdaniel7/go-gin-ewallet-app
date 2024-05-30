package entity

import (
	"database/sql"
	"time"
)

type PasswordToken struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
