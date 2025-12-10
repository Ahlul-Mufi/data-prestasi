package postgreroute

import (
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupPermissionRoutes(app *fiber.App, permissionService servicepostgre.PermissionService) {
	permission := app.Group("/permission")
	permission.Get("/", permissionService.GetAll)
	permission.Get("/:id", permissionService.GetByID)	
	permission.Post("/", permissionService.Create)
	permission.Put("/:id", permissionService.Update)
	permission.Delete("/:id", permissionService.Delete)
}