package main

import (
	"log"

	"github.com/MrMikeDevTech/mrmikedevs-os/middleware"
	"github.com/MrMikeDevTech/mrmikedevs-os/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use(middleware.CORSMiddleware)
	app.Use("/ws", middleware.WSUpgrade)

	routes.SystemRoutes(app)

	log.Fatal(app.Listen(":15800"))
}
