package main

import (
	"log"
	"os"

	repositorypostgre "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	"github.com/Ahlul-Mufi/data-prestasi/database"
	postgreroute "github.com/Ahlul-Mufi/data-prestasi/route/postgre"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
    database.ConnectDB()
    defer func() {
        if err := database.DB.Close(); err != nil {
            log.Printf("Error closing database connection: %v", err)
        }
    }()
    
    app := fiber.New()
    
    app.Use(logger.New())
    app.Use(cors.New())

    userRepo := repositorypostgre.NewUserRepository(database.DB)
    permissionRepo := repositorypostgre.NewPermissionRepository(database.DB)
    rolePermissionRepo := repositorypostgre.NewRolePermissionRepository(database.DB)

    userService := servicepostgre.NewUserService(userRepo)
    permissionService := servicepostgre.NewPermissionService(permissionRepo)
    rolePermissionService := servicepostgre.NewRolePermissionService(rolePermissionRepo)

    postgreroute.SetupRoutes(
        app, 
        userService, 
        userRepo, 
        permissionService, 
        permissionRepo,
        rolePermissionService,
    )
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    log.Printf("Server is running on port :%s", port)
    log.Fatal(app.Listen(":" + port))
}