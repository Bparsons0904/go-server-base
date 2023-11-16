package routes

import (
	"log"
	"runtime"
	"time"

	"github.com/bparsons094/go-server-base/database"
	"github.com/bparsons094/go-server-base/middleware"
	"github.com/bparsons094/go-server-base/utils"
	"github.com/bparsons094/go-server-base/websockets"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/websocket/v2"
)

var (
	startTime = time.Now().UTC()
)

func SetupRoutes(app *fiber.App, config utils.Config) {
	DB := database.GetDatabase()

	HealthRoutes(app)

	// Websocket routes
	service := websockets.NewWebSocketService()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("service", service)
		return c.Next()
	})
	app.Use("/ws", websocket.New(func(c *websocket.Conn) {
		service.HandleWebSocketConnection(c)
	}))

	// Internal routes
	api := app.Group("/api")
	api.Use(middleware.AuthenticateUser)
	api.Use(compress.New(compress.Config{
		Level: compress.LevelDefault,
	}))
	api.Use(middleware.LogMiddleware(DB))

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Route Not found"})
	})

}

func HealthRoutes(app *fiber.App) {
	app.Get("/health", getHealth)
	app.Get("/health/monitor", monitor.New(monitor.Config{
		Title: "Server Health Monitor",
	}))
}

func getHealth(c *fiber.Ctx) error {

	type Health struct {
		Uptime        string `json:"uptime"`
		AppVersion    string `json:"app_version"`
		MemoryUsage   uint64 `json:"memory_usage"`
		NumGoroutine  int    `json:"num_goroutine"`
		NumCPU        int    `json:"num_cpu"`
		DatabaseAlive bool   `json:"database_alive"`
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	sqlDB, err := database.GetDatabase().DB()
	if err != nil {
		log.Println("Error getting database: ", err)
	}
	dbAlive := sqlDB.Ping() == nil

	health := Health{
		Uptime:        time.Since(startTime).String(),
		AppVersion:    utils.GetEnv("VERSION"),
		MemoryUsage:   memStats.Alloc,
		NumGoroutine:  runtime.NumGoroutine(),
		NumCPU:        runtime.NumCPU(),
		DatabaseAlive: dbAlive,
	}

	return c.JSON(health)
}
