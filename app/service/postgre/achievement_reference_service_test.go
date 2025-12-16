package servicepostgre

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	modelmongo "github.com/Ahlul-Mufi/data-prestasi/app/model/mongo"
	modelpostgre "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockAchievementReferenceRepository struct {
	mock.Mock
}

func (m *MockAchievementReferenceRepository) Create(a modelpostgre.AchievementReference) (modelpostgre.AchievementReference, error) {
	args := m.Called(a)
	return args.Get(0).(modelpostgre.AchievementReference), args.Error(1)
}

func (m *MockAchievementReferenceRepository) GetByID(id uuid.UUID) (modelpostgre.AchievementReference, error) {
	args := m.Called(id)
	return args.Get(0).(modelpostgre.AchievementReference), args.Error(1)
}

func (m *MockAchievementReferenceRepository) Update(a modelpostgre.AchievementReference) (modelpostgre.AchievementReference, error) {
	args := m.Called(a)
	return args.Get(0).(modelpostgre.AchievementReference), args.Error(1)
}

func (m *MockAchievementReferenceRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAchievementReferenceRepository) GetFiltered(userID *uuid.UUID, status *modelpostgre.AchievementStatus) ([]modelpostgre.AchievementReference, error) {
	args := m.Called(userID, status)
	return args.Get(0).([]modelpostgre.AchievementReference), args.Error(1)
}

func (m *MockAchievementReferenceRepository) UpdateStatus(id, verifierID uuid.UUID, newStatus modelpostgre.AchievementStatus, rejectionNote *string) (modelpostgre.AchievementReference, error) {
	args := m.Called(id, verifierID, newStatus, rejectionNote)
	return args.Get(0).(modelpostgre.AchievementReference), args.Error(1)
}

func (m *MockAchievementReferenceRepository) FindByStudentID(studentID uuid.UUID) ([]modelpostgre.AchievementReference, error) {
	args := m.Called(studentID)
	return args.Get(0).([]modelpostgre.AchievementReference), args.Error(1)
}

func (m *MockAchievementReferenceRepository) CreateHistory(h modelpostgre.AchievementHistory) error {
	args := m.Called(h)
	return args.Error(0)
}

func (m *MockAchievementReferenceRepository) FindHistoryByAchievementRefID(refID uuid.UUID) ([]modelpostgre.AchievementHistory, error) {
	args := m.Called(refID)
	return args.Get(0).([]modelpostgre.AchievementHistory), args.Error(1)
}

type MockMongoAchievementRepository struct {
	mock.Mock
}

func (m *MockMongoAchievementRepository) Create(achievement modelmongo.Achievement) (modelmongo.Achievement, error) {
	args := m.Called(achievement)
	return args.Get(0).(modelmongo.Achievement), args.Error(1)
}

func (m *MockMongoAchievementRepository) GetByID(id primitive.ObjectID) (modelmongo.Achievement, error) {
	args := m.Called(id)
	return args.Get(0).(modelmongo.Achievement), args.Error(1)
}

func (m *MockMongoAchievementRepository) GetByStudentID(studentID string) ([]modelmongo.Achievement, error) {
	args := m.Called(studentID)
	return args.Get(0).([]modelmongo.Achievement), args.Error(1)
}

func (m *MockMongoAchievementRepository) GetMultipleByIDs(ids []primitive.ObjectID) ([]modelmongo.Achievement, error) {
	args := m.Called(ids)
	return args.Get(0).([]modelmongo.Achievement), args.Error(1)
}

func (m *MockMongoAchievementRepository) Update(id primitive.ObjectID, achievement modelmongo.Achievement) (modelmongo.Achievement, error) {
	args := m.Called(id, achievement)
	return args.Get(0).(modelmongo.Achievement), args.Error(1)
}

func (m *MockMongoAchievementRepository) SoftDelete(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMongoAchievementRepository) GetAll(filter bson.M, skip, limit int64) ([]modelmongo.Achievement, int64, error) {
	args := m.Called(filter, skip, limit)
	return args.Get(0).([]modelmongo.Achievement), args.Get(1).(int64), args.Error(2)
}

