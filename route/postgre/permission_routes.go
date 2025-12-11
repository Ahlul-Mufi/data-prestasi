package postgreroute

import (
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupPermissionRoutes(api fiber.Router, permissionService servicepostgre.PermissionService, userRepo repo.UserRepository) {
    authMiddleware := mw.AuthMiddleware()
    
    // Semua rute Permission hanya dapat diakses oleh Admin (membutuhkan permission:manage)
    permissionRoutes := api.Group("/permissions", 
        authMiddleware,
        mw.RBACMiddleware(userRepo, "permission:manage"),
    )
    
    // POST /api/v1/permissions (Create)
    permissionRoutes.Post("/", permissionService.Create)
    // GET /api/v1/permissions (Get All)
    permissionRoutes.Get("/", permissionService.GetAll)
    // GET /api/v1/permissions/:id (Get By ID)
    permissionRoutes.Get("/:id", permissionService.GetByID)
    // PUT /api/v1/permissions/:id (Update)
    permissionRoutes.Put("/:id", permissionService.Update)
    // DELETE /api/v1/permissions/:id (Delete)
    permissionRoutes.Delete("/:id", permissionService.Delete)
}