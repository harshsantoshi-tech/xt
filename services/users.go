package services

import (
	"database/sql"
	"errors"
	"expense-tracker/configs"
	"expense-tracker/models"
	sql2 "expense-tracker/services/sql"
	"fmt"
	"github.com/labstack/gommon/log"

	"github.com/golang-jwt/jwt/v5"
)

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := configs.AppConfig.DB.QueryRow(sql2.GET_USERS_FROM_EMAIL, email).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
	)

	if err != nil {
		log.Error("GetUserByEmail ", err.Error())
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

// GetUserIDFromToken parses a JWT string and returns the UserID
func GetUserIDFromToken(tokenString string, secretKey string) (int64, error) {

	if secretKey == "" {
		return 0, errors.New("internal server error: secret key not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		// Return the UserID from the payload
		return claims.UserID, nil
	}

	return 0, errors.New("invalid token: claims could not be verified")
}
