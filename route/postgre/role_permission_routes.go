package postgreroute

import (
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupRolePermissionRoutes(api fiber.Router, rpService servicepostgre.RolePermissionService) {
	rolePermissionGroup := api.Group("/role-permissions")
    rolePermissionGroup.Post("/", mw.AuthMiddleware(), rpService.Add)
    rolePermissionGroup.Delete("/", mw.AuthMiddleware(), rpService.Remove)
}