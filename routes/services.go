package routes

import (
	"github.com/MrMikeDevTech/mrmikedevs-os/middleware"
	"github.com/MrMikeDevTech/mrmikedevs-os/utils"
	"github.com/gofiber/fiber/v2"
)

func ServicesRoutes(app *fiber.App) {
	services := app.Group("/services", middleware.JwtMiddleware)

	services.Get("/", func(c *fiber.Ctx) error {
		list := utils.GetSystemServices()
		return c.JSON(fiber.Map{
			"status": "success",
			"data":   list,
		})
	})

	services.Post("/action/:action/:service", func(c *fiber.Ctx) error {
		action := c.Params("action")
		serviceID := c.Params("service")

		err := utils.HandleServiceAction(serviceID, action)
		if err != nil {
			if err.Error() == "ALREADY_ACTIVE" {
				return c.Status(400).JSON(fiber.Map{
					"status":  "info",
					"message": "El servicio " + serviceID + " ya se encuentra encendido",
				})
			}
			if err.Error() == "ALREADY_INACTIVE" {
				return c.Status(400).JSON(fiber.Map{
					"status":  "info",
					"message": "El servicio " + serviceID + " ya se encuentra apagado",
				})
			}

			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Error al ejecutar la acción",
				"error":   err.Error(),
			})
		}

		messages := map[string]string{
			"stop":    "Servicio detenido correctamente",
			"start":   "Servicio iniciado correctamente",
			"restart": "Servicio reiniciado correctamente",
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": messages[action],
		})
	})

	services.Get("/logs/:service", func(c *fiber.Ctx) error {
		serviceID := c.Params("service")

		logs, err := utils.GetServiceLogs(serviceID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "No se pudieron obtener los logs",
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"service": serviceID,
			"logs":    logs,
		})
	})

	NginxRoutes(services)
}
