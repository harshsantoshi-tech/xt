package handlers

import (
	"errors"
	"expense-tracker/configs"
	"expense-tracker/constants"
	"expense-tracker/models"
	"expense-tracker/services"
	"expense-tracker/services/redis"
	"expense-tracker/services/sql"
	"expense-tracker/services/utils"
	"fmt"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

func LoginHandler(email, password string) (string, string, error) {
	//Fetch user from DB
	user, err := services.GetUserByEmail(email)

	if err != nil {
		log.Error("Error getting user by email: "+err.Error(), " ", email)
		return "", "", fmt.Errorf("user not found")
	}

	//Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", fmt.Errorf("incorrect password")
	}

	//Generate Token
	accessToken, refreshToken, err := utils.GenerateTokens(user.Id)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate tokens")
	}

	expiresAt := time.Now().Add(constants.REFRESH_TOKEN_EXPIRY_TIME)
	err = sql.SaveUserSession(user.Id, refreshToken, expiresAt)
	if err != nil {
		log.Error("Session save error: ", err)
		return "", "", fmt.Errorf("failed to create session")
	}

	return accessToken, refreshToken, nil
}

// ForgotPasswordHandler sends OTP if user exists
func ForgotPasswordHandler(email string) error {

	_, err := services.GetUserByEmail(email)
	if err != nil {
		log.Error("ForgotPasswordHandler.GetUserByEmail ", email, " ", err)
		return errors.New("if this email is registered, an OTP has been sent")
	}

	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	otp = "123456"

	pending := models.PendingUser{OTP: otp}
	err = redis.SavePendingUser(email, pending)
	if err != nil {
		log.Error("Error saving pending user ", err, " ", email)
		return err
	}

	go func() {
		err := services.SendEmail(email, otp, "to reset your password:")
		if err != nil {
			fmt.Println("Async Email Error:", err)
		}
	}()

	return nil
}

// ResetPasswordHandler updates the password in DB
func ResetPasswordHandler(email string, newPassword string) error {

	resetAllowed, err := redis.GetRedis(fmt.Sprintf(redis.RESET_ALLOWED, email))
	if resetAllowed != "true" {
		return errors.New("unauthorized: please verify OTP first")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	_, err = configs.AppConfig.DB.Exec(sql.UPDATE_PASSWORD_BY_EMAIL, string(hashedPassword), email)

	if err == nil {
		redis.DeleteRedis(fmt.Sprintf(redis.RESET_ALLOWED, email))
	}

	return err
}

func ResendOTPHandler(email string) error {

	pending, err := redis.GetPendingUser(email)
	if err != nil {
		log.Error("ResendOTPHandler.GetPendingUser session expired or not found ", err, " ", email)
		return fmt.Errorf("session expired or not found, please start again")
	}

	newOtp := fmt.Sprintf("%06d", rand.Intn(1000000))
	pending.OTP = newOtp

	err = redis.SavePendingUser(email, *pending)
	if err != nil {
		log.Error("ResendOTPHandler.SavePendingUser Error saving pending user ", err, " ", email)
		return err
	}

	go func() {
		err := services.SendEmail(email, newOtp, "to complete your registration:")
		if err != nil {
			fmt.Println("Async Email Error:", err)
		}
	}()

	return nil
}