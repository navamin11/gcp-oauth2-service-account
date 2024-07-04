package routes

import (
	"gcp-access-token/handler"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRouter() *fiber.App {
	r := fiber.New()
	r.Use(logger.New())

	r.Get("/", handler.Healthcheck)
	r.Get("/command", handler.UseCommand)
	r.Get("/lib", handler.UseLib)
	r.Get("/manual", handler.NotUseLib)

	if err := r.Listen(":3000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	return r
}
