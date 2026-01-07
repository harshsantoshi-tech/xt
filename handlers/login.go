package handlers

import (
	"errors"
	"expense-tracker/configs"
	"expense-tracker/models"
	"expense-tracker/services"
	"expense-tracker/services/redis"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
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

// ForgotPasswordHandler sends OTP if user exists
func ForgotPasswordHandler(email string) error {

	_, err := services.GetUserByEmail(email)
	if err != nil {
		return errors.New("if this email is registered, an OTP has been sent")
	}

	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	pending := models.PendingUser{OTP: otp}
	err = redis.SavePendingUser(email, pending)
	if err != nil {
		return err
	}

	return services.SendEmail(email, otp, "to reset your password:")
}

// ResetPasswordHandler updates the password in DB
func ResetPasswordHandler(email string, newPassword string) error {

	resetAllowed, err := redis.GetRedis(fmt.Sprintf(redis.RESET_ALLOWED, email))
	if resetAllowed != "true" {
		return errors.New("unauthorized: please verify OTP first")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	query := "UPDATE users SET password_hash = ? WHERE email = ?"
	_, err = configs.AppConfig.DB.Exec(query, string(hashedPassword), email)

	if err == nil {
		redis.DeleteRedis(fmt.Sprintf(redis.RESET_ALLOWED, email))
	}

	return err
}
