package repository

import (
	"n0ctRnull/todo-api-go/database"
	"n0ctRnull/todo-api-go/models"
)

func InsertPost(post *models.Post, userId string) error {
	sql := "INSERT INTO posts(title,description,userId) VALUES($1,$2,$3)"

	_, err := database.Db.Exec(sql, post.Title, post.Description, userId)
	return err
}
