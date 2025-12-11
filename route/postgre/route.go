package postgreroute

import (
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App, 
	userService interface{}, 
	rolePermissionService interface{}, 
    achievementReferenceService interface{},
) {
    api := app.Group("/api/v1")
    us := userService.(servicepostgre.UserService)
    rps := rolePermissionService.(servicepostgre.RolePermissionService)
    ars := achievementReferenceService.(servicepostgre.AchievementReferenceService)

    SetupAuthRoutes(api, us)
    SetupRolePermissionRoutes(api, rps)
    SetupAchievementReferenceRoutes(api, ars)
}