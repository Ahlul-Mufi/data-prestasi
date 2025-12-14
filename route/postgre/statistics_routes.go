package postgreroute

import (
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	mw "github.com/Ahlul-Mufi/data-prestasi/middleware/postgre"
	"github.com/gofiber/fiber/v2"
)

func SetupStatisticsRoutes(api fiber.Router, statsService servicepostgre.StatisticsService) {
	reportsGroup := api.Group("/reports")

	reportsGroup.Get("/statistics",
		mw.AuthMiddleware(),
		statsService.GetStatistics,
	)

	reportsGroup.Get("/student/:id",
		mw.AuthMiddleware(),
		statsService.GetStudentStatistics,
	)
}
