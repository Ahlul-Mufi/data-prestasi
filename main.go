package main

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Ahlul-Mufi/data-prestasi/database"

	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	service "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	route "github.com/Ahlul-Mufi/data-prestasi/route/postgre"
)

func main() {
    app := fiber.New()

    database.ConnectDB()

    r := repo.NewRoleRepository(database.DB)

    s := service.NewRoleService(r)

    route.SetupRoutes(app, s)

    app.Listen(":3000")
}
