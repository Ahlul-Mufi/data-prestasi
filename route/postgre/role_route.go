package postgreroute

import (
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupRoleRoutes(api fiber.Router, roleService servicepostgre.RoleService, userRepo repo.UserRepository) {
    authMiddleware := mw.AuthMiddleware()
    
    roleGroup := api.Group("/roles", authMiddleware)
    roleGroup.Get("/", 
        mw.RBACMiddleware(userRepo, "role:manage", "role:read"), 
        roleService.GetAll,
    )
    roleGroup.Get("/:id", 
        mw.RBACMiddleware(userRepo, "role:manage", "role:read"), 
        roleService.GetByID,
    )
    roleGroup.Post("/", 
        mw.RBACMiddleware(userRepo, "role:manage", "role:create"), 
        roleService.Create,
    )
    roleGroup.Put("/:id", 
        mw.RBACMiddleware(userRepo, "role:manage", "role:update"), 
        roleService.Update,
    )
    roleGroup.Delete("/:id", 
        mw.RBACMiddleware(userRepo, "role:manage", "role:delete"), 
        roleService.Delete,
    )
}