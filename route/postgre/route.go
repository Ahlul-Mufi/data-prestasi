package postgreroute

import (
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userService interface{}) {
    api := app.Group("/api")
    us := userService.(servicepostgre.UserService)

    SetupAuthRoutes(api, us)
}
