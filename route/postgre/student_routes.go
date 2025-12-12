package postgreroute

import (
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupStudentRoutes(api fiber.Router, studentService servicepostgre.StudentService, userRepo repo.UserRepository) {
	studentGroup := api.Group("/students")

	studentGroup.Get("/", mw.AuthMiddleware(), studentService.GetAll)
	studentGroup.Get("/:id", mw.AuthMiddleware(), studentService.GetByID)
	studentGroup.Get("/:id/achievements", mw.AuthMiddleware(), studentService.GetAchievements)

	studentGroup.Post("/",
		mw.AuthMiddleware(),
		mw.RBACMiddleware(userRepo, "student:manage", "student:create"),
		studentService.CreateStudent,
	)

	studentGroup.Put("/:id",
		mw.AuthMiddleware(),
		mw.RBACMiddleware(userRepo, "student:manage", "student:update"),
		studentService.UpdateStudent,
	)

	studentGroup.Delete("/:id",
		mw.AuthMiddleware(),
		mw.RBACMiddleware(userRepo, "student:manage", "student:delete"),
		studentService.DeleteStudent,
	)

	studentGroup.Put("/:id/advisor",
		mw.AuthMiddleware(),
		mw.RBACMiddleware(userRepo, "student:manage", "student:update"),
		studentService.UpdateAdvisor,
	)
}
