package servicepostgre

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	modelpostgre "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	utils "github.com/Ahlul-Mufi/data-prestasi/utils/postgre"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByIdentity(identity string) (modelpostgre.User, error) {
	args := m.Called(identity)
	return args.Get(0).(modelpostgre.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (modelpostgre.User, error) {
	args := m.Called(id)
	return args.Get(0).(modelpostgre.User), args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]modelpostgre.User, error) {
	args := m.Called()
	return args.Get(0).([]modelpostgre.User), args.Error(1)
}

func (m *MockUserRepository) Create(user modelpostgre.User, roleName string) (modelpostgre.User, error) {
	args := m.Called(user, roleName)
	return args.Get(0).(modelpostgre.User), args.Error(1)
}

func (m *MockUserRepository) Update(user modelpostgre.User, roleName string) (modelpostgre.User, error) {
	args := m.Called(user, roleName)
	return args.Get(0).(modelpostgre.User), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateRole(userID uuid.UUID, roleID uuid.UUID) error {
	args := m.Called(userID, roleID)
	return args.Error(0)
}

func (m *MockUserRepository) FindRoleIDByName(roleName string) (uuid.UUID, error) {
	args := m.Called(roleName)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockUserRepository) GetPermissions(roleID uuid.UUID) ([]string, error) {
	args := m.Called(roleID)
	return args.Get(0).([]string), args.Error(1)
}

// Helper function to create a test user
func createTestUser() modelpostgre.User {
	roleID := uuid.New()
	return modelpostgre.User{
		ID:           uuid.New(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		FullName:     "Test User",
		RoleID:       &roleID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func TestLogin_Success(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser()
	hashedPassword, _ := utils.HashPassword("password123")
	user.PasswordHash = hashedPassword

	mockRepo.On("FindByIdentity", "testuser").Return(user, nil)

	loginReq := modelpostgre.LoginRequest{
		Identity: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Post("/login", service.Login)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidRequestBody(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	app.Post("/login", service.Login)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestLogin_UserNotFound(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByIdentity", "nonexistent").Return(modelpostgre.User{}, repo.ErrUserNotFound)

	loginReq := modelpostgre.LoginRequest{
		Identity: "nonexistent",
		Password: "password123",
	}
	body, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Post("/login", service.Login)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InactiveUser(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser()
	user.IsActive = false

	mockRepo.On("FindByIdentity", "testuser").Return(user, nil)

	loginReq := modelpostgre.LoginRequest{
		Identity: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Post("/login", service.Login)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser()
	hashedPassword, _ := utils.HashPassword("correctpassword")
	user.PasswordHash = hashedPassword

	mockRepo.On("FindByIdentity", "testuser").Return(user, nil)

	loginReq := modelpostgre.LoginRequest{
		Identity: "testuser",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Post("/login", service.Login)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestRefresh_Success(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser()
	refreshToken, _ := utils.GenerateRefreshToken(user)

	mockRepo.On("FindByID", user.ID).Return(user, nil)

	refreshReq := struct {
		RefreshToken string `json:"refresh_token"`
	}{
		RefreshToken: refreshToken,
	}
	body, _ := json.Marshal(refreshReq)

	req := httptest.NewRequest("POST", "/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Post("/refresh", service.Refresh)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestRefresh_InvalidToken(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	refreshReq := struct {
		RefreshToken string `json:"refresh_token"`
	}{
		RefreshToken: "invalid.token.here",
	}
	body, _ := json.Marshal(refreshReq)

	req := httptest.NewRequest("POST", "/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Post("/refresh", service.Refresh)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestLogout_Success(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	req := httptest.NewRequest("POST", "/logout", nil)

	app.Post("/logout", service.Logout)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestProfile_Success(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser()
	mockRepo.On("FindByID", user.ID).Return(user, nil)

	req := httptest.NewRequest("GET", "/profile", nil)

	app.Get("/profile", func(c *fiber.Ctx) error {
		c.Locals("user_id", user.ID.String())
		return service.Profile(c)
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestProfile_NoUserID(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	req := httptest.NewRequest("GET", "/profile", nil)

	app.Get("/profile", service.Profile)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestProfile_InvalidUserIDFormat(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	req := httptest.NewRequest("GET", "/profile", nil)

	app.Get("/profile", func(c *fiber.Ctx) error {
		c.Locals("user_id", "invalid-uuid")
		return service.Profile(c)
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestProfile_UserNotFound(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	mockRepo.On("FindByID", userID).Return(modelpostgre.User{}, repo.ErrUserNotFound)

	req := httptest.NewRequest("GET", "/profile", nil)

	app.Get("/profile", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return service.Profile(c)
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestGetUsers_Success(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	users := []modelpostgre.User{createTestUser(), createTestUser()}
	mockRepo.On("FindAll").Return(users, nil)

	req := httptest.NewRequest("GET", "/users", nil)

	app.Get("/users", service.GetUsers)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestGetUsers_DatabaseError(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("FindAll").Return([]modelpostgre.User{}, errors.New("database error"))

	req := httptest.NewRequest("GET", "/users", nil)

	app.Get("/users", service.GetUsers)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser()
	mockRepo.On("FindByID", user.ID).Return(user, nil)

	req := httptest.NewRequest("GET", "/users/"+user.ID.String(), nil)

	app.Get("/users/:id", service.GetUserByID)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_InvalidID(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	req := httptest.NewRequest("GET", "/users/invalid-id", nil)

	app.Get("/users/:id", service.GetUserByID)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetUserByID_UserNotFound(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	mockRepo.On("FindByID", userID).Return(modelpostgre.User{}, repo.ErrUserNotFound)

	req := httptest.NewRequest("GET", "/users/"+userID.String(), nil)

	app.Get("/users/:id", service.GetUserByID)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	createReq := modelpostgre.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
		FullName: "New User",
		RoleName: "user",
		IsActive: true,
	}

	mockRepo.On("Create", mock.AnythingOfType("modelpostgre.User"), "user").Return(createTestUser(), nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Post("/users", service.CreateUser)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_InvalidRequestBody(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	req := httptest.NewRequest("POST", "/users", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	app.Post("/users", service.CreateUser)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCreateUser_RoleNotFound(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	createReq := modelpostgre.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
		FullName: "New User",
		RoleName: "nonexistent",
		IsActive: true,
	}

	mockRepo.On("Create", mock.AnythingOfType("modelpostgre.User"), "nonexistent").Return(modelpostgre.User{}, repo.ErrRoleNotFound)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Post("/users", service.CreateUser)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_DuplicateUsernameOrEmail(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	createReq := modelpostgre.CreateUserRequest{
		Username: "existinguser",
		Email:    "existing@example.com",
		Password: "password123",
		FullName: "Existing User",
		RoleName: "user",
		IsActive: true,
	}

	mockRepo.On("Create", mock.AnythingOfType("modelpostgre.User"), "user").Return(modelpostgre.User{}, repo.ErrDuplicateUsernameOrEmail)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Post("/users", service.CreateUser)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := createTestUser()
	updateReq := modelpostgre.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
		FullName: "Updated User",
		RoleName: "admin",
	}

	mockRepo.On("FindByID", user.ID).Return(user, nil)
	mockRepo.On("Update", mock.AnythingOfType("modelpostgre.User"), "admin").Return(user, nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/users/"+user.ID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Put("/users/:id", service.UpdateUser)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_InvalidID(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	req := httptest.NewRequest("PUT", "/users/invalid-id", nil)

	app.Put("/users/:id", service.UpdateUser)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUpdateUser_UserNotFound(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	mockRepo.On("FindByID", userID).Return(modelpostgre.User{}, repo.ErrUserNotFound)

	updateReq := modelpostgre.UpdateUserRequest{
		Username: "updateduser",
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Put("/users/:id", service.UpdateUser)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	mockRepo.On("DeleteUser", userID).Return(nil)

	req := httptest.NewRequest("DELETE", "/users/"+userID.String(), nil)

	app.Delete("/users/:id", service.DeleteUser)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_InvalidID(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	req := httptest.NewRequest("DELETE", "/users/invalid-id", nil)

	app.Delete("/users/:id", service.DeleteUser)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestDeleteUser_UserNotFound(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	mockRepo.On("DeleteUser", userID).Return(repo.ErrUserNotFound)

	req := httptest.NewRequest("DELETE", "/users/"+userID.String(), nil)

	app.Delete("/users/:id", service.DeleteUser)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUserRole_Success(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	roleID := uuid.New()
	user := createTestUser()

	roleReq := modelpostgre.UpdateUserRoleRequest{
		RoleName: "admin",
	}

	mockRepo.On("FindRoleIDByName", "admin").Return(roleID, nil)
	mockRepo.On("UpdateRole", userID, roleID).Return(nil)
	mockRepo.On("FindByID", userID).Return(user, nil)

	body, _ := json.Marshal(roleReq)
	req := httptest.NewRequest("PATCH", "/users/"+userID.String()+"/role", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Patch("/users/:id/role", service.UpdateUserRole)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUserRole_InvalidID(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	req := httptest.NewRequest("PATCH", "/users/invalid-id/role", nil)

	app.Patch("/users/:id/role", service.UpdateUserRole)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUpdateUserRole_RoleNotFound(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	roleReq := modelpostgre.UpdateUserRoleRequest{
		RoleName: "nonexistent",
	}

	mockRepo.On("FindRoleIDByName", "nonexistent").Return(uuid.UUID{}, repo.ErrRoleNotFound)

	body, _ := json.Marshal(roleReq)
	req := httptest.NewRequest("PATCH", "/users/"+userID.String()+"/role", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Patch("/users/:id/role", service.UpdateUserRole)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUserRole_UserNotFound(t *testing.T) {
	app := fiber.New()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	roleID := uuid.New()
	roleReq := modelpostgre.UpdateUserRoleRequest{
		RoleName: "admin",
	}

	mockRepo.On("FindRoleIDByName", "admin").Return(roleID, nil)
	mockRepo.On("UpdateRole", userID, roleID).Return(repo.ErrUserNotFound)

	body, _ := json.Marshal(roleReq)
	req := httptest.NewRequest("PATCH", "/users/"+userID.String()+"/role", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Patch("/users/:id/role", service.UpdateUserRole)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}
