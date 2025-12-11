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
    
    app := fiber.New(fiber.Config{
    })

    app.Use(cors.New())
    app.Use(logger.New())

    userRepo := repositorypostgre.NewUserRepository(database.DB)
    rolePermissionRepo := repositorypostgre.NewRolePermissionRepository(database.DB) 
    achievementReferenceRepo := repositorypostgre.NewAchievementReferenceRepository(database.DB) 

    userService := servicepostgre.NewUserService(userRepo)
    rolePermissionService := servicepostgre.NewRolePermissionService(rolePermissionRepo)
    achievementReferenceService := servicepostgre.NewAchievementReferenceService(achievementReferenceRepo)

    postgreroute.SetupRoutes(
        app, 
        userService, 
        rolePermissionService, 
        achievementReferenceService,
    ) 

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }
    log.Fatal(app.Listen(":" + port))
}