package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func AllowHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		_, err := c.WriteString("gooood")
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return err
		}
		c.Status(http.StatusOK)
		return nil
	}
}
