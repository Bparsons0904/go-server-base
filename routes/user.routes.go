package routes

import (
	"github.com/bparsons094/go-server-base/controllers"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(api fiber.Router) {
	userRoutes := api.Group("/users")
	userRoutes.Get("/getMe", controllers.GetMe)
}
