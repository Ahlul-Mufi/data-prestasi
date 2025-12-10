package postgreroute

import (
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userService interface{}) {
    api := app.Group("/api/v1")
    us := userService.(servicepostgre.UserService)

    SetupAuthRoutes(api, us)
    SetupUserRoutes(api, us) 
}