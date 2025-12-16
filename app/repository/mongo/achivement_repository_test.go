package repositorymongo

import (
	"context"
	"testing"
	"time"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*mongo.Client, *mongo.Database) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		t.Skip("MongoDB not available for testing")
	}

	db := client.Database("prestasi_db")
	return client, db
}

func cleanupTestDB(t *testing.T, client *mongo.Client, db *mongo.Database) {
	if err := db.Drop(context.Background()); err != nil {
		t.Logf("Failed to drop test database: %v", err)
	}
	if err := client.Disconnect(context.Background()); err != nil {
		t.Logf("Failed to disconnect: %v", err)
	}
}

func createTestAchievement() m.Achievement {
	level := m.LevelNational
	return m.Achievement{
		StudentID:       "student123",
		AchievementType: m.TypeCompetition,
		Title:           "Test Competition",
		Description:     "Test Description",
		Details: m.AchievementDetails{
			CompetitionName:  "National Math Olympiad",
			CompetitionLevel: &level,
		},
		Tags:        []string{"math", "competition"},
		Points:      100,
		Attachments: []m.Attachment{},
	}
}

func TestCreate_Success(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	achievement := createTestAchievement()

	result, err := repo.Create(achievement)

	assert.NoError(t, err)
	assert.NotEqual(t, primitive.NilObjectID, result.ID)
	assert.Equal(t, achievement.Title, result.Title)
	assert.Equal(t, achievement.StudentID, result.StudentID)
	assert.False(t, result.IsDeleted)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)
}

func TestGetByID_Success(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	achievement := createTestAchievement()
	created, _ := repo.Create(achievement)

	result, err := repo.GetByID(created.ID)

	assert.NoError(t, err)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, created.Title, result.Title)
	assert.Equal(t, created.StudentID, result.StudentID)
}

