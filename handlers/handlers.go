package handlers

import (
	"fmt"
	"n0ctRnull/todo-api-go/helpers"
	"n0ctRnull/todo-api-go/models"
	"n0ctRnull/todo-api-go/repository"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RegisterUser(ctx *fiber.Ctx) error {

	newUser := models.User{}

	if err := ctx.BodyParser(&newUser); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	userExists, err := repository.FindUserByEmail(newUser.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Something went wrong"})

	}

	if userExists != nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"success": false, "message": "User already exists"})
	}

	hashedPwd, err := helpers.HashPassword(newUser.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Something went wrong"})
	}

	newUser.Password = hashedPwd

	err = repository.InsertUser(newUser)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Something went wrong"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "user has been added"})
}

func LoginUser(ctx *fiber.Ctx) error {

	loginReq := models.LoginRequest{}

	if err := ctx.BodyParser(&loginReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request",
		})
	}

	user, err := repository.FindUserByEmail(loginReq.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Something went wrong"})
	}

	if user == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "User not found"})
	}

	isValid := helpers.IsValidPassword(user.Password, loginReq.Password)

	if !isValid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid credentials"})
	}

	token, err := helpers.GenerateToken(user.Id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Something went wrong"})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		Secure:   true,
		HTTPOnly: true,
	}

	ctx.Cookie(&cookie)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "token": token})
}

func AddPost(ctx *fiber.Ctx) error {

	newPost := models.Post{}

	if err := ctx.BodyParser(&newPost); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"success": false, "message": "Bad request"})
	}

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"]
	userId := fmt.Sprintf("%.0f", id)

	err := repository.InsertPost(&newPost, userId)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Something went wrong"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "Post added"})
}

func UpdatePost(ctx *fiber.Ctx) error {

	updatedPost := models.Post{}
	if err := ctx.BodyParser(&updatedPost); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Bad request"})
	}

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"]
	userId, err := strconv.Atoi(fmt.Sprintf("%.0f", id))
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
		})
	}
	todoId := ctx.Params("todoId")

	post, err := repository.FindPostById(todoId)
	if err != nil {
		if err.Error() == "post not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "No such post found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
		})
	}

	if post.UserId != userId {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Forbidden",
		})
	}

	if err = repository.UpdatePost(&updatedPost, todoId); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": true, "message": "Post updated"})
}

func DeletePost(ctx *fiber.Ctx) error {

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"]
	userId, err := strconv.Atoi(fmt.Sprintf("%.0f", id))
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
		})
	}

	todoId := ctx.Params("todoId")

	post, err := repository.FindPostById(todoId)
	if err != nil {
		if err.Error() == "post not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "No such post found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
		})
	}

	if post.UserId != userId {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Forbidden",
		})
	}

	if err = repository.DeletePost(todoId); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
		})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func GetPosts(ctx *fiber.Ctx) error {

	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"]
	userId, err := strconv.Atoi(fmt.Sprintf("%.0f", id))
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
		})
	}

	posts, total, err := repository.FetchPosts(userId, page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    posts,
		"page":    page,
		"limit":   limit,
		"total":   total,
	})
}
