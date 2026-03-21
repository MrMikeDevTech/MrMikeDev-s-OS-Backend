package routes

import (
	"time"

	"github.com/MrMikeDevTech/mrmikedevs-os/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SystemRoutes(app *fiber.App) {

	app.Get("/ws/system", websocket.New(func(c *websocket.Conn) {
		for {
			metrics := fiber.Map{
				"cpu":     utils.GetCPU(),
				"ram":     utils.GetRAM(),
				"storage": utils.GetDisks(),
				"network": utils.GetNetwork(time.Second),
			}

			if err := c.WriteJSON(metrics); err != nil {
				break
			}
		}
	}))
}
