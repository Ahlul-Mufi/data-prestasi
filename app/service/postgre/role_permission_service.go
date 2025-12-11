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