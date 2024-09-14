package main

import (
	"n0ctRnull/todo-api-go/database"
	"n0ctRnull/todo-api-go/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func testHandler(ctx *fiber.Ctx) error {
	return ctx.SendString("Hello World")
}

func main() {
	app := fiber.New()
	database.ConnectDatabase()
	app.Use(logger.New())

	app.Get("/", testHandler)
	app.Post("/register", handlers.RegisterUser)

	app.Listen(":3000")
}
