package servicepostgre

import (
	modelmongo "github.com/Ahlul-Mufi/data-prestasi/app/model/mongo"
	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repomongo "github.com/Ahlul-Mufi/data-prestasi/app/repository/mongo"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	helper "github.com/Ahlul-Mufi/data-prestasi/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatisticsService interface {
	GetStatistics(c *fiber.Ctx) error
	GetStudentStatistics(c *fiber.Ctx) error
}

type statisticsService struct {
	statsRepo repo.StatisticsRepository
	mongoRepo repomongo.AchievementRepository
	userRepo  repo.UserRepository
}

func NewStatisticsService(
	statsRepo repo.StatisticsRepository,
	mongoRepo repomongo.AchievementRepository,
	userRepo repo.UserRepository,
) StatisticsService {
	return &statisticsService{
		statsRepo: statsRepo,
		mongoRepo: mongoRepo,
		userRepo:  userRepo,
	}
}

func (s *statisticsService) GetStatistics(c *fiber.Ctx) error {
	totalAchievements, err := s.statsRepo.GetTotalAchievements(nil)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get total achievements", err.Error())
	}

	achievementsByType, err := s.statsRepo.GetAchievementsByType(nil)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get achievements by type", err.Error())
	}

	achievementsByStatus, err := s.statsRepo.GetAchievementsByStatus(nil)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get achievements by status", err.Error())
	}

	achievementsByPeriod, err := s.statsRepo.GetAchievementsByPeriod(nil, "monthly")
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get achievements by period", err.Error())
	}

	topStudents, err := s.statsRepo.GetTopStudents(10)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get top students", err.Error())
	}

	competitionDistribution := []m.CompetitionLevelStatistic{
		{Level: "international", Count: 0},
		{Level: "national", Count: 0},
		{Level: "regional", Count: 0},
		{Level: "local", Count: 0},
	}

	response := m.StatisticsResponse{
		TotalAchievements:       totalAchievements,
		AchievementsByType:      achievementsByType,
		AchievementsByStatus:    achievementsByStatus,
		AchievementsByPeriod:    achievementsByPeriod,
		TopStudents:             topStudents,
		CompetitionDistribution: competitionDistribution,
	}

	return helper.SuccessResponse(c, fiber.StatusOK, response)
}

func (s *statisticsService) GetStudentStatistics(c *fiber.Ctx) error {
	idStr := c.Params("id")
	studentID, err := uuid.Parse(idStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid student ID format", err.Error())
	}

	student, err := s.statsRepo.GetStudentInfo(studentID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Student not found", err.Error())
	}

	user, err := s.userRepo.FindByID(student.UserID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get user info", err.Error())
	}

	totalAchievements, err := s.statsRepo.GetTotalAchievements(&studentID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get total achievements", err.Error())
	}

	achievementsByType, err := s.statsRepo.GetAchievementsByType(&studentID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get achievements by type", err.Error())
	}

	achievementsByStatus, err := s.statsRepo.GetAchievementsByStatus(&studentID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get achievements by status", err.Error())
	}

	recentRefs, err := s.statsRepo.GetRecentAchievements(studentID, 5)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get recent achievements", err.Error())
	}

	var mongoIDs []primitive.ObjectID
	for _, ref := range recentRefs {
		objID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
		if err == nil {
			mongoIDs = append(mongoIDs, objID)
		}
	}

	achievements, _ := s.mongoRepo.GetMultipleByIDs(mongoIDs)
	achievementMap := make(map[string]modelmongo.Achievement)
	for _, ach := range achievements {
		achievementMap[ach.ID.Hex()] = ach
	}

	var recentAchievements []m.RecentAchievement
	totalPoints := 0

	for _, ref := range recentRefs {
		if ach, exists := achievementMap[ref.MongoAchievementID]; exists {
			recentAchievements = append(recentAchievements, m.RecentAchievement{
				ID:              ref.ID.String(),
				Title:           ach.Title,
				AchievementType: string(ach.AchievementType),
				Status:          string(ref.Status),
				Points:          ach.Points,
				CreatedAt:       ref.CreatedAt,
			})
			if ref.Status == m.StatusVerified {
				totalPoints += ach.Points
			}
		}
	}

	monthlyTrend, err := s.statsRepo.GetAchievementsByPeriod(&studentID, "monthly")
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get monthly trend", err.Error())
	}

	var monthlyAchievements []m.MonthlyAchievement
	for _, stat := range monthlyTrend {
		monthlyAchievements = append(monthlyAchievements, m.MonthlyAchievement{
			Month: stat.Period,
			Count: stat.Count,
		})
	}

	response := m.StudentStatisticsResponse{
		StudentID:            student.StudentID,
		StudentName:          user.FullName,
		ProgramStudy:         student.ProgramStudy,
		AcademicYear:         student.AcademicYear,
		TotalAchievements:    totalAchievements,
		TotalPoints:          totalPoints,
		AchievementsByType:   achievementsByType,
		AchievementsByStatus: achievementsByStatus,
		RecentAchievements:   recentAchievements,
		MonthlyTrend:         monthlyAchievements,
	}

	return helper.SuccessResponse(c, fiber.StatusOK, response)
}
