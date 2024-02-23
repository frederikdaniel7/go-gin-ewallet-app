package dto

import (
	"time"
)

type PasswordToken struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

type CreatePasswordToken struct {
	Email string `json:"email" binding:"required,email"`
}

type PassTokenObj struct {
	Token PasswordToken `json:"token"`
}

type ResetPassword struct {
	Password string `json:"password" binding:"required,min=8"`
}

type PasswordTokenParam struct {
	Token string `uri:"token" binding:"required"`
}
