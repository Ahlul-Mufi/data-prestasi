package middleware

import (
	"strings"

	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	helper "github.com/Ahlul-Mufi/data-prestasi/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RBACMiddleware(userRepo repo.UserRepository, requiredPermissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleIDStr, ok := c.Locals("role_id").(string)
		if !ok || roleIDStr == "" {
			return helper.ErrorResponse(c, fiber.StatusForbidden, "Access Denied: Role information missing.", "")
		}

		roleID, err := uuid.Parse(roleIDStr)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid Role ID format in context.", err.Error())
		}

		userPermissions, err := userRepo.GetPermissions(roleID)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve user permissions.", err.Error())
		}

		hasPermission := false
		permissionsMap := make(map[string]bool)
		for _, p := range userPermissions {
			permissionsMap[p] = true
		}

		for _, requiredPerm := range requiredPermissions {
			if permissionsMap[requiredPerm] {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return helper.ErrorResponse(c, fiber.StatusForbidden, 
				"Access Denied: Insufficient permissions.", 
				"Required: " + strings.Join(requiredPermissions, " OR "))
		}

		return c.Next()
	}
}