package utils

import (
	"expense-tracker/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/gommon/log"
	"os"
	"time"
)

func GenerateToken(userID int64) (string, error) {
	// 1. Get the secret from environment variables
	secretKey := os.Getenv("JWT_SECRET")

	// 2. Create the Claims (the data inside the token)
	claims := &models.CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "xt-chat-app",
		},
	}

	// 3. Create the token object with the HS256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 4. Sign the token with your secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	log.Info("Generate token successfully ", tokenString)

	return tokenString, nil
}

