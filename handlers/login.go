package handlers

import (
	"errors"
	"expense-tracker/configs"
	"expense-tracker/models"
	"expense-tracker/services"
	"expense-tracker/services/redis"
	"expense-tracker/services/sql"
	"expense-tracker/services/utils"
	"fmt"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

func LoginHandler(email, password string) (string, error) {
	// 1. Fetch user from DB
	user, err := services.GetUserByEmail(email)

	if err != nil {
		log.Error("Error getting user by email: "+err.Error(), " ", email)
		return "", fmt.Errorf("user not found")
	}

	// 2. Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("incorrect password")
	}

	// 3. Generate Token
	token, err := utils.GenerateToken(user.Id)
	if err != nil {
		log.Error("Error generating token ", user.Id, " ", err)
		return "", fmt.Errorf("failed to generate token")
	}

	return token, nil
}

// ForgotPasswordHandler sends OTP if user exists
func ForgotPasswordHandler(email string) error {

	_, err := services.GetUserByEmail(email)
	if err != nil {
		log.Error("ForgotPasswordHandler.GetUserByEmail ", email, " ", err)
		return errors.New("if this email is registered, an OTP has been sent")
	}

	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	pending := models.PendingUser{OTP: otp}
	err = redis.SavePendingUser(email, pending)
	if err != nil {
		log.Error("Error saving pending user ", err, " ", email)
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

	return services.SendEmail(email, newOtp,"")
}