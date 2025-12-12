package postgreroute

import (
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App, 
	userService interface{}, 
	roleService interface{},
	permissionService interface{},
	rolePermissionService interface{}, 
	achievementReferenceService interface{}, 
	studentService interface{},
	lecturerService interface{},
) {
    api := app.Group("/api/v1")
    us := userService.(servicepostgre.UserService)
    rs := roleService.(servicepostgre.RoleService)
    rps := rolePermissionService.(servicepostgre.RolePermissionService)
    ars := achievementReferenceService.(servicepostgre.AchievementReferenceService) 
    ss := studentService.(servicepostgre.StudentService)
    ls := lecturerService.(servicepostgre.LecturerService)

    SetupAuthRoutes(api, us)
    SetupRoleRoutes(api, rs) 
    SetupRolePermissionRoutes(api, rps)
    SetupAchievementReferenceRoutes(api, ars) 
    SetupStudentRoutes(api, ss) 
    SetupLecturerRoutes(api, ls) 
}