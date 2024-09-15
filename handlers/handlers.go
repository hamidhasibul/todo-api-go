package handlers

import (
	"n0ctRnull/todo-api-go/helpers"
	"n0ctRnull/todo-api-go/models"
	"n0ctRnull/todo-api-go/repository"
	"time"

	"github.com/gofiber/fiber/v2"
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
