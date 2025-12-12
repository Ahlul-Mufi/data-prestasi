package postgreroute

import (
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupLecturerRoutes(api fiber.Router, lecturerService servicepostgre.LecturerService) {
	lecturerGroup := api.Group("/lecturers")
	lecturerGroup.Get("/", mw.AuthMiddleware(), lecturerService.GetAll)
	lecturerGroup.Get("/:id/advisees", mw.AuthMiddleware(), lecturerService.GetAdvisees)
}