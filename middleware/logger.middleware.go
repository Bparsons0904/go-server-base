package middleware

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bparsons094/go-server-base/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func LogMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestTime := time.Now()

		// Check if the Content-Type is for multipart form data
		contentType := string(c.Request().Header.ContentType())
		isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

		// Only read and log the body if it's not a multipart request
		var bodyBytes []byte
		if !isMultipart {
			bodyBytes = c.Body()
		}

		// Proceed with the actual request
		c.Next()

		// Capture the response set by the handler
		responseBytes := c.Response().Body()

		responseTime := time.Now()
		duration := responseTime.Sub(requestTime).Seconds()

		// Safely attempt to get the current user
		// If no user is found, set userID to nil
		var userID *uuid.UUID
		contextUserID := c.Locals("UserID")

		if contextUserID != nil {
			contextUserID := contextUserID.(uuid.UUID)
			userID = &contextUserID
		}

		logEntry := models.RequestLog{
			RequestTime:  requestTime,
			ResponseTime: responseTime,
			UserID:       userID,
			Duration:     duration,
			Method:       c.Method(),
			Path:         c.Path(),
			Headers:      fmt.Sprintf("%v", string(c.Request().Header.Header())),
			Body:         string(bodyBytes),
			Response:     string(responseBytes),
		}

		// Write the log to the database in a separate goroutine
		go func(db *gorm.DB, logEntry models.RequestLog) {
			if err := db.Create(&logEntry).Error; err != nil {
				log.Println("Error creating log entry:", err, logEntry.Path)
			}
		}(db, logEntry)

		return c.Send(responseBytes)
	}
}
