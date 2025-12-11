package servicepostgre

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	helper "github.com/Ahlul-Mufi/data-prestasi/helper"
	utils "github.com/Ahlul-Mufi/data-prestasi/utils/postgre"
)

type UserService interface {
    Login(c *fiber.Ctx) error
    Profile(c *fiber.Ctx) error
    Refresh(c *fiber.Ctx) error     
    Logout(c *fiber.Ctx) error  
    GetUsers(c *fiber.Ctx) error        
    GetUserByID(c *fiber.Ctx) error     
    CreateUser(c *fiber.Ctx) error
    UpdateUser(c *fiber.Ctx) error
    DeleteUser(c *fiber.Ctx) error
    UpdateUserRole(c *fiber.Ctx) error  
}

type userService struct {
    repo repo.UserRepository
}

func NewUserService(r repo.UserRepository) UserService {
    return &userService{r}
}

func (s *userService) Login(c *fiber.Ctx) error {
    var req m.LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
    }

    user, err := s.repo.FindByIdentity(req.Identity)
    if err != nil {
        if errors.Is(err, repo.ErrUserNotFound) || err == sql.ErrNoRows {
            return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid identity or password", "")
        }
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Internal Server Error", err.Error())
    }
    if !user.IsActive {
        return helper.ErrorResponse(c, fiber.StatusUnauthorized, "User is inactive", "")
    }

    if !utils.CheckPassword(user.PasswordHash, req.Password) {
        return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid identity or password", "")
    }
    
    accessToken, err := utils.GenerateToken(user) 
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Access Token generation failed", err.Error())
    }
    
    refreshToken, err := utils.GenerateRefreshToken(user)
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Refresh Token generation failed", err.Error())
    }

    user.PasswordHash = ""

    return helper.SuccessResponse(c, fiber.StatusOK, m.LoginResponse{
        Token: accessToken,
        RefreshToken: refreshToken,
        User:  user,
    })
}

func (s *userService) Refresh(c *fiber.Ctx) error {
    var req struct {
        RefreshToken string `json:"refresh_token"`
    }
    if err := c.BodyParser(&req); err != nil {
        return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
    }

    claims, err := utils.ValidateRefreshToken(req.RefreshToken)
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired refresh token", err.Error())
    }

    user, err := s.repo.FindByID(claims.UserID)
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusUnauthorized, "User not found", "")
    }

    newAccessToken, err := utils.GenerateToken(user)
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Token generation failed", err.Error())
    }

    return helper.SuccessResponse(c, fiber.StatusOK, fiber.Map{
        "token": newAccessToken,
    })
}

func (s *userService) Logout(c *fiber.Ctx) error { 
    return helper.SuccessResponse(c, fiber.StatusOK, fiber.Map{
        "message": "Logged out successfully. Please delete your tokens. The current session will expire in 2 minutes.",
    })
}


func (s *userService) Profile(c *fiber.Ctx) error {
    userIDStr := c.Locals("user_id")

    if userIDStr == nil {
        return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", "User ID not found in context")
    }

    userID, err := uuid.Parse(userIDStr.(string))
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID format", err.Error())
    }

    user, err := s.repo.FindByID(userID) 
    if err != nil {
        if errors.Is(err, repo.ErrUserNotFound) || err == sql.ErrNoRows {
            return helper.ErrorResponse(c, fiber.StatusNotFound, "User not found", "")
        }
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Internal Server Error", err.Error())
    }
    
    user.PasswordHash = ""
    return helper.SuccessResponse(c, fiber.StatusOK, user)
}

func (s *userService) GetUsers(c *fiber.Ctx) error {
    users, err := s.repo.FindAll()
    if err != nil {
        log.Println("Database error in GetUsers:", err)
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve users", err.Error())
    }
    
    for i := range users {
        users[i].PasswordHash = ""
    }

    return helper.SuccessResponse(c, fiber.StatusOK, users)
}

func (s *userService) GetUserByID(c *fiber.Ctx) error {
    idStr := c.Params("id")
    userID, err := uuid.Parse(idStr)
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid User ID format", err.Error())
    }

    user, err := s.repo.FindByID(userID)
    if err != nil {
        if errors.Is(err, repo.ErrUserNotFound) {
            return helper.ErrorResponse(c, fiber.StatusNotFound, "User not found", "")
        }
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Internal Server Error", err.Error())
    }

    user.PasswordHash = ""
    return helper.SuccessResponse(c, fiber.StatusOK, user)
}

