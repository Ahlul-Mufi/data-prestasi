package servicepostgre

import (
	"database/sql"
	"errors"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	helper "github.com/Ahlul-Mufi/data-prestasi/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RolePermissionService interface {
	AddPermissionToRole(c *fiber.Ctx) error
	RemovePermissionFromRole(c *fiber.Ctx) error
	GetPermissionsByRoleID(c *fiber.Ctx) error
}

type rolePermissionService struct {
	repo repo.RolePermissionRepository
}

func NewRolePermissionService(r repo.RolePermissionRepository) RolePermissionService {
	return &rolePermissionService{r}
}

func (s *rolePermissionService) AddPermissionToRole(c *fiber.Ctx) error {
	var req m.AddRolePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	rp := m.RolePermission{
		RoleID:       req.RoleID,
		PermissionID: req.PermissionID,
	}

	newRP, err := s.repo.AddPermissionToRole(rp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, errors.New("permission already linked to role")) {
			return helper.ErrorResponse(c, fiber.StatusConflict, "Permission already linked to role", err.Error())
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add permission to role", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, newRP)
}

func (s *rolePermissionService) RemovePermissionFromRole(c *fiber.Ctx) error {
	roleIDStr := c.Params("roleId")
	permissionIDStr := c.Params("permissionId")
	
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role ID format", err.Error())
	}
	
	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid permission ID format", err.Error())
	}

	err = s.repo.RemovePermissionFromRole(roleID, permissionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Role-Permission relation not found", err.Error())
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to remove permission from role", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Permission removed from role successfully",
	})
}

func (s *rolePermissionService) GetPermissionsByRoleID(c *fiber.Ctx) error {
	roleIDStr := c.Params("roleId")
	
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid role ID format", err.Error())
	}
	
	perms, err := s.repo.GetPermissionsByRoleID(roleID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch permissions", err.Error())
	}
	
	return helper.SuccessResponse(c, fiber.StatusOK, perms)
}