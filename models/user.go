package models

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type User struct {
	Id           int64     `json:"id"`
	Username     *string   `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	CreatedAt    time.Time `json:"created_at"`
}

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// PendingUser is the structure we Marshal into JSON to store in Redis
type PendingUser struct {
	Password string `json:"password"`
	OTP      string `json:"otp"`
}
