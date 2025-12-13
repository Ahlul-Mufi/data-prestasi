package servicepostgre

import (
	"errors"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	helper "github.com/Ahlul-Mufi/data-prestasi/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StudentService interface {
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetAchievements(c *fiber.Ctx) error
	CreateStudent(c *fiber.Ctx) error
	UpdateStudent(c *fiber.Ctx) error
	DeleteStudent(c *fiber.Ctx) error
	UpdateAdvisor(c *fiber.Ctx) error
}

type studentService struct {
	repo     repo.StudentRepository
	arRepo   repo.AchievementReferenceRepository
	userRepo repo.UserRepository
	lecRepo  repo.LecturerRepository
}

func NewStudentService(r repo.StudentRepository, arRepo repo.AchievementReferenceRepository, userRepo repo.UserRepository, lecRepo repo.LecturerRepository) StudentService {
	return &studentService{r, arRepo, userRepo, lecRepo}
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

func (s *studentService) CreateStudent(c *fiber.Ctx) error {
	var req m.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	_, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "User not found", "")
	}

	if req.AdvisorID != nil {
		_, err := s.lecRepo.GetByID(*req.AdvisorID)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Advisor (Lecturer) not found", "")
		}
	}

	newStudent := m.Student{
		ID:           uuid.New(),
		UserID:       req.UserID,
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
		AdvisorID:    req.AdvisorID,
	}

	createdStudent, err := s.repo.Create(newStudent)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create student profile", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, createdStudent)
}

func (s *studentService) UpdateStudent(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", "ID must be a valid UUID")
	}

	var req m.UpdateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	existingStudent, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, errors.New("student not found")) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Student not found", "")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch student", err.Error())
	}

	if req.StudentID != "" {
		existingStudent.StudentID = req.StudentID
	}
	if req.ProgramStudy != "" {
		existingStudent.ProgramStudy = req.ProgramStudy
	}
	if req.AcademicYear != "" {
		existingStudent.AcademicYear = req.AcademicYear
	}
	if req.AdvisorID != nil {
		_, err := s.lecRepo.GetByID(*req.AdvisorID)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Advisor (Lecturer) not found", "")
		}
		existingStudent.AdvisorID = req.AdvisorID
	}

	updatedStudent, err := s.repo.Update(existingStudent)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update student profile", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, updatedStudent)
}

func (s *studentService) DeleteStudent(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", "ID must be a valid UUID")
	}

	err = s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, errors.New("student not found")) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Student not found", "")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete student profile", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Student profile deleted successfully"})
}

func (s *studentService) UpdateAdvisor(c *fiber.Ctx) error {
	studentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Student ID format", "ID must be a valid UUID")
	}

	var req m.UpdateAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if req.AdvisorID == uuid.Nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Advisor ID", "Advisor ID must be provided")
	}

	_, err = s.repo.GetByID(studentID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Student not found", "")
	}

	_, err = s.lecRepo.GetByID(req.AdvisorID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Advisor (Lecturer) not found", "")
	}

	err = s.repo.UpdateAdvisor(studentID, req.AdvisorID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update advisor", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Advisor updated successfully"})
}
