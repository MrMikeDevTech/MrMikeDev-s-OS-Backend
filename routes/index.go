package routes

import (
	"github.com/MrMikeDevTech/mrmikedevs-os/middleware"
	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	app.Use(middleware.ApiKeyMiddleware)

	app.Get("/", middleware.JwtMiddleware, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"endpoints": []string{
				"GET /health",
				"WS /ws/system",
				"POST /auth/login",
				"POST /auth/register",
				"GET /auth/validate",
				"POST /auth/logout",
				"POST /auth/refresh",
			},
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	AuthRoutes(app)

	SystemRoutes(app)
}