func (m *MockMongoAchievementRepository) AddAttachment(id primitive.ObjectID, attachment modelmongo.Attachment) error {
	args := m.Called(id, attachment)
	return args.Error(0)
}

type MockStudentRepository struct {
	mock.Mock
}

func (m *MockStudentRepository) UpdateAdvisor(id uuid.UUID, advisorID uuid.UUID) error {
	panic("unimplemented")
}

func (m *MockStudentRepository) GetByUserID(userID uuid.UUID) (modelpostgre.Student, error) {
	args := m.Called(userID)
	return args.Get(0).(modelpostgre.Student), args.Error(1)
}

func (m *MockStudentRepository) Create(student modelpostgre.Student) (modelpostgre.Student, error) {
	args := m.Called(student)
	return args.Get(0).(modelpostgre.Student), args.Error(1)
}

func (m *MockStudentRepository) GetByID(id uuid.UUID) (modelpostgre.Student, error) {
	args := m.Called(id)
	return args.Get(0).(modelpostgre.Student), args.Error(1)
}

func (m *MockStudentRepository) GetAll() ([]modelpostgre.Student, error) {
	args := m.Called()
	return args.Get(0).([]modelpostgre.Student), args.Error(1)
}

func (m *MockStudentRepository) Update(student modelpostgre.Student) (modelpostgre.Student, error) {
	args := m.Called(student)
	return args.Get(0).(modelpostgre.Student), args.Error(1)
}

