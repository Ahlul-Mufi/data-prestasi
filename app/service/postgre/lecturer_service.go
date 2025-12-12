package servicepostgre

import (
	"errors"

	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	helper "github.com/Ahlul-Mufi/data-prestasi/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LecturerService interface {
	GetAll(c *fiber.Ctx) error
	GetAdvisees(c *fiber.Ctx) error
}

type lecturerService struct {
	repo repo.LecturerRepository
}

func NewLecturerService(r repo.LecturerRepository) LecturerService {
	return &lecturerService{r}
}

func (s *lecturerService) GetAll(c *fiber.Ctx) error {
	lecturers, err := s.repo.GetAll()
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch lecturers", err.Error())
	}
	return helper.SuccessResponse(c, fiber.StatusOK, lecturers)
}

func (s *lecturerService) GetAdvisees(c *fiber.Ctx) error {
	lecturerID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", "ID must be a valid UUID")
	}
	_, err = s.repo.GetByID(lecturerID)
	if err != nil {
		if errors.Is(err, errors.New("lecturer not found")) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "Lecturer not found")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check lecturer", err.Error())
	}

	advisees, err := s.repo.GetAdvisees(lecturerID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch advisees", err.Error())
	}
	return helper.SuccessResponse(c, fiber.StatusOK, advisees)
}