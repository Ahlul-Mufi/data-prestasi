package postgreroute

import (
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupPermissionRoutes(api fiber.Router, permissionService servicepostgre.PermissionService, userRepo repo.UserRepository) {
    authMiddleware := mw.AuthMiddleware()
    
    permissionRoutes := api.Group("/permissions", 
        authMiddleware,
        mw.RBACMiddleware(userRepo, "permission:manage"),
    )
    
    permissionRoutes.Post("/", permissionService.Create)
    permissionRoutes.Get("/", permissionService.GetAll)
    permissionRoutes.Get("/:id", permissionService.GetByID)
    permissionRoutes.Delete("/:id", permissionService.Delete)
}