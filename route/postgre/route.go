package postgreroute

import (
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userService interface{}, userRepo interface{}, permissionService interface{}, permissionRepo interface{}, rolePermissionService interface{}) {
    api := app.Group("/api/v1") 
    
    us := userService.(servicepostgre.UserService)
    ur := userRepo.(repo.UserRepository)
    ps := permissionService.(servicepostgre.PermissionService)
    rps := rolePermissionService.(servicepostgre.RolePermissionService)
    
    SetupAuthRoutes(api, us)
    SetupUserRoutes(api, us, ur) 
    SetupPermissionRoutes(api, ps, ur)
    SetupRolePermissionRoutes(api, rps, ur)
}