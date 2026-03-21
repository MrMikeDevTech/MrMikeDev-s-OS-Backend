package routes

import (
	"os"
	"time"

	"github.com/MrMikeDevTech/mrmikedevs-os/database"
	"github.com/MrMikeDevTech/mrmikedevs-os/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func AuthRoutes(app *fiber.App) {
	app.Post("/login", func(c *fiber.Ctx) error {
		type LoginInput struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var input LoginInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Revisa tu entrada", "data": err})
		}

		var user models.User
		if err := database.DB.Where(&models.User{Username: input.Username}).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Usuario o contraseña inválidos", "data": nil})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Usuario o contraseña inválidos", "data": nil})
		}

		claims := jwt.MapClaims{
			"user_id":  user.ID,
			"username": user.Username,
			"exp":      time.Now().Add(time.Hour * 72).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "No se pudo generar el token", "data": err})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Login exitoso",
			"token":   t,
			"user":    user,
		})
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		var user models.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Revisa tu entrada", "data": err})
		}

		if err := database.DB.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "No se pudo crear el usuario", "data": err})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Usuario creado",
			"user":    user,
		})
	})

	// app.Post("/logout", func(c *fiber.Ctx) error {
	// })

	// app.Post("/refresh", func(c *fiber.Ctx) error {
	// })
}
