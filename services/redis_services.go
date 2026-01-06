package services

import (
	"context"
	"encoding/json"
	"expense-tracker/configs"
	"expense-tracker/models"
	"fmt"
	"time"
)

var ctx = context.Background()

// SavePendingUser stores the hashed password and OTP with a 5-minute TTL
func SavePendingUser(email string, pending models.PendingUser) error {
	data, err := json.Marshal(pending)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("pending_user:%s", email)
	return configs.AppConfig.Redis.Set(ctx, key, data, 5*time.Minute).Err()
}

// GetPendingUser retrieves the temporary signup data
func GetPendingUser(email string) (*models.PendingUser, error) {
	key := fmt.Sprintf("pending_user:%s", email)
	val, err := configs.AppConfig.Redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var pending models.PendingUser
	if err := json.Unmarshal([]byte(val), &pending); err != nil {
		return nil, err
	}

	return &pending, nil
}

// IncrementOTPRetry tracks failed attempts and returns the current count
func IncrementOTPRetry(email string) (int64, error) {
	key := fmt.Sprintf("otp_retries:%s", email)

	// Increment the counter
	count, err := configs.AppConfig.Redis.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// Set expiration on first failure
	if count == 1 {
		configs.AppConfig.Redis.Expire(ctx, key, 5*time.Minute)
	}

	return count, nil
}

// ClearSignupData removes all OTP related data from Redis after success or lockout
func ClearSignupData(email string) {
	configs.AppConfig.Redis.Del(ctx,
		fmt.Sprintf("pending_user:%s", email),
		fmt.Sprintf("otp_retries:%s", email),
	)
}
