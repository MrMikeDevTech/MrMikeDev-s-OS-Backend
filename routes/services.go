package routes

import (
	"github.com/MrMikeDevTech/mrmikedevs-os/middleware"
	"github.com/MrMikeDevTech/mrmikedevs-os/utils"
	"github.com/gofiber/fiber/v2"
)

func ServicesRoutes(app *fiber.App) {
	services := app.Group("/services", middleware.JwtMiddleware)

	nginx := services.Group("/nginx", middleware.JwtMiddleware)

	nginx.Get("/config", func(c *fiber.Ctx) error {
		content, err := utils.ReadNginxConfig()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "No se pudo leer el archivo"})
		}
		return c.JSON(fiber.Map{"status": "success", "content": content})
	})

	nginx.Post("/test", func(c *fiber.Ctx) error {
		type Request struct {
			Content string `json:"content"`
		}
		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Cuerpo inválido"})
		}

		output, err := utils.TestNginxConfig(req.Content)
		if err != nil {
			return c.Status(422).JSON(fiber.Map{
				"status":  "error",
				"message": "Sintaxis de Nginx inválida",
				"output":  output,
			})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Configuración válida", "output": output})
	})

	nginx.Post("/save", func(c *fiber.Ctx) error {
		type Request struct {
			Content string `json:"content"`
		}
		var req Request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Cuerpo inválido"})
		}

		_, err := utils.TestNginxConfig(req.Content)
		if err != nil {
			return c.Status(422).JSON(fiber.Map{"status": "error", "message": "No se puede guardar: Sintaxis inválida"})
		}

		if err := utils.SaveNginxConfig(req.Content); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error al guardar el archivo"})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Configuración guardada y aplicada"})
	})
}
