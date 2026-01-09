package sql

import (
	"database/sql"
	"errors"
	"expense-tracker/configs"
	"time"
)

func UpsertUser(email , hashedPassword string)error{
	var status string
	err := configs.AppConfig.DB.QueryRow("SELECT status FROM users WHERE email = ?", email).Scan(&status)

	if err == nil {
		// User exists - check their status
		if status != "pending" {
			return errors.New("user already registered with this email")
		}
		// If they are pending, update their password (in case they changed it on retry)
		_, err = configs.AppConfig.DB.Exec(
			"UPDATE users SET password_hash = ?, created_at = ? WHERE email = ? AND status = 'pending'",
			string(hashedPassword), time.Now(), email,
		)
	} else if errors.Is(err, sql.ErrNoRows) {
		// New user - perform initial insert
		_, err = configs.AppConfig.DB.Exec(
			"INSERT INTO users (email, password_hash, status, created_at) VALUES (?, ?, 'pending', ?)",
			email, string(hashedPassword), time.Now(),
		)
	}

	return err
}


func UpdateUserStatus(email string, status string)error{

	query := "UPDATE users SET status = ? WHERE email = ? AND status = 'pending'"
	result, err := configs.AppConfig.DB.Exec(query,status, email)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		err = errors.New("registration failed: user not found or already verified")
	}
	return err
}