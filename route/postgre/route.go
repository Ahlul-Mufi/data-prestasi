package postgreroute

import (
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	userService interface{},
	userRepo interface{},
	roleService interface{},
	permissionService interface{},
	rolePermissionService interface{},
	achievementReferenceService interface{},
	studentService interface{},
	lecturerService interface{},
) {
	api := app.Group("/api/v1")

	us := userService.(servicepostgre.UserService)
	ur := userRepo.(repo.UserRepository)
	rs := roleService.(servicepostgre.RoleService)
	rps := rolePermissionService.(servicepostgre.RolePermissionService)
	ars := achievementReferenceService.(servicepostgre.AchievementReferenceService)
	ss := studentService.(servicepostgre.StudentService)
	ls := lecturerService.(servicepostgre.LecturerService)

	SetupAuthRoutes(api, us)
	SetupUserRoutes(api, us, ur)
	SetupRoleRoutes(api, rs)
	SetupRolePermissionRoutes(api, rps)
	SetupAchievementReferenceRoutes(api, ars)
	SetupStudentRoutes(api, ss, ur)
	SetupLecturerRoutes(api, ls, ur)
}
