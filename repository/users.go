package repository

import (
	"database/sql"
	"n0ctRnull/todo-api-go/database"
	"n0ctRnull/todo-api-go/models"
)

func FindUserByEmail(email string) (*models.User, error) {
	user := models.User{}
	query := "SELECT id, name, email,password FROM users WHERE email=$1"
	err := database.Db.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func InsertUser(user models.User) error {
	query := "INSERT INTO users(name,email,password) VALUES ($1,$2,$3)"
	_, err := database.Db.Exec(query, user.Name, user.Email, user.Password)
	return err

}
