package main

import (
	"log"
	"os"

	_ "github.com/Ahlul-Mufi/data-prestasi/docs"

	repomongo "github.com/Ahlul-Mufi/data-prestasi/app/repository/mongo"
	repositorypostgre "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	"github.com/Ahlul-Mufi/data-prestasi/database"
	postgreroute "github.com/Ahlul-Mufi/data-prestasi/route/postgre"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Data Prestasi API
// @version 1.0
// @description API Backend Data Prestasi (PostgreSQL & MongoDB)
// @host localhost:3000
// @BasePath /
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	database.ConnectDB()
	defer database.DisconnectDB()

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	userRepo := repositorypostgre.NewUserRepository(database.DB)
	roleRepo := repositorypostgre.NewRoleRepository(database.DB)
	permissionRepo := repositorypostgre.NewPermissionRepository(database.DB)
	rolePermissionRepo := repositorypostgre.NewRolePermissionRepository(database.DB)
	achievementReferenceRepo := repositorypostgre.NewAchievementReferenceRepository(database.DB)
	studentRepo := repositorypostgre.NewStudentRepository(database.DB)
	lecturerRepo := repositorypostgre.NewLecturerRepository(database.DB)
	statisticsRepo := repositorypostgre.NewStatisticsRepository(database.DB)

	achievementMongoRepo := repomongo.NewAchievementRepository()

	userService := servicepostgre.NewUserService(userRepo)
	roleService := servicepostgre.NewRoleService(roleRepo)
	permissionService := servicepostgre.NewPermissionService(permissionRepo)
	rolePermissionService := servicepostgre.NewRolePermissionService(rolePermissionRepo)

	achievementReferenceService := servicepostgre.NewAchievementReferenceService(
		achievementReferenceRepo,
		achievementMongoRepo,
		studentRepo,
		lecturerRepo,
	)

	studentService := servicepostgre.NewStudentService(studentRepo, achievementReferenceRepo, userRepo, lecturerRepo)
	lecturerService := servicepostgre.NewLecturerService(lecturerRepo, userRepo)

	statisticsService := servicepostgre.NewStatisticsService(
		statisticsRepo,
		achievementMongoRepo,
		userRepo,
	)

	postgreroute.SetupRoutes(
		app,
		userService,
		userRepo,
		roleService,
		permissionService,
		rolePermissionService,
		achievementReferenceService,
		studentService,
		lecturerService,
		statisticsService,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("Swagger UI tersedia di http://localhost:3000/swagger/index.html")
	log.Fatal(app.Listen(":" + port))
}
