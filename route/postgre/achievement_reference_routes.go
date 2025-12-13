package postgreroute

import (
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupAchievementReferenceRoutes(api fiber.Router, arService servicepostgre.AchievementReferenceService, userRepo repo.UserRepository) {
	arGroup := api.Group("/achievements")

	arGroup.Get("/:id", mw.AuthMiddleware(), arService.GetByID)

	arGroup.Post("/", mw.AuthMiddleware(), arService.Create)
	arGroup.Get("/me", mw.AuthMiddleware(), arService.GetMyAchievements)
	arGroup.Put("/:id", mw.AuthMiddleware(), arService.Update)
	arGroup.Delete("/:id", mw.AuthMiddleware(), arService.Delete)
	arGroup.Post("/:id/submit", mw.AuthMiddleware(), arService.Submit)

	arGroup.Get("/pending",
		mw.AuthMiddleware(),
		mw.RBACMiddleware(userRepo, "achievement:verify"),
		arService.GetPendingAchievements,
	)
	arGroup.Post("/:id/verify",
		mw.AuthMiddleware(),
		mw.RBACMiddleware(userRepo, "achievement:verify"),
		arService.Verify,
	)
	arGroup.Post("/:id/reject",
		mw.AuthMiddleware(),
		mw.RBACMiddleware(userRepo, "achievement:verify"),
		arService.Reject,
	)

	arGroup.Get("/",
		mw.AuthMiddleware(),
		mw.RBACMiddleware(userRepo, "achievement:read", "achievement:manage"),
		arService.GetAllAchievements,
	)

	arGroup.Get("/:id/history", mw.AuthMiddleware(), arService.GetHistory)
}
