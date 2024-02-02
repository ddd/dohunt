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

		// TODO: Switch to DOCHKEY env variable
		if key != os.Getenv("DOCHKEY") {
			return c.Status(fiber.StatusUnauthorized).SendString("401 Unauthorized")
		}

		// Convert the map to JSON
		jsonData, err := json.Marshal(*domainData)
		if err != nil {
			return err
		}

		// Set the response content type as JSON
		c.Set("Content-Type", "application/json")

		// Return the JSON data
		return c.Send(jsonData)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
