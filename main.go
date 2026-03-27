package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/jiehui555/meaw-oa/internal/config"
	"github.com/jiehui555/meaw-oa/internal/database"
)

func main() {
	cfg := config.Load()

	db := database.Init(cfg.DBPath)
	_ = db

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
