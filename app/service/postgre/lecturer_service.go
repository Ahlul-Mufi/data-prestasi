package servicepostgre

import (
	"errors"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	helper "github.com/Ahlul-Mufi/data-prestasi/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LecturerService interface {
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetAdvisees(c *fiber.Ctx) error
	CreateLecturer(c *fiber.Ctx) error
	UpdateLecturer(c *fiber.Ctx) error
	DeleteLecturer(c *fiber.Ctx) error
}

type lecturerService struct {
	repo     repo.LecturerRepository
	userRepo repo.UserRepository
}

func NewLecturerService(r repo.LecturerRepository, userRepo repo.UserRepository) LecturerService {
	return &lecturerService{r, userRepo}
}

func (s *lecturerService) GetAll(c *fiber.Ctx) error {
	lecturers, err := s.repo.GetAll()
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch lecturers", err.Error())
	}
	return helper.SuccessResponse(c, fiber.StatusOK, lecturers)
}

func (s *lecturerService) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", "ID must be a valid UUID")
	}

	lecturer, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, errors.New("lecturer not found")) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Lecturer not found", "")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch lecturer", err.Error())
	}
	return helper.SuccessResponse(c, fiber.StatusOK, lecturer)
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

func (s *lecturerService) CreateLecturer(c *fiber.Ctx) error {
	var req m.CreateLecturerRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	_, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "User not found", "")
	}

	newLecturer := m.Lecturer{
		ID:         uuid.New(),
		UserID:     req.UserID,
		LecturerID: req.LecturerID,
		Department: req.Department,
	}

	createdLecturer, err := s.repo.Create(newLecturer)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create lecturer profile", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, createdLecturer)
}

func (s *lecturerService) UpdateLecturer(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", "ID must be a valid UUID")
	}

	var req m.UpdateLecturerRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	existingLecturer, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, errors.New("lecturer not found")) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Lecturer not found", "")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch lecturer", err.Error())
	}

	if req.LecturerID != "" {
		existingLecturer.LecturerID = req.LecturerID
	}
	if req.Department != "" {
		existingLecturer.Department = req.Department
	}

	updatedLecturer, err := s.repo.Update(existingLecturer)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update lecturer profile", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, updatedLecturer)
}

func (s *lecturerService) DeleteLecturer(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", "ID must be a valid UUID")
	}

	err = s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, errors.New("lecturer not found")) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Lecturer not found", "")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete lecturer profile", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Lecturer profile deleted successfully"})
}
