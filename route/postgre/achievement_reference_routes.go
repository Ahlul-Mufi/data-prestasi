package postgreroute

import (
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupAchievementReferenceRoutes(api fiber.Router, arService servicepostgre.AchievementReferenceService) {
	arGroup := api.Group("/achievements")

	arGroup.Get("/", mw.AuthMiddleware(), arService.GetAllAchievements)
	arGroup.Get("/:id", mw.AuthMiddleware(), arService.GetByID)
	arGroup.Post("/", mw.AuthMiddleware(), arService.Create)
	arGroup.Put("/:id", mw.AuthMiddleware(), arService.Update)
	arGroup.Delete("/:id", mw.AuthMiddleware(), arService.Delete)

	arGroup.Post("/:id/submit", mw.AuthMiddleware(), arService.Submit)
	arGroup.Post("/:id/verify", mw.AuthMiddleware(), arService.Verify)
	arGroup.Post("/:id/reject", mw.AuthMiddleware(), arService.Reject)
	arGroup.Get("/:id/history", mw.AuthMiddleware(), arService.GetHistory)

	arGroup.Get("/pending", mw.AuthMiddleware(), arService.GetPendingAchievements)
	arGroup.Get("/me", mw.AuthMiddleware(), arService.GetMyAchievements)
}
