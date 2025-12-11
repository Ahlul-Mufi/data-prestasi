package modelpostgre

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID  `json:"id"`
    Username     string     `json:"username"`
    Email        string     `json:"email"`
    PasswordHash string     `json:"password_hash,omitempty"` 
    FullName     string     `json:"full_name"`
    RoleID       *uuid.UUID `json:"role_id"`
    IsActive     bool       `json:"is_active"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}

type LoginRequest struct {
    Identity string `json:"identity"` 
    Password string `json:"password"`
}

type LoginResponse struct {
    Token        string `json:"token"`
    RefreshToken string `json:"refresh_token"`
    User         User   `json:"user"`
}

type JWTClaims struct {
    UserID   uuid.UUID  `json:"user_id"` 
    Username string     `json:"username"`
    RoleID   *uuid.UUID `json:"role_id"` 
    jwt.RegisteredClaims
}

type RefreshClaims struct {
    UserID           uuid.UUID `json:"user_id"`
    jwt.RegisteredClaims
}

type CreateUserRequest struct {
    Username string `json:"username"`
    Email string `json:"email"`
    Password string `json:"password"`
    FullName string `json:"full_name"`
    RoleName string `json:"role_name"` 
    IsActive bool `json:"is_active"`
}
type UpdateUserRequest struct {
    Username string `json:"username"`
    Email string `json:"email"`
    Password *string `json:"password"` 
    FullName string `json:"full_name"`
    RoleName string `json:"role_name"`
    IsActive *bool `json:"is_active"` 
}

type UpdateUserRoleRequest struct {
    RoleName string `json:"role_name"`
}