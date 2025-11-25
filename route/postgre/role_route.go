package postgreroute

import (
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, roleService servicepostgre.RoleService) {
    role := app.Group("/role")

    role.Get("/", roleService.GetAll)
    role.Get("/:id", roleService.GetByID)
    role.Post("/", roleService.Create)
    role.Put("/:id", roleService.Update)
    role.Delete("/:id", roleService.Delete)
}
