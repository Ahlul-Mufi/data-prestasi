package servicepostgre

import (
	"errors"
	"strings"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	helper "github.com/Ahlul-Mufi/data-prestasi/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)
type RolePermissionService interface {
	Add(c *fiber.Ctx) error
	Remove(c *fiber.Ctx) error
}

type rolePermissionService struct {
	repo repo.RolePermissionRepository
}

func NewRolePermissionService(r repo.RolePermissionRepository) RolePermissionService {
	return &rolePermissionService{r}
}

// @Summary Tambahkan Izin ke Peran
// @Description Menghubungkan izin tertentu dengan peran tertentu.
// @Tags Roles & Permissions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param addRolePermissionRequest body m.AddRolePermissionRequest true "Role ID dan Permission ID"
// @Success 201 {object} modelpostgre.RolePermission "Izin berhasil ditambahkan ke peran"
// @Failure 400 {object} map[string]interface{} "Body request tidak valid / ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Role atau Permission tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal menambahkan relasi"
// @Router /api/v1/role-permissions [post]
func (s *rolePermissionService) Add(c *fiber.Ctx) error {
	var req m.AddRolePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if req.RoleID == uuid.Nil || req.PermissionID == uuid.Nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid IDs", "RoleID and PermissionID must be provided.")
	}

    rp := m.RolePermission(req)

	result, err := s.repo.Add(rp)
	if err != nil {
        if strings.Contains(err.Error(), "role or permission not found") {
            return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "One or both IDs (Role/Permission) do not exist.")
        }
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add role permission", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, result)
}

// @Summary Hapus Izin dari Peran
// @Description Memutus hubungan izin dari peran tertentu.
// @Tags Roles & Permissions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param removeRolePermissionRequest body m.RolePermission true "Role ID dan Permission ID"
// @Success 200 {object} modelpostgre.RolePermission "Izin berhasil dihapus dari peran"
// @Failure 400 {object} map[string]interface{} "Body request tidak valid / ID tidak valid"
// @Failure 404 {object} map[string]interface{} "Relasi Role-Permission tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Gagal menghapus relasi"
// @Router /api/v1/role-permissions [delete]
func (s *rolePermissionService) Remove(c *fiber.Ctx) error {
	var req m.RolePermission
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	
	if req.RoleID == uuid.Nil || req.PermissionID == uuid.Nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid IDs", "RoleID and PermissionID must be provided.")
	}

	err := s.repo.Remove(req.RoleID, req.PermissionID)

	if err != nil {
		if errors.Is(err, errors.New("role permission not found")) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "The specified role permission link does not exist.")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to remove role permission", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Role permission removed successfully"})
}