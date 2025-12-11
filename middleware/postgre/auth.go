package middleware

import (
	"strings"

	helper "github.com/Ahlul-Mufi/data-prestasi/helper"
	"github.com/gofiber/fiber/v2"

	utils "github.com/Ahlul-Mufi/data-prestasi/utils/postgre"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Missing authorization header", "")
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid authorization header format", "Expected: Bearer <token>")
		}
		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token", err.Error())
		}
		
		c.Locals("user_id", claims.UserID.String())

        if claims.RoleID != nil {
            c.Locals("role_id", claims.RoleID.String())
        }
        
		return c.Next()
	}
}