package servicepostgre

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
)

type RoleService interface {
    GetAll(c *fiber.Ctx) error
    GetByID(c *fiber.Ctx) error
    Create(c *fiber.Ctx) error
    Update(c *fiber.Ctx) error
    Delete(c *fiber.Ctx) error
}

type roleService struct {
    repo repo.RoleRepository
}

func NewRoleService(r repo.RoleRepository) RoleService {
    return &roleService{r}
}

// @Summary Ambil semua peran
// @Description Mengambil daftar semua peran yang tersedia.
// @Tags Roles
// @Security BearerAuth
// @Produce json
// @Success 200 {array} modelpostgre.Role "Daftar peran"
// @Failure 500 {object} map[string]interface{} "Gagal mengambil peran"
// @Router /api/v1/roles [get]
func (s *roleService) GetAll(c *fiber.Ctx) error {
    data, err := s.repo.GetAll()
    if err != nil {
        return c.Status(500).JSON(err.Error())
    }
    return c.JSON(data)
}

// @Summary Ambil peran berdasarkan ID
// @Description Mengambil detail peran berdasarkan ID.
// @Tags Roles
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID Peran (UUID)"
// @Success 200 {object} modelpostgre.Role "Peran ditemukan"
// @Failure 400 {object} map[string]interface{} "Format ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Peran tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal mengambil peran"
// @Router /api/v1/roles/{id} [get]
func (s *roleService) GetByID(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid uuid"})
    }

    data, err := s.repo.GetByID(id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "not found"})
    }

    return c.JSON(data)
}

// @Summary Buat peran baru
// @Description Membuat peran baru.
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param createRoleRequest body m.CreateRoleRequest true "Detail Peran Baru"
// @Success 201 {object} modelpostgre.Role "Peran berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "Body request tidak valid"
// @Failure 500 {object} map[string]interface{} "Gagal membuat peran"
// @Router /api/v1/roles [post]
func (s *roleService) Create(c *fiber.Ctx) error {
    var req m.CreateRoleRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON("invalid body")
    }

    role := m.Role{
        ID:          uuid.New(),
        Name:        req.Name,
        Description: req.Description,
    }

    newRole, err := s.repo.Create(role)
    if err != nil {
        return c.Status(500).JSON(err.Error())
    }

    return c.Status(201).JSON(newRole)
}

// @Summary Perbarui peran
// @Description Memperbarui nama atau deskripsi peran.
// @Tags Roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID Peran (UUID)"
// @Param updateRoleRequest body m.UpdateRoleRequest true "Detail Peran yang diperbarui"
// @Success 200 {object} modelpostgre.Role "Peran berhasil diperbarui"
// @Failure 400 {object} map[string]interface{} "Format ID tidak valid / Body request tidak valid"
// @Failure 404 {object} map[string]interface{} "Peran tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal memperbarui peran"
// @Router /api/v1/roles/{id} [put]
func (s *roleService) Update(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return c.Status(400).JSON("invalid uuid")
    }

    var req m.UpdateRoleRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON("invalid body")
    }

    update := m.Role{
        Name:        req.Name,
        Description: req.Description,
    }

    updated, err := s.repo.Update(id, update)
    if err != nil {
        return c.Status(500).JSON(err.Error())
    }

    return c.JSON(updated)
}

// @Summary Hapus peran
// @Description Menghapus peran berdasarkan ID.
// @Tags Roles
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID Peran (UUID)"
// @Success 200 {object} modelpostgre.Role "Peran berhasil dihapus"
// @Failure 400 {object} map[string]interface{} "Format ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Peran tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal menghapus peran"
// @Router /api/v1/roles/{id} [delete]
func (s *roleService) Delete(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return c.Status(400).JSON("invalid uuid")
    }

    if err := s.repo.Delete(id); err != nil {
        return c.Status(404).JSON("not found")
    }

    return c.JSON(fiber.Map{"message": "deleted"})
}