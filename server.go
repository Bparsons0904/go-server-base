package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bparsons094/go-server-base/controllers"
	"github.com/bparsons094/go-server-base/database"
	"github.com/bparsons094/go-server-base/routes"
	"github.com/bparsons094/go-server-base/scheduler"
	"github.com/bparsons094/go-server-base/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	server *fiber.App
	config utils.Config
)

func init() {
	log.SetOutput(os.Stdout)

	config = utils.LoadConfig("./")

	db := database.ConnectDB(config)
	controllers.SetDb(db)

	if config.Environment == "local" {
		server = fiber.New(fiber.Config{
			ReadBufferSize:    16384,
			StreamRequestBody: true,
		})
	} else {
		server = fiber.New(fiber.Config{
			DisableStartupMessage: true,
			StreamRequestBody:     true,
			ReadBufferSize:        16384,
		})
	}

	go scheduler.InitScheduler(db)
}

func LoadEnvMiddleware(c *fiber.Ctx) error {
	c.Locals("clientOrigin", config.ClientOrigin)
	return c.Next()
}

func main() {
	server.Use(cors.New(cors.Config{
		AllowOrigins:     config.ClientOrigin,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
	server.Use(LoadEnvMiddleware)
	if config.Environment == "local" {
		server.Use(logger.New())
	}
	server.Use(recover.New())
	server.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", config.ClientOrigin)
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Credentials", "true")

		// Handle Preflight Request
		if c.Method() == "OPTIONS" {
			c.Status(fiber.StatusNoContent)
			return nil
		}

		return c.Next()
	})

	routes.SetupRoutes(server, config)

	// Creates a channel to listen for a shutdown signal
	go setupGracefulShutdown(server)

	log.Println("Server is running in", config.Environment, "mode on port", config.Port)
	log.Fatal(server.Listen(":" + config.Port))
}

func setupGracefulShutdown(server *fiber.App) {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-channel
		log.Println("Received termination signal, gracefully shutting down...")

		if err := server.Shutdown(); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
	}()
}
