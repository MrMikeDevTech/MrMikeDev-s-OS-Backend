package middleware

import (
	"strings"

	"github.com/MrMikeDevTech/mrmikedevs-os/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
)

func WSUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func CORSMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:15805, http://mrmikedev-vps:15805, http://mrmikedev-vps",
		AllowHeaders:     "Origin, Content-Type, Accept, x-api-key, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	})
}

func ApiKeyMiddleware(c *fiber.Ctx) error {
	apiKey := utils.GetEnv("API_KEY")
	if apiKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "API_KEY no configurado en el entorno"})
	}

	clientKey := c.Get("x-api-key")
	if clientKey == "" {
		clientKey = c.Query("x-api-key")
	}

	if clientKey != apiKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "API Key inválida"})
	}

	return c.Next()
}

func JwtMiddleware(c *fiber.Ctx) error {
	var tokenString string

	authHeader := c.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = authHeader[7:]
	} else if c.Query("token") != "" {
		tokenString = c.Query("token")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Token inválido",
		})
	}

	secret := []byte(utils.GetEnv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Sesión inválida o expirada",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.Locals("user_id", claims["user_id"])
		c.Locals("username", claims["username"])
	}

	return c.Next()
}
