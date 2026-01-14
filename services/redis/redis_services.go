package redis

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

	key := fmt.Sprintf(PENDING_USER, email)
	return configs.AppConfig.Redis.Set(ctx, key, data, 15*time.Minute).Err()
}

// GetPendingUser retrieves the temporary signup data
func GetPendingUser(email string) (*models.PendingUser, error) {
	key := fmt.Sprintf(PENDING_USER, email)
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
	key := fmt.Sprintf(OTP_RETRIES, email)

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
	configs.AppConfig.Redis.Del(ctx, fmt.Sprintf(PENDING_USER, email), fmt.Sprintf(OTP_RETRIES, email))
}

// IsResetAuthorized checks if the user has successfully verified their OTP recently
func IsResetAuthorized(email string) bool {
	key := fmt.Sprintf(RESET_ALLOWED, email)
	val, _ := configs.AppConfig.Redis.Get(ctx, key).Result()
	return val == "true"
}

func SetRedis(key string, value string, expiration time.Duration) error {
	return configs.AppConfig.Redis.Set(ctx, key, value, expiration).Err()
}

func GetRedis(key string) (string, error) {
	return configs.AppConfig.Redis.Get(ctx, key).Result()
}

func DeleteRedis(key string) error {
	return configs.AppConfig.Redis.Del(ctx, key).Err()
}
