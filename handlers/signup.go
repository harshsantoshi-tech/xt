package handlers

import (
	"errors"
	"expense-tracker/constants"
	"expense-tracker/models"
	"expense-tracker/services"
	"expense-tracker/services/redis"
	"expense-tracker/services/sql"
	"fmt"
	redis2 "github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
)

// SendOTPHandler  to generate and send OTP
func SendOTPHandler(email string, password string) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	err := sql.UpsertUser(email, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("database operation failed: %v", err)
	}

	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	err = redis.SavePendingUser(email, models.PendingUser{
		Password: string(hashedPassword),
		OTP:      otp,
	})

	if err != nil {
		log.Printf("Redis failed for %s. Manual cleanup not strictly needed due to 'pending' logic.", email)
		return errors.New("failed to generate verification code, please try again")
	}

	go func() {
		err := services.SendEmail(email, otp, "to complete your registration:")
		if err != nil {
			fmt.Println("Async Email Error:", err)
		}
	}()

	return nil
}

// VerifySignupHandler checks the OTP and moves the user to MySQL
func VerifySignupHandler(email string, otp string) *models.AppError {
	// 1. Check Redis
	pending, err := redis.GetPendingUser(email)
	if err != nil {
		if errors.Is(err, redis2.Nil) {
			return models.Unauthorized("OTP expired, please try again")
		}
		return models.InternalError("Redis connection failed")
	}

	// 2. Check OTP
	if pending.OTP != otp {
		return models.Unauthorized("OTP verification failed")
	}

	// 3. Update Database
	err = sql.UpdateUserStatus(email, constants.ONLINE)
	if err != nil {
		return models.InternalError("Database update failed")
	}

	return nil
}
