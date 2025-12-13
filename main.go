package main

import (
	"log"
	"os"

	repomongo "github.com/Ahlul-Mufi/data-prestasi/app/repository/mongo"
	repositorypostgre "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	servicepostgre "github.com/Ahlul-Mufi/data-prestasi/app/service/postgre"
	"github.com/Ahlul-Mufi/data-prestasi/database"
	postgreroute "github.com/Ahlul-Mufi/data-prestasi/route/postgre"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	database.ConnectDB()
	defer database.DisconnectDB()

	app := fiber.New(fiber.Config{})

	app.Use(cors.New())
	app.Use(logger.New())

	userRepo := repositorypostgre.NewUserRepository(database.DB)
	roleRepo := repositorypostgre.NewRoleRepository(database.DB)
	permissionRepo := repositorypostgre.NewPermissionRepository(database.DB)
	rolePermissionRepo := repositorypostgre.NewRolePermissionRepository(database.DB)
	achievementReferenceRepo := repositorypostgre.NewAchievementReferenceRepository(database.DB)
	studentRepo := repositorypostgre.NewStudentRepository(database.DB)
	lecturerRepo := repositorypostgre.NewLecturerRepository(database.DB)

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
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
