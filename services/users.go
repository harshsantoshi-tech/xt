package services

import (
	"database/sql"
	"errors"
	"expense-tracker/configs"
	"expense-tracker/models"
	"fmt"
	"github.com/labstack/gommon/log"

	"github.com/golang-jwt/jwt/v5"
	
)

var DB = configs.AppConfig.DB

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	query := `SELECT id, username, email, password_hash FROM users WHERE email = ? LIMIT 1`

	err := DB.QueryRow(query, email).Scan(
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
    // 1. Safety check for empty secret
    if secretKey == "" {
        return 0, errors.New("internal server error: secret key not set")
    }

    // 2. Parse the token with CustomClaims (assuming models.CustomClaims)
    token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Validate the algorithm is HMAC
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secretKey), nil
    })

    if err != nil {
        return 0, err
    }

    // 3. Extract the claims - TYPE MUST MATCH Step 1 exactly
    if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
        // Return the UserID from the payload
        return claims.UserID, nil
    }

    return 0, errors.New("invalid token: claims could not be verified")
}