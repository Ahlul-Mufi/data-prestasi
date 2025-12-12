package postgreroute

import (
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupLecturerRoutes(api fiber.Router, lecturerService servicepostgre.LecturerService, userRepo repo.UserRepository) {
	lecturerGroup := api.Group("/lecturers")

	lecturerGroup.Get("/", mw.AuthMiddleware(), lecturerService.GetAll)
	lecturerGroup.Get("/:id", mw.AuthMiddleware(), lecturerService.GetByID)
	lecturerGroup.Get("/:id/advisees", mw.AuthMiddleware(), lecturerService.GetAdvisees)

	lecturerGroup.Post("/",
		mw.AuthMiddleware(),
		mw.RBACMiddleware(userRepo, "lecturer:manage", "lecturer:create"),
		lecturerService.CreateLecturer,
	)

	lecturerGroup.Put("/:id",
		mw.AuthMiddleware(),
		mw.RBACMiddleware(userRepo, "lecturer:manage", "lecturer:update"),
		lecturerService.UpdateLecturer,
	)

	lecturerGroup.Delete("/:id",
		mw.AuthMiddleware(),
		mw.RBACMiddleware(userRepo, "lecturer:manage", "lecturer:delete"),
		lecturerService.DeleteLecturer,
	)
}
