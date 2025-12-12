package postgreroute

import (
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupRoleRoutes(api fiber.Router, roleService servicepostgre.RoleService) {
	roleGroup := api.Group("/roles")
	
	roleGroup.Get("/", mw.AuthMiddleware(), roleService.GetAll)
	roleGroup.Get("/:id", mw.AuthMiddleware(), roleService.GetByID)
	roleGroup.Post("/", mw.AuthMiddleware(), roleService.Create)
	roleGroup.Put("/:id", mw.AuthMiddleware(), roleService.Update)
	roleGroup.Delete("/:id", mw.AuthMiddleware(), roleService.Delete)
}