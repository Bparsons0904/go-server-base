package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func GetMe(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": "user"})
}
