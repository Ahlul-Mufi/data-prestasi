package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	utils "github.com/Ahlul-Mufi/data-prestasi/utils/postgre"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format. Expected: Bearer <token>",
			})
		}
		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}
		c.Locals("user_id", claims.UserID.String())
		return c.Next()
	}
}