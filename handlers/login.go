package handlers

import (
	"database/sql"
	"errors"
	"expense-tracker/models"
	"expense-tracker/services"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func LoginHandler(email, password string) (models.User, error) {
	var user models.User

	// Fetch user by email
	user, err := services.GetUserFromEmail(email)
	if err != nil {
		// Use errors.Is for proper error comparison
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}
		return user, fmt.Errorf("failed to fetch user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, ErrInvalidCredentials
	}

	fmt.Printf("User %s (ID %d) logged in successfully\n", user.Email, user.Id)
	return user, nil
}