func TestGetByID_NotFound(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	fakeID := primitive.NewObjectID()

	_, err := repo.GetByID(fakeID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestGetByID_DeletedAchievement(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	achievement := createTestAchievement()
	created, _ := repo.Create(achievement)

	_ = repo.SoftDelete(created.ID)

	_, err := repo.GetByID(created.ID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestGetByStudentID_Success(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	studentID := "student123"
	achievement1 := createTestAchievement()
	achievement1.StudentID = studentID
	achievement2 := createTestAchievement()
	achievement2.StudentID = studentID
	achievement2.Title = "Second Achievement"

	_, _ = repo.Create(achievement1)
	time.Sleep(10 * time.Millisecond)
	_, _ = repo.Create(achievement2)

	results, err := repo.GetByStudentID(studentID)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "Second Achievement", results[0].Title)
}

func TestGetByStudentID_Empty(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	results, err := repo.GetByStudentID("nonexistent")

	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestGetMultipleByIDs_Success(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	achievement1 := createTestAchievement()
	achievement2 := createTestAchievement()
	achievement2.Title = "Second Achievement"

	created1, _ := repo.Create(achievement1)
	created2, _ := repo.Create(achievement2)

	ids := []primitive.ObjectID{created1.ID, created2.ID}

	results, err := repo.GetMultipleByIDs(ids)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestGetMultipleByIDs_EmptyInput(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	results, err := repo.GetMultipleByIDs([]primitive.ObjectID{})

	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestGetMultipleByIDs_PartialMatch(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	achievement := createTestAchievement()
	created, _ := repo.Create(achievement)

	fakeID := primitive.NewObjectID()
	ids := []primitive.ObjectID{created.ID, fakeID}

	results, err := repo.GetMultipleByIDs(ids)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, created.ID, results[0].ID)
}

func TestUpdate_Success(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	achievement := createTestAchievement()
	created, _ := repo.Create(achievement)

	created.Title = "Updated Title"
	created.Description = "Updated Description"
	created.Points = 200

	result, err := repo.Update(created.ID, created)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", result.Title)
	assert.Equal(t, "Updated Description", result.Description)
	assert.Equal(t, 200, result.Points)
	assert.True(t, result.UpdatedAt.After(created.CreatedAt))
}

func TestUpdate_NotFound(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	achievement := createTestAchievement()
	fakeID := primitive.NewObjectID()

	_, err := repo.Update(fakeID, achievement)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestUpdate_DeletedAchievement(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	achievement := createTestAchievement()
	created, _ := repo.Create(achievement)

	_ = repo.SoftDelete(created.ID)

	created.Title = "Should Not Update"
	_, err := repo.Update(created.ID, created)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestSoftDelete_Success(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	achievement := createTestAchievement()
	created, _ := repo.Create(achievement)

	err := repo.SoftDelete(created.ID)

	assert.NoError(t, err)

	var deleted m.Achievement
	filter := bson.M{"_id": created.ID}
	_ = collection.FindOne(context.Background(), filter).Decode(&deleted)
	assert.True(t, deleted.IsDeleted)
}

func TestSoftDelete_NotFound(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	fakeID := primitive.NewObjectID()

	err := repo.SoftDelete(fakeID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestSoftDelete_AlreadyDeleted(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	achievement := createTestAchievement()
	created, _ := repo.Create(achievement)

	_ = repo.SoftDelete(created.ID)
	err := repo.SoftDelete(created.ID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already deleted")
}

func TestGetAll_Success(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	for i := 0; i < 5; i++ {
		achievement := createTestAchievement()
		_, _ = repo.Create(achievement)
		time.Sleep(5 * time.Millisecond)
	}

	results, total, err := repo.GetAll(nil, 0, 10)

	assert.NoError(t, err)
	assert.Len(t, results, 5)
	assert.Equal(t, int64(5), total)
}

func TestGetAll_WithPagination(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	for i := 0; i < 10; i++ {
		achievement := createTestAchievement()
		_, _ = repo.Create(achievement)
		time.Sleep(5 * time.Millisecond)
	}

	results, total, err := repo.GetAll(nil, 0, 5)

	assert.NoError(t, err)
	assert.Len(t, results, 5)
	assert.Equal(t, int64(10), total)

	results2, total2, err2 := repo.GetAll(nil, 5, 5)

	assert.NoError(t, err2)
	assert.Len(t, results2, 5)
	assert.Equal(t, int64(10), total2)
	assert.NotEqual(t, results[0].ID, results2[0].ID)
}

func TestGetAll_WithFilter(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	student1 := createTestAchievement()
	student1.StudentID = "student1"
	student2 := createTestAchievement()
	student2.StudentID = "student2"

	_, _ = repo.Create(student1)
	_, _ = repo.Create(student2)
	_, _ = repo.Create(student1)

	filter := bson.M{"studentId": "student1"}
	results, total, err := repo.GetAll(filter, 0, 10)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, int64(2), total)
}

func TestGetAll_ExcludesDeleted(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	achievement1 := createTestAchievement()
	achievement2 := createTestAchievement()

	created1, _ := repo.Create(achievement1)
	_, _ = repo.Create(achievement2)

	_ = repo.SoftDelete(created1.ID)

	results, total, err := repo.GetAll(nil, 0, 10)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, int64(1), total)
}

func TestAddAttachment_Success(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	achievement := createTestAchievement()
	created, _ := repo.Create(achievement)

	attachment := m.Attachment{
		FileName:   "certificate.pdf",
		FileURL:    "https://example.com/file.pdf",
		FileType:   "application/pdf",
		FileSize:   1024,
		UploadedAt: time.Now(),
	}

	err := repo.AddAttachment(created.ID, attachment)

	assert.NoError(t, err)

	updated, _ := repo.GetByID(created.ID)
	assert.Len(t, updated.Attachments, 1)
	assert.Equal(t, "certificate.pdf", updated.Attachments[0].FileName)
}

func TestAddAttachment_NotFound(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	fakeID := primitive.NewObjectID()
	attachment := m.Attachment{
		FileName: "test.pdf",
		FileURL:  "https://example.com/test.pdf",
	}

	err := repo.AddAttachment(fakeID, attachment)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestAddAttachment_MultipleAttachments(t *testing.T) {
	client, db := setupTestDB(t)
	defer cleanupTestDB(t, client, db)

	collection := db.Collection("achievements")
	repo := &achievementRepository{collection: collection}

	achievement := createTestAchievement()
	created, _ := repo.Create(achievement)

	attachment1 := m.Attachment{
		FileName: "file1.pdf",
		FileURL:  "https://example.com/file1.pdf",
	}
	attachment2 := m.Attachment{
		FileName: "file2.pdf",
		FileURL:  "https://example.com/file2.pdf",
	}

	_ = repo.AddAttachment(created.ID, attachment1)
	_ = repo.AddAttachment(created.ID, attachment2)

	updated, _ := repo.GetByID(created.ID)
	assert.Len(t, updated.Attachments, 2)
}
