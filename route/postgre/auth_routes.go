package postgreroute

import (
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(api fiber.Router, userService servicepostgre.UserService) {
    api.Post("/login", userService.Login)

    api.Get("/profile", mw.AuthMiddleware(), userService.Profile)
}
