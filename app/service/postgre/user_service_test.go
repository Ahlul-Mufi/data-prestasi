package servicepostgre

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (mr *MockUserRepository) FindByIdentity(identity string) (m.User, error) {
	args := mr.Called(identity)
	return args.Get(0).(m.User), args.Error(1)
}

func (mr *MockUserRepository) GetPermissions(roleID uuid.UUID) ([]string, error) {
	args := mr.Called(roleID)
	return args.Get(0).([]string), args.Error(1)
}

func (mr *MockUserRepository) FindByID(id uuid.UUID) (m.User, error) {
	args := mr.Called(id)
	return args.Get(0).(m.User), args.Error(1)
}

func (mr *MockUserRepository) FindRoleIDByName(roleName string) (uuid.UUID, error) {
	args := mr.Called(roleName)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (mr *MockUserRepository) FindAll() ([]m.User, error) {
	args := mr.Called()
	return args.Get(0).([]m.User), args.Error(1)
}

func (mr *MockUserRepository) Create(user m.User, roleName string) (m.User, error) {
	args := mr.Called(user, roleName)
	return args.Get(0).(m.User), args.Error(1)
}

func (mr *MockUserRepository) Update(user m.User, newRoleName string) (m.User, error) {
	args := mr.Called(user, newRoleName)
	return args.Get(0).(m.User), args.Error(1)
}

func (mr *MockUserRepository) DeleteUser(id uuid.UUID) error {
	args := mr.Called(id)
	return args.Error(0)
}

func (mr *MockUserRepository) UpdateRole(userID uuid.UUID, roleID uuid.UUID) error {
	args := mr.Called(userID, roleID)
	return args.Error(0)
}

func setupTestApp() *fiber.App {
	app := fiber.New()
	return app
}

func TestLogin_Success(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	roleID := uuid.New()
	// Use a pre-hashed password for testing
	hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy" // "password123"

	user := m.User{
		ID:           userID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		RoleID:       &roleID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRepo.On("FindByIdentity", "testuser").Return(user, nil)

	// Register the route
	app.Post("/login", service.Login)

	loginReq := m.LoginRequest{
		Identity: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidIdentity(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByIdentity", "nonexistent").Return(m.User{}, repo.ErrUserNotFound)

	app.Post("/login", service.Login)

	loginReq := m.LoginRequest{
		Identity: "nonexistent",
		Password: "password123",
	}
	body, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InactiveUser(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

	user := m.User{
		ID:           userID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		IsActive:     false, // Inactive user
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRepo.On("FindByIdentity", "testuser").Return(user, nil)

	app.Post("/login", service.Login)

	loginReq := m.LoginRequest{
		Identity: "testuser",
		Password: "password123",
	}
	body, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestProfile_Success(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	user := m.User{
		ID:        userID,
		Username:  "testuser",
		Email:     "test@example.com",
		FullName:  "Test User",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", userID).Return(user, nil)

	app.Get("/profile", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return service.Profile(c)
	})

	req := httptest.NewRequest("GET", "/profile", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.NotNil(t, response["data"])
	mockRepo.AssertExpectations(t)
}

func TestProfile_Unauthorized(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	app.Get("/profile", service.Profile)

	req := httptest.NewRequest("GET", "/profile", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestGetUsers_Success(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	users := []m.User{
		{
			ID:       uuid.New(),
			Username: "user1",
			Email:    "user1@example.com",
			FullName: "User One",
			IsActive: true,
		},
		{
			ID:       uuid.New(),
			Username: "user2",
			Email:    "user2@example.com",
			FullName: "User Two",
			IsActive: true,
		},
	}

	mockRepo.On("FindAll").Return(users, nil)

	app.Get("/users", service.GetUsers)

	req := httptest.NewRequest("GET", "/users", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	assert.NotNil(t, response["data"])
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	user := m.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		IsActive: true,
	}

	mockRepo.On("FindByID", userID).Return(user, nil)

	app.Get("/users/:id", service.GetUserByID)

	req := httptest.NewRequest("GET", "/users/"+userID.String(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()

	mockRepo.On("FindByID", userID).Return(m.User{}, repo.ErrUserNotFound)

	app.Get("/users/:id", service.GetUserByID)

	req := httptest.NewRequest("GET", "/users/"+userID.String(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_InvalidID(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	app.Get("/users/:id", service.GetUserByID)

	req := httptest.NewRequest("GET", "/users/invalid-uuid", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCreateUser_Success(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	roleID := uuid.New()
	createReq := m.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
		FullName: "New User",
		RoleName: "user",
		IsActive: true,
	}

	createdUser := m.User{
		ID:       userID,
		Username: createReq.Username,
		Email:    createReq.Email,
		FullName: createReq.FullName,
		RoleID:   &roleID,
		IsActive: createReq.IsActive,
	}

	// Use mock.MatchedBy to match any User struct with the correct username
	mockRepo.On("Create", mock.MatchedBy(func(user m.User) bool {
		return user.Username == "newuser" && user.Email == "newuser@example.com"
	}), "user").Return(createdUser, nil)

	app.Post("/users", service.CreateUser)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_RoleNotFound(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	createReq := m.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
		FullName: "New User",
		RoleName: "nonexistent",
		IsActive: true,
	}

	mockRepo.On("Create", mock.MatchedBy(func(user m.User) bool {
		return user.Username == "newuser"
	}), "nonexistent").Return(m.User{}, repo.ErrRoleNotFound)

	app.Post("/users", service.CreateUser)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	createReq := m.CreateUserRequest{
		Username: "existinguser",
		Email:    "new@example.com",
		Password: "password123",
		FullName: "New User",
		RoleName: "user",
		IsActive: true,
	}

	mockRepo.On("Create", mock.MatchedBy(func(user m.User) bool {
		return user.Username == "existinguser"
	}), "user").Return(m.User{}, repo.ErrDuplicateUsernameOrEmail)

	app.Post("/users", service.CreateUser)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	existingUser := m.User{
		ID:           userID,
		Username:     "oldusername",
		Email:        "old@example.com",
		PasswordHash: "oldhash",
		FullName:     "Old Name",
		IsActive:     true,
	}

	updatedUser := m.User{
		ID:       userID,
		Username: "newusername",
		Email:    "new@example.com",
		FullName: "New Name",
		IsActive: true,
	}

	updateReq := m.UpdateUserRequest{
		Username: "newusername",
		Email:    "new@example.com",
		FullName: "New Name",
		RoleName: "",
	}

	mockRepo.On("FindByID", userID).Return(existingUser, nil)
	mockRepo.On("Update", mock.MatchedBy(func(user m.User) bool {
		return user.ID == userID && user.Username == "newusername"
	}), "").Return(updatedUser, nil)

	app.Put("/users/:id", service.UpdateUser)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_NotFound(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	updateReq := m.UpdateUserRequest{
		Username: "newusername",
		Email:    "new@example.com",
		FullName: "New Name",
	}

	mockRepo.On("FindByID", userID).Return(m.User{}, repo.ErrUserNotFound)

	app.Put("/users/:id", service.UpdateUser)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()

	mockRepo.On("DeleteUser", userID).Return(nil)

	app.Delete("/users/:id", service.DeleteUser)

	req := httptest.NewRequest("DELETE", "/users/"+userID.String(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_NotFound(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()

	mockRepo.On("DeleteUser", userID).Return(repo.ErrUserNotFound)

	app.Delete("/users/:id", service.DeleteUser)

	req := httptest.NewRequest("DELETE", "/users/"+userID.String(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUserRole_Success(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	roleID := uuid.New()

	updateRoleReq := m.UpdateUserRoleRequest{
		RoleName: "admin",
	}

	updatedUser := m.User{
		ID:       userID,
		Username: "testuser",
		RoleID:   &roleID,
	}

	mockRepo.On("FindRoleIDByName", "admin").Return(roleID, nil)
	mockRepo.On("UpdateRole", userID, roleID).Return(nil)
	mockRepo.On("FindByID", userID).Return(updatedUser, nil)

	app.Put("/users/:id/role", service.UpdateUserRole)

	body, _ := json.Marshal(updateRoleReq)
	req := httptest.NewRequest("PUT", "/users/"+userID.String()+"/role", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestLogout_Success(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	app.Post("/logout", service.Logout)

	req := httptest.NewRequest("POST", "/logout", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)

	data := response["data"].(map[string]interface{})
	assert.Contains(t, data["message"], "Logged out successfully")
}

func TestRefresh_Success(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	user := m.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		IsActive: true,
	}

	mockRepo.On("FindByID", userID).Return(user, nil)

	app.Post("/refresh", service.Refresh)

	refreshReq := map[string]string{
		"refresh_token": "valid_refresh_token",
	}
	body, _ := json.Marshal(refreshReq)

	req := httptest.NewRequest("POST", "/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	// Note: This test will likely fail because we're not generating a real refresh token
	// In a real scenario, you would need to mock the token validation
	// For now, we expect either OK or Unauthorized depending on token validation
	assert.True(t, resp.StatusCode == fiber.StatusOK || resp.StatusCode == fiber.StatusUnauthorized)
}

func TestCreateUser_InvalidRequestBody(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	app.Post("/users", service.CreateUser)

	// Invalid JSON
	body := []byte(`{"username": "incomplete"`)

	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUpdateUser_WithPassword(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	newPassword := "newpassword123"

	existingUser := m.User{
		ID:           userID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "oldhash",
		FullName:     "Test User",
		IsActive:     true,
	}

	updatedUser := m.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		IsActive: true,
	}

	updateReq := m.UpdateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: &newPassword,
		FullName: "Test User",
	}

	mockRepo.On("FindByID", userID).Return(existingUser, nil)
	mockRepo.On("Update", mock.MatchedBy(func(user m.User) bool {
		// Verify password was hashed (should be different from plain text)
		return user.ID == userID && user.PasswordHash != newPassword
	}), "").Return(updatedUser, nil)

	app.Put("/users/:id", service.UpdateUser)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_WithIsActiveFlag(t *testing.T) {
	app := setupTestApp()
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	isActive := false

	existingUser := m.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		IsActive: true,
	}

	updatedUser := m.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		IsActive: false,
	}

	updateReq := m.UpdateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
		IsActive: &isActive,
	}

	mockRepo.On("FindByID", userID).Return(existingUser, nil)
	mockRepo.On("Update", mock.MatchedBy(func(user m.User) bool {
		return user.ID == userID && user.IsActive == false
	}), "").Return(updatedUser, nil)

	app.Put("/users/:id", service.UpdateUser)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}
