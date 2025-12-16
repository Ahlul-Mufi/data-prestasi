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

// @Summary Ambil semua Mahasiswa
// @Description Mengambil daftar semua profil mahasiswa.
// @Tags Students
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} modelpostgre.Student "Daftar mahasiswa"
// @Failure 500 {object} map[string]interface{} "Gagal mengambil data mahasiswa"
// @Router /api/v1/students [get]
func (s *studentService) GetAll(c *fiber.Ctx) error {
	students, err := s.repo.GetAll()
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch students", err.Error())
	}
	return helper.SuccessResponse(c, fiber.StatusOK, students)
}

// @Summary Ambil Mahasiswa berdasarkan ID
// @Description Mengambil detail profil mahasiswa berdasarkan ID (UUID).
// @Tags Students
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID Pengguna Mahasiswa (UUID)"
// @Success 200 {object} modelpostgre.Student "Mahasiswa ditemukan"
// @Failure 400 {object} map[string]interface{} "Format ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Mahasiswa tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal mengambil data mahasiswa"
// @Router /api/v1/students/{id} [get]
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

// @Summary Ambil Prestasi Mahasiswa
// @Description Mengambil daftar semua prestasi (Achievement References) yang dimiliki oleh mahasiswa tertentu.
// @Tags Students
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID Pengguna Mahasiswa (UUID)"
// @Success 200 {array} modelpostgre.AchievementReference "Daftar prestasi mahasiswa"
// @Failure 400 {object} map[string]interface{} "Format ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Mahasiswa tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal mengambil prestasi"
// @Router /api/v1/students/{id}/achievements [get]
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

// @Summary Buat profil Mahasiswa baru
// @Description Membuat profil mahasiswa baru (Membutuhkan User ID yang sudah ada).
// @Tags Students
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param createStudentRequest body m.CreateStudentRequest true "Detail Mahasiswa Baru"
// @Success 201 {object} modelpostgre.Student "Mahasiswa berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "Body request tidak valid / User ID tidak ditemukan"
// @Failure 409 {object} map[string]interface{} "Mahasiswa sudah ada (Duplicate Entry)"
// @Failure 500 {object} map[string]interface{} "Gagal membuat profil mahasiswa"
// @Router /api/v1/students [post]
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

// @Summary Perbarui profil Mahasiswa
// @Description Memperbarui detail profil mahasiswa (NIM, Program Studi, Tahun Akademik).
// @Tags Students
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "ID Pengguna Mahasiswa (UUID)"
// @Param updateStudentRequest body m.UpdateStudentRequest true "Detail Mahasiswa yang diperbarui"
// @Success 200 {object} modelpostgre.Student "Profil mahasiswa berhasil diperbarui"
// @Failure 400 {object} map[string]interface{} "Format ID tidak valid / Body request tidak valid"
// @Failure 404 {object} map[string]interface{} "Mahasiswa tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal memperbarui profil mahasiswa"
// @Router /api/v1/students/{id} [put]
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

// @Summary Hapus profil Mahasiswa
// @Description Menghapus profil mahasiswa berdasarkan ID (UUID).
// @Tags Students
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID Pengguna Mahasiswa (UUID)"
// @Success 200 {object} modelpostgre.Student "Profil mahasiswa berhasil dihapus"
// @Failure 400 {object} map[string]interface{} "Format ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Mahasiswa tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal menghapus profil mahasiswa"
// @Router /api/v1/students/{id} [delete]
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

// @Summary Perbarui Dosen Pembimbing Mahasiswa
// @Description Memperbarui Dosen Pembimbing (Advisor) untuk mahasiswa tertentu.
// @Tags Students
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "ID Pengguna Mahasiswa (UUID)"
// @Param updateAdvisorRequest body m.UpdateAdvisorRequest true "ID Dosen Pembimbing (Advisor ID)"
// @Success 200 {object} modelpostgre.Student "Dosen Pembimbing berhasil diperbarui"
// @Failure 400 {object} map[string]interface{} "Format ID tidak valid / Body request tidak valid"
// @Failure 404 {object} map[string]interface{} "Mahasiswa atau Dosen Pembimbing tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal memperbarui Dosen Pembimbing"
// @Router /api/v1/students/{id}/advisor [patch]
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
