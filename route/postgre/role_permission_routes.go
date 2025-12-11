package postgreroute

import (
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupRolePermissionRoutes(api fiber.Router, rpService servicepostgre.RolePermissionService, userRepo repo.UserRepository) {
	authMiddleware := mw.AuthMiddleware()
	rolePermissionRoutes := api.Group("/role-permissions", 
		authMiddleware,
		mw.RBACMiddleware(userRepo, "role:manage"),
	)
	rolePermissionRoutes.Post("/", rpService.AddPermissionToRole)
	rolePermissionRoutes.Get("/role/:roleId", rpService.GetPermissionsByRoleID) 
	rolePermissionRoutes.Delete("/role/:roleId/permission/:permissionId", rpService.RemovePermissionFromRole) 
}