package services

import (
	"errors"
	"expense-tracker/configs"
	"expense-tracker/models"
)

var DB = configs.AppConfig.DB

func GetUserFromEmail(email string) (models.User, error) {
	var user models.User

	// err := DB.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).
	// 	Scan(&user.Id, &user.Email, &user.Password)

	// if err != nil {
	// 	log.Error("GetUserFromEmail.QueryRow.ERROR ", email)
	// 	return user, err
	// }
	return user, nil

}

func CreateUser(user *models.User) error {

	if user == nil {
		return errors.New("user is nil")
	}

	// query := "INSERT INTO users (email, password,username) VALUES (?, ?,?)"

	// // Execute the insert query
	// result, err := DB.Exec(query, user.Email, user.Password, user.Username)
	// if err != nil {
	// 	log.Error("CreateUser.Exec.ERROR ", err)
	// 	return err
	// }

	// // Get the inserted ID and assign it back to the user
	// id, err := result.LastInsertId()
	// if err != nil {
	// 	log.Error("CreateUser.LastInsertId.ERROR ", err)
	// 	return err
	// }
	// user.Id = id

	return nil
}
