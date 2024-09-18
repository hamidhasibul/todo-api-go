package repository

import (
	"database/sql"
	"fmt"
	"n0ctRnull/todo-api-go/database"
	"n0ctRnull/todo-api-go/models"
)

func InsertPost(post *models.Post, userId string) error {
	sql := "INSERT INTO posts(title,description,userId) VALUES($1,$2,$3)"

	_, err := database.Db.Exec(sql, post.Title, post.Description, userId)
	return err
}

func UpdatePost(post *models.Post, postId string) error {
	query := "UPDATE posts SET title=$1, description=$2 WHERE id=$3"

	_, err := database.Db.Exec(query, post.Title, post.Description, postId)
	return err
}

func FindPostById(postId string) (*models.Post, error) {
	post := models.Post{}
	query := "SELECT * from posts WHERE id=$1"
	err := database.Db.QueryRow(query, postId).Scan(&post.Id, &post.Description, &post.Title, &post.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, err
	}
	return &post, nil
}
