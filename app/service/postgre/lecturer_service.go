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

// @Summary Ambil semua Dosen
// @Description Mengambil daftar semua profil dosen.
// @Tags Lecturers
// @Security BearerAuth
// @Produce json
// @Success 200 {array} modelpostgre.Lecturer "Daftar dosen"
// @Failure 500 {object} map[string]interface{} "Gagal mengambil data dosen"
// @Router /api/v1/lecturers [get]
func (s *lecturerService) GetAll(c *fiber.Ctx) error {
	lecturers, err := s.repo.GetAll()
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch lecturers", err.Error())
	}
	return helper.SuccessResponse(c, fiber.StatusOK, lecturers)
}

// @Summary Ambil Dosen berdasarkan ID
// @Description Mengambil detail profil dosen berdasarkan ID (UUID).
// @Tags Lecturers
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID Pengguna Dosen (UUID)"
// @Success 200 {object} modelpostgre.Lecturer "Dosen ditemukan"
// @Failure 400 {object} map[string]interface{} "Format ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Dosen tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal mengambil data dosen"
// @Router /api/v1/lecturers/{id} [get]
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

// @Summary Ambil Mahasiswa Bimbingan Dosen
// @Description Mengambil daftar mahasiswa yang dibimbing oleh dosen tertentu.
// @Tags Lecturers
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID Pengguna Dosen (UUID)"
// @Success 200 {array} modelpostgre.Student "Daftar mahasiswa bimbingan"
// @Failure 400 {object} map[string]interface{} "Format ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Dosen tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal mengambil mahasiswa bimbingan"
// @Router /api/v1/lecturers/{id}/advisees [get]
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

// @Summary Buat profil Dosen baru
// @Description Membuat profil dosen baru (Membutuhkan User ID yang sudah ada).
// @Tags Lecturers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param createLecturerRequest body m.CreateLecturerRequest true "Detail Dosen Baru"
// @Success 201 {object} modelpostgre.Lecturer "Dosen berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "Body request tidak valid / User ID tidak ditemukan"
// @Failure 409 {object} map[string]interface{} "Dosen sudah ada (Duplicate Entry)"
// @Failure 500 {object} map[string]interface{} "Gagal membuat profil dosen"
// @Router /api/v1/lecturers [post]
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

// @Summary Perbarui profil Dosen
// @Description Memperbarui detail profil dosen (NIDN, Departemen).
// @Tags Lecturers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID Pengguna Dosen (UUID)"
// @Param updateLecturerRequest body m.UpdateLecturerRequest true "Detail Dosen yang diperbarui"
// @Success 200 {object} modelpostgre.Lecturer "Profil dosen berhasil diperbarui"
// @Failure 400 {object} map[string]interface{} "Format ID tidak valid / Body request tidak valid"
// @Failure 404 {object} map[string]interface{} "Dosen tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal memperbarui profil dosen"
// @Router /api/v1/lecturers/{id} [put]
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

// @Summary Hapus profil Dosen
// @Description Menghapus profil dosen berdasarkan ID (UUID).
// @Tags Lecturers
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID Pengguna Dosen (UUID)"
// @Success 200 {object} modelpostgre.Lecturer "Profil dosen berhasil dihapus"
// @Failure 400 {object} map[string]interface{} "Format ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Dosen tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal menghapus profil dosen"
// @Router /api/v1/lecturers/{id} [delete]
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
