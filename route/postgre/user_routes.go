package postgreroute

import (
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router, userService servicepostgre.UserService, userRepo repo.UserRepository) {
    
    authMiddleware := mw.AuthMiddleware()
    userRoutesCreate := api.Group("/users", 
        authMiddleware, 
        mw.RBACMiddleware(userRepo, "user:manage", "user:create"),
    )
    userRoutesCreate.Post("/", userService.CreateUser) 
    userRoutesRead := api.Group("/users", 
        authMiddleware, 
        mw.RBACMiddleware(userRepo, "user:manage", "user:read"),
    )
    userRoutesRead.Get("/", userService.GetUsers) 
    userRoutesRead.Get("/:id", userService.GetUserByID) 
    
    userRoutesUpdate := api.Group("/users", 
        authMiddleware, 
        mw.RBACMiddleware(userRepo, "user:manage", "user:update"),
    )
    userRoutesUpdate.Put("/:id", userService.UpdateUser) 
    userRoutesUpdate.Put("/:id/role", userService.UpdateUserRole) 

    userRoutesDelete := api.Group("/users", 
        authMiddleware, 
        mw.RBACMiddleware(userRepo, "user:manage", "user:delete"),
    )
    userRoutesDelete.Delete("/:id", userService.DeleteUser) 
}