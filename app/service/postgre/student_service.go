package servicepostgre

import (
	"errors"

	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	helper "github.com/Ahlul-Mufi/data-prestasi/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StudentService interface {
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetAchievements(c *fiber.Ctx) error
	UpdateAdvisor(c *fiber.Ctx) error
}

type studentService struct {
	repo repo.StudentRepository
	arRepo repo.AchievementReferenceRepository 
}

func NewStudentService(r repo.StudentRepository, arRepo repo.AchievementReferenceRepository) StudentService {
	return &studentService{r, arRepo}
}

func (s *studentService) GetAll(c *fiber.Ctx) error {
	students, err := s.repo.GetAll()
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch students", err.Error())
	}
	return helper.SuccessResponse(c, fiber.StatusOK, students)
}

func (s *studentService) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", "ID must be a valid UUID")
	}

	student, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, errors.New("student not found")) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", err.Error())
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch student", err.Error())
	}
	return helper.SuccessResponse(c, fiber.StatusOK, student)
}

func (s *studentService) GetAchievements(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", "ID must be a valid UUID")
	}

	_, err = s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, errors.New("student not found")) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "Student not found")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check student", err.Error())
	}
	
	achievements, err := s.arRepo.FindByStudentID(id)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch achievements", err.Error())
	}
	return helper.SuccessResponse(c, fiber.StatusOK, achievements)
}

func (s *studentService) UpdateAdvisor(c *fiber.Ctx) error {
	studentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Student ID format", "ID must be a valid UUID")
	}

	var req struct {
		AdvisorID uuid.UUID `json:"advisor_id" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if req.AdvisorID == uuid.Nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Advisor ID", "Advisor ID must be provided")
	}

	err = s.repo.UpdateAdvisor(studentID, req.AdvisorID)
	if err != nil {
		if errors.Is(err, errors.New("student not found or advisor already set")) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "Student or Advisor not found.")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update advisor", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Advisor updated successfully"})
}