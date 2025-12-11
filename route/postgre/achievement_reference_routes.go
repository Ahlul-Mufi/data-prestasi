package postgreroute

import (
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupAchievementReferenceRoutes(api fiber.Router, arService servicepostgre.AchievementReferenceService) {
	arGroup := api.Group("/achievement-references")

	arGroup.Post("/", mw.AuthMiddleware(), arService.Create) 
	arGroup.Get("/me", mw.AuthMiddleware(), arService.GetMyAchievements) 
	arGroup.Put("/:id", mw.AuthMiddleware(), arService.Update) 
	arGroup.Delete("/:id", mw.AuthMiddleware(), arService.Delete) 


	arGroup.Get("/pending", mw.AuthMiddleware(), arService.GetPendingAchievements) 
	arGroup.Post("/:id/verify", mw.AuthMiddleware(), arService.Verify) 
	arGroup.Post("/:id/reject", mw.AuthMiddleware(), arService.Reject) 
    
	
	arGroup.Get("/", mw.AuthMiddleware(), arService.GetAllAchievements) 
}