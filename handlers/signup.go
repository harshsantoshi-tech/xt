package handlers

import (
	"errors"
	"expense-tracker/configs"
	"expense-tracker/models"
	"expense-tracker/services"
	"expense-tracker/services/redis"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

// SendOTPHandler handles the initial signup request
func SendOTPHandler(email string, password string) error {
	// 1. Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// 2. Generate a 6-digit OTP
	// Seed is handled automatically in Go 1.20+, otherwise use rand.Seed
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 3. Store in Redis using the modular service
	pendingUser := models.PendingUser{
		Password: string(hashedPassword),
		OTP:      otp,
	}

	err = redis.SavePendingUser(email, pendingUser)
	if err != nil {
		return fmt.Errorf("failed to save temporary user: %v", err)
	}

	// 4. Send the OTP via Email
	return services.SendEmail(email, otp, "to complete your registration:")
}

// VerifySignupHandler checks the OTP and moves the user to MySQL
func VerifySignupHandler(email string, inputOTP string) error {
	// 1. Retrieve data from Redis via service
	pending, err := redis.GetPendingUser(email)
	if err != nil {
		return errors.New("OTP expired or signup not initiated")
	}

	// 2. Check OTP
	if pending.OTP != inputOTP {
		// Increment retry count via service
		count, _ := redis.IncrementOTPRetry(email)

		if count >= 3 {
			// Brute force detected: Wipe the data via service
			redis.ClearSignupData(email)
			return errors.New("too many failed attempts, please request a new OTP")
		}

		return fmt.Errorf("invalid OTP, %d attempts remaining", 3-count)
	}

	// 3. OTP is correct! Save to MySQL
	query := "INSERT INTO users (username , email, password_hash, created_at) VALUES (?,?, ?, ?)"
	_, err = configs.AppConfig.DB.Exec(query, "success", email, pending.Password, time.Now())

	if err != nil {
		return fmt.Errorf("failed to create user in database: %v", err)
	}

	// 4. Cleanup Redis via service
	redis.ClearSignupData(email)

	return nil
}
