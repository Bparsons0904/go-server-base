package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/bparsons094/go-server-base/database"
	"github.com/bparsons094/go-server-base/models"
	"github.com/bparsons094/go-server-base/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AuthenticateUser(c *fiber.Ctx) error {
	var token string

	authorizationHeader := c.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) != 0 && fields[0] == "Bearer" {
		token = fields[1]
	}

	if token == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "You are not logged in",
		})
	}

	sub, err := utils.ValidateToken(token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	if cachedUser, found := utils.GetUser(sub); found {
		c.Locals("currentUser", cachedUser)
		return c.Next()
	}

	var user models.User
	err = database.DB.Where("id = ?", sub).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "User not found",
			})
		} else {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Token Error",
			})
		}
	}

	c.Locals("currentUser", user)
	return c.Next()
}
