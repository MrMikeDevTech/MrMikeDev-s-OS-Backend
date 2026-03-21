package main

import (
	"log"

	"github.com/MrMikeDevTech/mrmikedevs-os/database"
	"github.com/MrMikeDevTech/mrmikedevs-os/middleware"
	"github.com/MrMikeDevTech/mrmikedevs-os/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file or file not found")
	}

	database.ConnectDB()

	app := fiber.New()

	app.Use(middleware.CORSMiddleware)
	app.Use("/ws", middleware.WSUpgrade)

	routes.Router(app)

	log.Fatal(app.Listen(":15800"))
}
