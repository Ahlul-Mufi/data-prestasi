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
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            code := fiber.StatusInternalServerError
            if e, ok := err.(*fiber.Error); ok {
                code = e.Code
            }
            c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
            return c.Status(code).JSON(fiber.Map{
                "error": err.Error(),
            })
        },
    })
    
    app.Use(logger.New())
    app.Use(cors.New())

    userRepo := repositorypostgre.NewUserRepository(database.DB)
    userService := servicepostgre.NewUserService(userRepo)

    postgreroute.SetupRoutes(app, userService)
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    log.Printf("Server is running on port :%s", port)
    log.Fatal(app.Listen(":" + port))
}