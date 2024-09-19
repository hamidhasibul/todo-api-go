package repository

import (
	"database/sql"
	"fmt"
	"n0ctRnull/todo-api-go/database"
	"n0ctRnull/todo-api-go/models"
)

func InsertPost(post *models.Post, userId string) error {
	query := "INSERT INTO posts(title,description,userId) VALUES($1,$2,$3)"

	_, err := database.Db.Exec(query, post.Title, post.Description, userId)
	return err
}

func FetchPosts(userId, page, limit int) ([]models.Post, int, error) {
	posts := []models.Post{}
	offset := (page - 1) * limit
	query := "SELECT id,title,description,userId FROM posts WHERE userId=$1 LIMIT $2 OFFSET $3"
	rows, err := database.Db.Query(query, userId, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		post := models.Post{}
		if err := rows.Scan(&post.Id, &post.Title, &post.Description, &post.UserId); err != nil {
			return nil, 0, err
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM posts WHERE userId=$1"
	err = database.Db.QueryRow(countQuery, userId).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func UpdatePost(post *models.Post, postId string) error {
	query := "UPDATE posts SET title=$1, description=$2 WHERE id=$3"

	_, err := database.Db.Exec(query, post.Title, post.Description, postId)
	return err
}

func DeletePost(postId string) error {
	query := "DELETE FROM posts WHERE id=$1"

	_, err := database.Db.Exec(query, postId)

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
