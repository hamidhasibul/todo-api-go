package main

import (
	"fmt"
	"n0ctRnull/todo-api-go/database"
	"n0ctRnull/todo-api-go/handlers"
	"n0ctRnull/todo-api-go/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
)

func testHandler(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"]
	fmt.Println(id)
	return ctx.SendString("Hello World")

}

func main() {
	app := fiber.New()
	database.ConnectDatabase()
	app.Use(logger.New())

	app.Get("/", middlewares.VerifyJWT(), testHandler)
	app.Post("/register", handlers.RegisterUser)
	app.Post("/login", handlers.LoginUser)

	app.Post("/posts", middlewares.VerifyJWT(), handlers.AddPost)

	app.Listen(":3000")
}
