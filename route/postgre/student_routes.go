package postgreroute

import (
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupStudentRoutes(api fiber.Router, studentService servicepostgre.StudentService) {
	studentGroup := api.Group("/students")

	studentGroup.Get("/", mw.AuthMiddleware(), studentService.GetAll)
	studentGroup.Get("/:id", mw.AuthMiddleware(), studentService.GetByID)
	studentGroup.Get("/:id/achievements", mw.AuthMiddleware(), studentService.GetAchievements)
	studentGroup.Put("/:id/advisor", mw.AuthMiddleware(), studentService.UpdateAdvisor)
}