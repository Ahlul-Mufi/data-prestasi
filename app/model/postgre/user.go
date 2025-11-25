package modelpostgre

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"password_hash"`
    FullName     string    `json:"full_name"`
    RoleID       *uuid.UUID `json:"role_id"`
    IsActive     bool      `json:"is_active"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
    Username     string     `json:"username"`
    Email        string     `json:"email"`
    PasswordHash string     `json:"password_hash"`
    FullName     string     `json:"full_name"`
    RoleID       *uuid.UUID `json:"role_id"`
}

type UpdateUserRequest struct {
    FullName *string     `json:"full_name"`
    RoleID   *uuid.UUID  `json:"role_id"`
    IsActive *bool       `json:"is_active"`
}
