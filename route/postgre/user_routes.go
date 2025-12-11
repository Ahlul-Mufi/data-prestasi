package postgreroute

import (
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router, userService servicepostgre.UserService) {
    userRoutes := api.Group("/users", mw.AuthMiddleware()) 
    
    userRoutes.Get("/", userService.GetUsers) 
    userRoutes.Get("/:id", userService.GetUserByID) 
    userRoutes.Post("/", userService.CreateUser) 
    userRoutes.Put("/:id", userService.UpdateUser) 
    userRoutes.Delete("/:id", userService.DeleteUser) 
    userRoutes.Put("/:id/role", userService.UpdateUserRole) 
}