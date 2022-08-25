package main

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/zyzmoz/go-image-recognition/vision"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		var r = vision.CloudVision()

		results, err := json.Marshal(r)
		if err != nil {
			return c.Status(500).SendString("Error parsing data")
		}

		return c.Send(results)
	})

	app.Listen(":3000")
}
