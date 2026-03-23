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
	app.Post("/auth/login", func(c *fiber.Ctx) error {
		type LoginInput struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var input LoginInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Cuerpo de la petición inválido",
			})
		}

		var identity string
		if input.Email != "" {
			identity = input.Email
		} else {
			identity = input.Username
		}

		if identity == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Debes proporcionar un usuario o un correo electrónico",
			})
		}

		var user models.User
		if err := database.DB.Where("username = ? OR email = ?", identity, identity).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Credenciales inválidas",
			})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Credenciales inválidas",
			})
		}

		claims := jwt.MapClaims{
			"user_id":  user.ID,
			"username": user.Username,
			"exp":      time.Now().Add(time.Hour * 5).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Error al generar el token de acceso",
			})
		}

		user.Password = ""
		return c.JSON(fiber.Map{
			"status": "success",
			"token":  t,
			"user":   user,
		})
	})

	app.Post("/auth/register", func(c *fiber.Ctx) error {
		type RegisterInput struct {
			Username        string `json:"username"`
			FullName        string `json:"full_name"`
			Email           string `json:"email"`
			Password        string `json:"password"`
			ConfirmPassword string `json:"confirm_password"`
		}

		var input RegisterInput

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Datos de entrada inválidos",
			})
		}

		if len(input.Password) < 8 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "La contraseña debe tener al menos 8 caracteres",
			})
		}

		if input.Password != input.ConfirmPassword {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Las contraseñas no coinciden",
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Error al procesar la seguridad",
			})
		}

		user := models.User{
			Username: input.Username,
			Email:    input.Email,
			Password: string(hashedPassword),
		}

		if err := database.DB.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status":  "error",
				"message": "El nombre de usuario o email ya están registrados",
			})
		}

		user.Password = ""
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Usuario registrado exitosamente en MrMikeDev-s-OS",
			"user":    user,
		})
	})

	app.Get("/auth/validate", func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "No hay token"})
		}

		tokenString := authHeader[7:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Token inválido o expirado"})
		}

		claims := token.Claims.(jwt.MapClaims)

		return c.JSON(fiber.Map{
			"status":   "success",
			"message":  "Sesión válida",
			"user_id":  claims["user_id"],
			"username": claims["username"],
		})
	})

	app.Post("/auth/refresh", func(c *fiber.Ctx) error {
		type RefreshInput struct {
			Token string `json:"token"`
		}
		var input RefreshInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Token requerido"})
		}

		token, err := jwt.Parse(input.Token, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil && !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Token inválido"})
		}

		claims := token.Claims.(jwt.MapClaims)

		newClaims := jwt.MapClaims{
			"user_id":  claims["user_id"],
			"username": claims["username"],
			"exp":      time.Now().Add(time.Hour * 72).Unix(),
		}

		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
		t, err := newToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Error al renovar"})
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"token":  t,
		})
	})

	app.Post("/auth/logout", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Saliendo del sistema... borra el token en el cliente",
		})
	})
}