func (m *MockStudentRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockLecturerRepository struct {
	mock.Mock
}

func (m *MockLecturerRepository) Create(lecturer modelpostgre.Lecturer) (modelpostgre.Lecturer, error) {
	panic("unimplemented")
}

func (m *MockLecturerRepository) Delete(id uuid.UUID) error {
	panic("unimplemented")
}

func (m *MockLecturerRepository) GetAdvisees(lecturerID uuid.UUID) ([]modelpostgre.Student, error) {
	panic("unimplemented")
}

func (m *MockLecturerRepository) GetAll() ([]modelpostgre.Lecturer, error) {
	panic("unimplemented")
}

func (m *MockLecturerRepository) GetByID(id uuid.UUID) (modelpostgre.Lecturer, error) {
	panic("unimplemented")
}

func (m *MockLecturerRepository) Update(lecturer modelpostgre.Lecturer) (modelpostgre.Lecturer, error) {
	panic("unimplemented")
}

func (m *MockLecturerRepository) GetByUserID(userID uuid.UUID) (modelpostgre.Lecturer, error) {
	args := m.Called(userID)
	return args.Get(0).(modelpostgre.Lecturer), args.Error(1)
}

func createTestStudent() modelpostgre.Student {
	return modelpostgre.Student{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}
}

func createTestAchievement() modelmongo.Achievement {
	return modelmongo.Achievement{
		ID:              primitive.NewObjectID(),
		StudentID:       uuid.New().String(),
		AchievementType: modelmongo.TypeCompetition,
		Title:           "Test Achievement",
		Description:     "Test Description",
		Points:          100,
	}
}

func createTestAchievementRef() modelpostgre.AchievementReference {
	return modelpostgre.AchievementReference{
		ID:                 uuid.New(),
		StudentID:          uuid.New(),
		MongoAchievementID: primitive.NewObjectID().Hex(),
		Status:             modelpostgre.StatusDraft,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

func TestCreate_Success(t *testing.T) {
	app := fiber.New()
	mockPostgresRepo := new(MockAchievementReferenceRepository)
	mockMongoRepo := new(MockMongoAchievementRepository)
	mockStudentRepo := new(MockStudentRepository)
	mockLecturerRepo := new(MockLecturerRepository)

	service := NewAchievementReferenceService(mockPostgresRepo, mockMongoRepo, mockStudentRepo, mockLecturerRepo)

	userID := uuid.New()
	student := createTestStudent()
	student.UserID = userID

	mockStudentRepo.On("GetByUserID", userID).Return(student, nil)

	mongoAchievement := createTestAchievement()
	mockMongoRepo.On("Create", mock.AnythingOfType("modelmongo.Achievement")).Return(mongoAchievement, nil)

	refResult := createTestAchievementRef()
	refResult.StudentID = student.ID
	mockPostgresRepo.On("Create", mock.AnythingOfType("modelpostgre.AchievementReference")).Return(refResult, nil)
	mockPostgresRepo.On("CreateHistory", mock.AnythingOfType("modelpostgre.AchievementHistory")).Return(nil)

	createReq := modelmongo.CreateAchievementRequest{
		AchievementType: modelmongo.TypeCompetition,
		Title:           "Test Achievement",
		Description:     "Test Description",
		Points:          100,
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/achievements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return service.Create(c)
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockMongoRepo.AssertExpectations(t)
	mockPostgresRepo.AssertExpectations(t)
}

func TestCreate_InvalidUser(t *testing.T) {
	app := fiber.New()
	mockPostgresRepo := new(MockAchievementReferenceRepository)
	mockMongoRepo := new(MockMongoAchievementRepository)
	mockStudentRepo := new(MockStudentRepository)
	mockLecturerRepo := new(MockLecturerRepository)

	service := NewAchievementReferenceService(mockPostgresRepo, mockMongoRepo, mockStudentRepo, mockLecturerRepo)

	userID := uuid.New()
	mockStudentRepo.On("GetByUserID", userID).Return(modelpostgre.Student{}, errors.New("not found"))

	createReq := modelmongo.CreateAchievementRequest{
		Title:       "Test",
		Description: "Test",
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/achievements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return service.Create(c)
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
}

func TestCreate_ValidationError(t *testing.T) {
	app := fiber.New()
	mockPostgresRepo := new(MockAchievementReferenceRepository)
	mockMongoRepo := new(MockMongoAchievementRepository)
	mockStudentRepo := new(MockStudentRepository)
	mockLecturerRepo := new(MockLecturerRepository)

	service := NewAchievementReferenceService(mockPostgresRepo, mockMongoRepo, mockStudentRepo, mockLecturerRepo)

	userID := uuid.New()
	student := createTestStudent()
	mockStudentRepo.On("GetByUserID", userID).Return(student, nil)

	createReq := modelmongo.CreateAchievementRequest{
		Title: "",
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/achievements", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return service.Create(c)
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
}

func TestGetByID_Success(t *testing.T) {
	app := fiber.New()
	mockPostgresRepo := new(MockAchievementReferenceRepository)
	mockMongoRepo := new(MockMongoAchievementRepository)
	mockStudentRepo := new(MockStudentRepository)
	mockLecturerRepo := new(MockLecturerRepository)

	service := NewAchievementReferenceService(mockPostgresRepo, mockMongoRepo, mockStudentRepo, mockLecturerRepo)

	refID := uuid.New()
	ref := createTestAchievementRef()
	ref.ID = refID

	mongoAchievement := createTestAchievement()
	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)

	mockPostgresRepo.On("GetByID", refID).Return(ref, nil)
	mockMongoRepo.On("GetByID", mongoID).Return(mongoAchievement, nil)

	req := httptest.NewRequest("GET", "/achievements/"+refID.String(), nil)

	app.Get("/achievements/:id", service.GetByID)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockPostgresRepo.AssertExpectations(t)
	mockMongoRepo.AssertExpectations(t)
}

func TestGetByID_NotFound(t *testing.T) {
	app := fiber.New()
	mockPostgresRepo := new(MockAchievementReferenceRepository)
	mockMongoRepo := new(MockMongoAchievementRepository)
	mockStudentRepo := new(MockStudentRepository)
	mockLecturerRepo := new(MockLecturerRepository)

	service := NewAchievementReferenceService(mockPostgresRepo, mockMongoRepo, mockStudentRepo, mockLecturerRepo)

	refID := uuid.New()
	mockPostgresRepo.On("GetByID", refID).Return(modelpostgre.AchievementReference{}, errors.New("achievement reference not found"))

	req := httptest.NewRequest("GET", "/achievements/"+refID.String(), nil)

	app.Get("/achievements/:id", service.GetByID)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockPostgresRepo.AssertExpectations(t)
}

func TestGetMyAchievements_Success(t *testing.T) {
	app := fiber.New()
	mockPostgresRepo := new(MockAchievementReferenceRepository)
	mockMongoRepo := new(MockMongoAchievementRepository)
	mockStudentRepo := new(MockStudentRepository)
	mockLecturerRepo := new(MockLecturerRepository)

	service := NewAchievementReferenceService(mockPostgresRepo, mockMongoRepo, mockStudentRepo, mockLecturerRepo)

	userID := uuid.New()
	student := createTestStudent()
	student.UserID = userID

	refs := []modelpostgre.AchievementReference{createTestAchievementRef(), createTestAchievementRef()}
	achievements := []modelmongo.Achievement{createTestAchievement(), createTestAchievement()}

	mockStudentRepo.On("GetByUserID", userID).Return(student, nil)
	mockPostgresRepo.On("FindByStudentID", student.ID).Return(refs, nil)
	mockMongoRepo.On("GetMultipleByIDs", mock.AnythingOfType("[]primitive.ObjectID")).Return(achievements, nil)

	req := httptest.NewRequest("GET", "/my-achievements", nil)

	app.Get("/my-achievements", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return service.GetMyAchievements(c)
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockPostgresRepo.AssertExpectations(t)
}

func TestSubmit_Success(t *testing.T) {
	app := fiber.New()
	mockPostgresRepo := new(MockAchievementReferenceRepository)
	mockMongoRepo := new(MockMongoAchievementRepository)
	mockStudentRepo := new(MockStudentRepository)
	mockLecturerRepo := new(MockLecturerRepository)

	service := NewAchievementReferenceService(mockPostgresRepo, mockMongoRepo, mockStudentRepo, mockLecturerRepo)

	userID := uuid.New()
	student := createTestStudent()
	student.UserID = userID

	refID := uuid.New()
	ref := createTestAchievementRef()
	ref.ID = refID
	ref.StudentID = student.ID
	ref.Status = modelpostgre.StatusDraft

	updatedRef := ref
	updatedRef.Status = modelpostgre.StatusSubmitted

	mockStudentRepo.On("GetByUserID", userID).Return(student, nil)
	mockPostgresRepo.On("GetByID", refID).Return(ref, nil)
	mockPostgresRepo.On("Update", mock.AnythingOfType("modelpostgre.AchievementReference")).Return(updatedRef, nil)
	mockPostgresRepo.On("CreateHistory", mock.AnythingOfType("modelpostgre.AchievementHistory")).Return(nil)

	req := httptest.NewRequest("POST", "/achievements/"+refID.String()+"/submit", nil)

	app.Post("/achievements/:id/submit", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return service.Submit(c)
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockPostgresRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
}

func TestSubmit_InvalidStatus(t *testing.T) {
	app := fiber.New()
	mockPostgresRepo := new(MockAchievementReferenceRepository)
	mockMongoRepo := new(MockMongoAchievementRepository)
	mockStudentRepo := new(MockStudentRepository)
	mockLecturerRepo := new(MockLecturerRepository)

	service := NewAchievementReferenceService(mockPostgresRepo, mockMongoRepo, mockStudentRepo, mockLecturerRepo)

	userID := uuid.New()
	student := createTestStudent()
	student.UserID = userID

	refID := uuid.New()
	ref := createTestAchievementRef()
	ref.ID = refID
	ref.StudentID = student.ID
	ref.Status = modelpostgre.StatusVerified

	mockStudentRepo.On("GetByUserID", userID).Return(student, nil)
	mockPostgresRepo.On("GetByID", refID).Return(ref, nil)

	req := httptest.NewRequest("POST", "/achievements/"+refID.String()+"/submit", nil)

	app.Post("/achievements/:id/submit", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return service.Submit(c)
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
	mockPostgresRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
}

func TestSubmit_Unauthorized(t *testing.T) {
	app := fiber.New()
	mockPostgresRepo := new(MockAchievementReferenceRepository)
	mockMongoRepo := new(MockMongoAchievementRepository)
	mockStudentRepo := new(MockStudentRepository)
	mockLecturerRepo := new(MockLecturerRepository)

	service := NewAchievementReferenceService(mockPostgresRepo, mockMongoRepo, mockStudentRepo, mockLecturerRepo)

	userID := uuid.New()
	student := createTestStudent()
	student.UserID = userID

	refID := uuid.New()
	ref := createTestAchievementRef()
	ref.ID = refID
	ref.StudentID = uuid.New()
	ref.Status = modelpostgre.StatusDraft

	mockStudentRepo.On("GetByUserID", userID).Return(student, nil)
	mockPostgresRepo.On("GetByID", refID).Return(ref, nil)

	req := httptest.NewRequest("POST", "/achievements/"+refID.String()+"/submit", nil)

	app.Post("/achievements/:id/submit", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return service.Submit(c)
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	mockPostgresRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	app := fiber.New()
	mockPostgresRepo := new(MockAchievementReferenceRepository)
	mockMongoRepo := new(MockMongoAchievementRepository)
	mockStudentRepo := new(MockStudentRepository)
	mockLecturerRepo := new(MockLecturerRepository)

	service := NewAchievementReferenceService(mockPostgresRepo, mockMongoRepo, mockStudentRepo, mockLecturerRepo)

	userID := uuid.New()
	student := createTestStudent()
	student.UserID = userID

	refID := uuid.New()
	ref := createTestAchievementRef()
	ref.ID = refID
	ref.StudentID = student.ID
	ref.Status = modelpostgre.StatusDraft

	mongoID, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
	mongoAchievement := createTestAchievement()
	mongoAchievement.ID = mongoID

	updatedAchievement := mongoAchievement
	updatedAchievement.Title = "Updated Title"

	mockStudentRepo.On("GetByUserID", userID).Return(student, nil)
	mockPostgresRepo.On("GetByID", refID).Return(ref, nil)
	mockMongoRepo.On("GetByID", mongoID).Return(mongoAchievement, nil)
	mockMongoRepo.On("Update", mongoID, mock.AnythingOfType("modelmongo.Achievement")).Return(updatedAchievement, nil)

	updateReq := modelmongo.UpdateAchievementRequest{
		Title: stringPtr("Updated Title"),
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/achievements/"+refID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Put("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return service.Update(c)
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockPostgresRepo.AssertExpectations(t)
	mockMongoRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
}

func stringPtr(s string) *string {
	return &s
}

func TestUpdate_InvalidStatus(t *testing.T) {
	app := fiber.New()
	mockPostgresRepo := new(MockAchievementReferenceRepository)
	mockMongoRepo := new(MockMongoAchievementRepository)
	mockStudentRepo := new(MockStudentRepository)
	mockLecturerRepo := new(MockLecturerRepository)

	service := NewAchievementReferenceService(mockPostgresRepo, mockMongoRepo, mockStudentRepo, mockLecturerRepo)

	userID := uuid.New()
	student := createTestStudent()
	student.UserID = userID

	refID := uuid.New()
	ref := createTestAchievementRef()
	ref.ID = refID
	ref.StudentID = student.ID
	ref.Status = modelpostgre.StatusVerified

	mockStudentRepo.On("GetByUserID", userID).Return(student, nil)
	mockPostgresRepo.On("GetByID", refID).Return(ref, nil)

	updateReq := modelmongo.UpdateAchievementRequest{}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/achievements/"+refID.String(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Put("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("user_id", userID.String())
		return service.Update(c)
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
	mockPostgresRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
}

func TestVerify_Success(t *testing.T) {
	app := fiber.New()
	mockPostgresRepo := new(MockAchievementReferenceRepository)
	mockMongoRepo := new(MockMongoAchievementRepository)
	mockStudentRepo := new(MockStudentRepository)
	mockLecturerRepo := new(MockLecturerRepository)

	service := NewAchievementReferenceService(mockPostgresRepo, mockMongoRepo, mockStudentRepo, mockLecturerRepo)

	verifierID := uuid.New()
	refID := uuid.New()
	ref := createTestAchievementRef()
	ref.ID = refID
	ref.Status = modelpostgre.StatusSubmitted

	verifiedRef := ref
	verifiedRef.Status = modelpostgre.StatusVerified

	mockPostgresRepo.On("GetByID", refID).Return(ref, nil)
	mockPostgresRepo.On("UpdateStatus", refID, verifierID, modelpostgre.StatusVerified, (*string)(nil)).Return(verifiedRef, nil)
	mockPostgresRepo.On("CreateHistory", mock.AnythingOfType("modelpostgre.AchievementHistory")).Return(nil)

	req := httptest.NewRequest("POST", "/achievements/"+refID.String()+"/verify", nil)

	app.Post("/achievements/:id/verify", func(c *fiber.Ctx) error {
		c.Locals("user_id", verifierID.String())
		return service.Verify(c)
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockPostgresRepo.AssertExpectations(t)
}

func TestReject_Success(t *testing.T) {
	app := fiber.New()
	mockPostgresRepo := new(MockAchievementReferenceRepository)
	mockMongoRepo := new(MockMongoAchievementRepository)
	mockStudentRepo := new(MockStudentRepository)
	mockLecturerRepo := new(MockLecturerRepository)

	service := NewAchievementReferenceService(mockPostgresRepo, mockMongoRepo, mockStudentRepo, mockLecturerRepo)

	verifierID := uuid.New()
	refID := uuid.New()
	ref := createTestAchievementRef()
	ref.ID = refID
	ref.Status = modelpostgre.StatusSubmitted

	rejectionNote := "Incomplete documentation"
	rejectedRef := ref
	rejectedRef.Status = modelpostgre.StatusRejected

	mockPostgresRepo.On("GetByID", refID).Return(ref, nil)
	mockPostgresRepo.On("UpdateStatus", refID, verifierID, modelpostgre.StatusRejected, &rejectionNote).Return(rejectedRef, nil)
	mockPostgresRepo.On("CreateHistory", mock.AnythingOfType("modelpostgre.AchievementHistory")).Return(nil)

	verifyReq := modelpostgre.VerificationRequest{
		RejectionNote: &rejectionNote,
	}

	body, _ := json.Marshal(verifyReq)
	req := httptest.NewRequest("POST", "/achievements/"+refID.String()+"/reject", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	app.Post("/achievements/:id/reject", func(c *fiber.Ctx) error {
		c.Locals("user_id", verifierID.String())
		return service.Reject(c)
	})

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockPostgresRepo.AssertExpectations(t)
}

func TestReject_MissingNote(t *testing.T) {
	app := fiber.New()
	mockPostgresRepo := new(MockAchievementReferenceRepository)
	mockMongoRepo := new(MockMongoAchievementRepository)
	mockStudentRepo := new(MockStudentRepository)
	mockLecturerRepo := new(MockLecturerRepository)

	service := NewAchievementReferenceService(mockPostgresRepo, mockMongoRepo, mockStudentRepo, mockLecturerRepo)

	verifierID := uuid.New()
	refID := uuid.New()
	ref := createTestAchievementRef()
	ref.ID = refID
	ref.Status = modelpostgre.StatusSubmitted
	mockPostgresRepo.On("GetByID", refID).Return(ref, nil)

	verifyReq := modelpostgre.VerificationRequest{
		RejectionNote: nil,
	}
	body, _ := json.Marshal(verifyReq)
	req := httptest.NewRequest("POST", "/achievements/"+refID.String()+"/reject", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	app.Post("/achievements/:id/reject", func(c *fiber.Ctx) error {
		c.Locals("user_id", verifierID.String())
		return service.Reject(c)
	})
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockPostgresRepo.AssertExpectations(t)
}
