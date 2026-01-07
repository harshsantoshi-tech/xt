package handlers

import (
	"errors"
	"expense-tracker/configs"
	"expense-tracker/models"
	"expense-tracker/services"
	"expense-tracker/services/redis"
	"expense-tracker/services/sql"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

// SendOTPHandler  to generate and send OTP
func SendOTPHandler(email string, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	pendingUser := models.PendingUser{
		Password: string(hashedPassword),
		OTP:      otp,
	}

	err = redis.SavePendingUser(email, pendingUser)
	if err != nil {
		log.Printf("failed to save pending user: %v", err)
		return fmt.Errorf("failed to save temporary user: %v", err)
	}

	return services.SendEmail(email, otp, "to complete your registration:")
}

// VerifySignupHandler checks the OTP and moves the user to MySQL
func VerifySignupHandler(email string, inputOTP string) error {

	pending, err := redis.GetPendingUser(email)
	if err != nil {
		return errors.New("OTP expired or signup not initiated")
	}

	// 2. Check OTP
	if pending.OTP != inputOTP {
		count, _ := redis.IncrementOTPRetry(email)

		if count >= 3 {
			redis.ClearSignupData(email)
			return errors.New("too many failed attempts, please request a new OTP")
		}

		return fmt.Errorf("invalid OTP, %d attempts remaining", 3-count)
	}

	_, err = configs.AppConfig.DB.Exec(sql.INSERT_USERS, "user", email, pending.Password, time.Now())

	if err != nil {
		return fmt.Errorf("failed to create user in database: %v", err)
	}

	redis.ClearSignupData(email)

	return nil
}
