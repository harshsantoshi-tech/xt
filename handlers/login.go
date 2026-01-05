package handlers

import (
	"errors"
	"expense-tracker/services"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func LoginHandler(email, password string) (string, error) {
	// 1. Fetch user from DB
	user, err := services.GetUserByEmail(email)

	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	// 2. Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("incorrect password")
	}

	// 3. Generate Token
	token, err := services.GenerateToken(user.Id)
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}

	return token, nil
}