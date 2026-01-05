package models

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}