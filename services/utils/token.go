package utils

import (
	"expense-tracker/constants"
	"expense-tracker/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/gommon/log"
	"os"
	"time"
)

func GenerateToken(userID int64 ,expiry time.Duration) (string, error) {
	// 1. Get the secret from environment variables
	secretKey := os.Getenv("JWT_SECRET")

	// 2. Create the Claims (the data inside the token)
	claims := &models.CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
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

func GenerateTokens(userID int64 ) (string, string, error) {
	accessToken, err := GenerateToken(userID ,constants.JWT_EXPIRY_TIME)
	if err != nil {
		log.Error("GenerateTokens.GenerateToken "," access token ", userID)
		return "", "", err
	}

	refreshToken , err := GenerateToken(userID,constants.REFRESH_TOKEN_EXPIRY_TIME )
	if err != nil {
		log.Error("GenerateTokens.GenerateToken "," refresh token ", userID)
		return "", "", err
	}
	return accessToken, refreshToken, err
}