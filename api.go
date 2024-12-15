package main

import (
	"dohunt/pkg/domains"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func startAPI(domainData *map[string]domains.Domain, port int) {
	// Pass the engine to the Views
	app := fiber.New()

	app.Static("/", "/static")

	app.Get("/api/domains", func(c *fiber.Ctx) error {

		key := c.Query("key")

		if key != os.Getenv("DOCHKEY") {
			return c.Status(fiber.StatusUnauthorized).SendString("401 Unauthorized")
		}

		jsonData, err := json.Marshal(*domainData)
		if err != nil {
			return err
		}

		c.Set("Content-Type", "application/json")
		return c.Send(jsonData)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