func (s *userService) CreateUser(c *fiber.Ctx) error {
    var req m.CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
    }
    
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash password", err.Error())
    }

    newUser := m.User{
        ID:           uuid.New(),
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: hashedPassword,
        FullName:     req.FullName,
        IsActive:     req.IsActive,
    }

    createdUser, err := s.repo.Create(newUser, req.RoleName)
    if err != nil {
        if errors.Is(err, repo.ErrRoleNotFound) {
            return helper.ErrorResponse(c, fiber.StatusBadRequest, "Role not found", "")
        }
        if errors.Is(err, repo.ErrDuplicateUsernameOrEmail) {
            return helper.ErrorResponse(c, fiber.StatusConflict, "Username or Email already exists", "")
        }
        log.Println("Database error in CreateUser:", err)
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user", err.Error())
    }

    createdUser.PasswordHash = "" 
    return helper.SuccessResponse(c, fiber.StatusCreated, createdUser)
}

func (s *userService) UpdateUser(c *fiber.Ctx) error {
    idStr := c.Params("id")
    userID, err := uuid.Parse(idStr)
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid User ID format", err.Error())
    }

    var req m.UpdateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
    }

    existingUser, err := s.repo.FindByID(userID)
    if err != nil {
        if errors.Is(err, repo.ErrUserNotFound) {
            return helper.ErrorResponse(c, fiber.StatusNotFound, "User not found", "")
        }
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Internal Server Error", err.Error())
    }
    
    if req.Username != "" { existingUser.Username = req.Username }
    if req.Email != "" { existingUser.Email = req.Email }
    if req.FullName != "" { existingUser.FullName = req.FullName }
    
    if req.Password != nil && *req.Password != "" {
        hashedPassword, err := utils.HashPassword(*req.Password)
        if err != nil {
             return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash new password", err.Error())
        }
        existingUser.PasswordHash = hashedPassword
    }
    
    if req.IsActive != nil {
        existingUser.IsActive = *req.IsActive
    }
    
    updatedUser, err := s.repo.Update(existingUser, req.RoleName)
    if err != nil {
        if errors.Is(err, repo.ErrUserNotFound) {
             return helper.ErrorResponse(c, fiber.StatusNotFound, "User not found", "")
        }
        if errors.Is(err, repo.ErrRoleNotFound) {
            return helper.ErrorResponse(c, fiber.StatusBadRequest, "Role not found", "")
        }
        log.Println("Database error in UpdateUser:", err)
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user", err.Error())
    }

    updatedUser.PasswordHash = "" 
    return helper.SuccessResponse(c, fiber.StatusOK, updatedUser)
}

func (s *userService) DeleteUser(c *fiber.Ctx) error {
    idStr := c.Params("id")
    userID, err := uuid.Parse(idStr)
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid User ID format", err.Error())
    }

    err = s.repo.DeleteUser(userID)
    if err != nil {
        if errors.Is(err, repo.ErrUserNotFound) {
            return helper.ErrorResponse(c, fiber.StatusNotFound, "User not found", "")
        }
        log.Println("Database error in DeleteUser:", err)
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete user", err.Error())
    }

    return helper.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "User deleted successfully"})
}

func (s *userService) UpdateUserRole(c *fiber.Ctx) error {
    idStr := c.Params("id")
    userID, err := uuid.Parse(idStr)
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid User ID format", err.Error())
    }

    var req m.UpdateUserRoleRequest
    if err := c.BodyParser(&req); err != nil {
        return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
    }
    
    roleID, err := s.repo.FindRoleIDByName(req.RoleName)
    if err != nil {
        if errors.Is(err, repo.ErrRoleNotFound) {
            return helper.ErrorResponse(c, fiber.StatusBadRequest, "Role not found", "")
        }
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve role", err.Error())
    }
    
    err = s.repo.UpdateRole(userID, roleID)
    if err != nil {
        if errors.Is(err, repo.ErrUserNotFound) {
            return helper.ErrorResponse(c, fiber.StatusNotFound, "User not found", "")
        }
        log.Println("Database error in UpdateUserRole:", err)
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user role", err.Error())
    }

    updatedUser, _ := s.repo.FindByID(userID)
    updatedUser.PasswordHash = ""
    
    return helper.SuccessResponse(c, fiber.StatusOK, updatedUser)
}