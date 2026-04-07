package routes

import (
	"time"

	"github.com/MrMikeDevTech/mrmikedevs-os/middleware"
	"github.com/MrMikeDevTech/mrmikedevs-os/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func WebsocketTunnelRoutes(app *fiber.App) {
	ws := app.Group("/ws", middleware.JwtMiddleware)

	ws.Get("/system", websocket.New(func(c *websocket.Conn) {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			metrics := fiber.Map{
				"cpu":     utils.GetCPU(),
				"ram":     utils.GetRAM(),
				"storage": utils.GetDisks(),
				"network": utils.GetNetwork(time.Second),
			}

			services := utils.GetSystemServices()

			payload := fiber.Map{
				"event": "system_update",
				"data": fiber.Map{
					"metrics":  metrics,
					"services": services,
				},
			}

			if err := c.WriteJSON(payload); err != nil {
				return
			}

			<-ticker.C
		}
	}))
}
