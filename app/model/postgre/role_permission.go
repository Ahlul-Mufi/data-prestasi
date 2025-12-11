package modelpostgre

import "github.com/google/uuid"

type RolePermission struct {
    RoleID       uuid.UUID `json:"role_id"`
    PermissionID uuid.UUID `json:"permission_id"`
}

type AddRolePermissionRequest struct {
    RoleID       uuid.UUID `json:"role_id" validate:"required"`
    PermissionID uuid.UUID `json:"permission_id" validate:"required"`
}