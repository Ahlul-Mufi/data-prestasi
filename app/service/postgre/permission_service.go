package servicepostgre

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
)

type PermissionService interface {
    GetAll(c *fiber.Ctx) error
    GetByID(c *fiber.Ctx) error
    Create(c *fiber.Ctx) error
    Update(c *fiber.Ctx) error
    Delete(c *fiber.Ctx) error
}

type permissionService struct {
    repo repo.PermissionRepository
}

func NewPermissionService(r repo.PermissionRepository) PermissionService {
    return &permissionService{r}
}

// @Summary Ambil semua izin
// @Description Mengambil daftar semua izin yang tersedia (Permission).
// @Tags Permissions
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} modelpostgre.Permission "Daftar izin"
// @Failure 500 {object} map[string]interface{} "Gagal mengambil izin"
// @Router /api/v1/permissions [get]
func (s *permissionService) GetAll(c *fiber.Ctx) error {
    data, err := s.repo.GetAll()
    if err != nil {
        return c.Status(500).JSON(err.Error())
    }
    return c.JSON(data)
}

// @Summary Ambil izin berdasarkan ID
// @Description Mengambil detail izin berdasarkan ID.
// @Tags Permissions
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID Izin (UUID)"
// @Success 200 {object} modelpostgre.Permission "Izin ditemukan"
// @Failure 400 {object} map[string]interface{} "Format ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Izin tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal mengambil izin"
// @Router /api/v1/permissions/{id} [get]
func (s *permissionService) GetByID(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return c.Status(400).JSON("invalid uuid")
    }

    data, err := s.repo.GetByID(id)
    if err != nil {
        return c.Status(404).JSON("not found")
    }

    return c.JSON(data)
}

// @Summary Buat izin baru
// @Description Membuat izin (Permission) baru.
// @Tags Permissions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param createPermissionRequest body m.CreatePermissionRequest true "Detail Izin Baru"
// @Success 201 {object} modelpostgre.Permission "Izin berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "Body request tidak valid"
// @Failure 500 {object} map[string]interface{} "Gagal membuat izin"
// @Router /api/v1/permissions [post]
func (s *permissionService) Create(c *fiber.Ctx) error {
    var req m.CreatePermissionRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON("invalid body")
    }

    p := m.Permission{
        Name:        req.Name,
        Resource:    req.Resource,
        Action:      req.Action,
        Description: req.Description,
    }

    newP, err := s.repo.Create(p)
    if err != nil {
        return c.Status(500).JSON(err.Error())
    }

    return c.Status(201).JSON(newP)
}

// @Summary Perbarui izin
// @Description Memperbarui detail izin.
// @Tags Permissions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "ID Izin (UUID)"
// @Param updatePermissionRequest body m.UpdatePermissionRequest true "Detail Izin yang diperbarui"
// @Success 200 {object} modelpostgre.Permission "Izin berhasil diperbarui"
// @Failure 400 {object} map[string]interface{} "Format ID tidak valid / Body request tidak valid"
// @Failure 404 {object} map[string]interface{} "Izin tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal memperbarui izin"
// @Router /api/v1/permissions/{id} [put]
func (s *permissionService) Update(c *fiber.Ctx) error {
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON("invalid uuid")
    }

    var req m.UpdatePermissionRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON("invalid body")
    }

    p := m.Permission{
        Name:        req.Name,
        Resource:    req.Resource,
        Action:      req.Action,
        Description: req.Description,
    }

    updated, err := s.repo.Update(id, p)
    if err != nil {
        return c.Status(500).JSON(err.Error())
    }

    return c.JSON(updated)
}

// @Summary Hapus izin
// @Description Menghapus izin berdasarkan ID.
// @Tags Permissions
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID Izin (UUID)"
// @Success 200 {object} modelpostgre.Permission "Izin berhasil dihapus"
// @Failure 400 {object} map[string]interface{} "Format ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Izin tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal menghapus izin"
// @Router /api/v1/permissions/{id} [delete]
func (s *permissionService) Delete(c *fiber.Ctx) error {
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON("invalid uuid")
    }

    if err := s.repo.Delete(id); err != nil {
        return c.Status(404).JSON("not found")
    }

    return c.JSON(fiber.Map{"message": "deleted"})
}
