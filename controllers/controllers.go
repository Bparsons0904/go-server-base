package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetDb(db *gorm.DB) {
	DB = db
}

func getUserId(c *fiber.Ctx) uuid.UUID {
	userID := c.Locals("UserID").(uuid.UUID)
	return userID
}

func getSubAccountId(c *fiber.Ctx) (uuid.UUID, error) {
	type RequestData struct {
		SubAccountId string `json:"subAccountId"`
	}

	var data RequestData
	if err := c.BodyParser(&data); err != nil {
		return uuid.Nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Failed parsing request body"})
	}

	subAccountId, err := uuid.Parse(data.SubAccountId)
	if err != nil {
		return uuid.Nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Failed parsing sub account id"})
	}

	return subAccountId, nil
}

func stringToUuid(str string) (uuid.UUID, error) {
	id, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
