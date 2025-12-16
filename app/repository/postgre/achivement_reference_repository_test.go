package repositorypostgre

import (
	"database/sql"
	"testing"
	"time"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func setupTestPostgresDB(t *testing.T) *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=prestasi_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err, "Failed to connect to PostgreSQL. Make sure PostgreSQL is running on localhost:5432")

	err = db.Ping()
	require.NoError(t, err, "Failed to ping PostgreSQL database")

	setupTestTables(t, db)

	return db
}

func setupTestTables(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS achievement_references (
			id UUID PRIMARY KEY,
			student_id UUID NOT NULL,
			mongo_achievement_id VARCHAR(255) NOT NULL,
			status VARCHAR(50) NOT NULL,
			submitted_at TIMESTAMP,
			verified_at TIMESTAMP,
			verified_by UUID,
			rejection_note TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS achievement_histories (
			id UUID PRIMARY KEY,
			achievement_ref_id UUID NOT NULL,
			previous_status VARCHAR(50) NOT NULL,
			new_status VARCHAR(50) NOT NULL,
			changed_by_user_id UUID NOT NULL,
			note TEXT,
			created_at TIMESTAMP NOT NULL,
			FOREIGN KEY (achievement_ref_id) REFERENCES achievement_references(id) ON DELETE CASCADE
		)
	`)
	require.NoError(t, err)
}

func cleanupTestPostgresDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`DROP TABLE IF EXISTS achievement_histories CASCADE`)
	if err != nil {
		t.Logf("Failed to drop achievement_histories table: %v", err)
	}

	_, err = db.Exec(`DROP TABLE IF EXISTS achievement_references CASCADE`)
	if err != nil {
		t.Logf("Failed to drop achievement_references table: %v", err)
	}

	if err := db.Close(); err != nil {
		t.Logf("Failed to close database: %v", err)
	}
}

func createTestAchievementReference() m.AchievementReference {
	return m.AchievementReference{
		StudentID:          uuid.New(),
		MongoAchievementID: "mongo_id_123",
		Status:             m.StatusDraft,
	}
}

func TestAchievementReference_Create_Success(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)
	ref := createTestAchievementReference()

	result, err := repo.Create(ref)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, ref.StudentID, result.StudentID)
	assert.Equal(t, ref.MongoAchievementID, result.MongoAchievementID)
	assert.Equal(t, m.StatusDraft, result.Status)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)
	assert.Nil(t, result.SubmittedAt)
}

func TestAchievementReference_Create_WithSubmittedStatus(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)
	ref := createTestAchievementReference()
	ref.Status = m.StatusSubmitted

	result, err := repo.Create(ref)

	assert.NoError(t, err)
	assert.Equal(t, m.StatusSubmitted, result.Status)
	assert.NotNil(t, result.SubmittedAt)
}

func TestAchievementReference_GetByID_Success(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)
	ref := createTestAchievementReference()
	created, _ := repo.Create(ref)

	result, err := repo.GetByID(created.ID)

	assert.NoError(t, err)
	assert.Equal(t, created.ID, result.ID)
	assert.Equal(t, created.StudentID, result.StudentID)
	assert.Equal(t, created.MongoAchievementID, result.MongoAchievementID)
}

func TestAchievementReference_GetByID_NotFound(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)
	fakeID := uuid.New()

	_, err := repo.GetByID(fakeID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestAchievementReference_Update_Success(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)
	ref := createTestAchievementReference()
	created, _ := repo.Create(ref)

	created.MongoAchievementID = "new_mongo_id"
	created.Status = m.StatusSubmitted
	now := time.Now()
	created.SubmittedAt = &now

	result, err := repo.Update(created)

	assert.NoError(t, err)
	assert.Equal(t, "new_mongo_id", result.MongoAchievementID)
	assert.Equal(t, m.StatusSubmitted, result.Status)
	assert.NotNil(t, result.SubmittedAt)
}

func TestAchievementReference_Update_NotFound(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)
	ref := createTestAchievementReference()
	ref.ID = uuid.New()

	_, err := repo.Update(ref)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
func TestAchievementReference_Delete_Success(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)
	ref := createTestAchievementReference()
	created, _ := repo.Create(ref)

	err := repo.Delete(created.ID)

	assert.NoError(t, err)

	_, err = repo.GetByID(created.ID)
	assert.Error(t, err)
}

func TestAchievementReference_Delete_NotFound(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)
	fakeID := uuid.New()

	err := repo.Delete(fakeID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestAchievementReference_GetFiltered_NoFilters(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	for i := 0; i < 3; i++ {
		ref := createTestAchievementReference()
		_, _ = repo.Create(ref)
	}

	results, err := repo.GetFiltered(nil, nil)

	assert.NoError(t, err)
	assert.Len(t, results, 3)
}

func TestAchievementReference_GetFiltered_ByStudentID(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	studentID1 := uuid.New()
	studentID2 := uuid.New()

	ref1 := createTestAchievementReference()
	ref1.StudentID = studentID1
	ref2 := createTestAchievementReference()
	ref2.StudentID = studentID1
	ref3 := createTestAchievementReference()
	ref3.StudentID = studentID2

	_, _ = repo.Create(ref1)
	_, _ = repo.Create(ref2)
	_, _ = repo.Create(ref3)

	results, err := repo.GetFiltered(&studentID1, nil)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	for _, r := range results {
		assert.Equal(t, studentID1, r.StudentID)
	}
}

func TestAchievementReference_GetFiltered_ByStatus(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	ref1 := createTestAchievementReference()
	ref1.Status = m.StatusDraft
	ref2 := createTestAchievementReference()
	ref2.Status = m.StatusSubmitted
	ref3 := createTestAchievementReference()
	ref3.Status = m.StatusSubmitted

	_, _ = repo.Create(ref1)
	_, _ = repo.Create(ref2)
	_, _ = repo.Create(ref3)

	status := m.StatusSubmitted
	results, err := repo.GetFiltered(nil, &status)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	for _, r := range results {
		assert.Equal(t, m.StatusSubmitted, r.Status)
	}
}

func TestAchievementReference_GetFiltered_ByStudentIDAndStatus(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	studentID := uuid.New()
	ref1 := createTestAchievementReference()
	ref1.StudentID = studentID
	ref1.Status = m.StatusDraft
	ref2 := createTestAchievementReference()
	ref2.StudentID = studentID
	ref2.Status = m.StatusSubmitted

	_, _ = repo.Create(ref1)
	_, _ = repo.Create(ref2)

	status := m.StatusSubmitted
	results, err := repo.GetFiltered(&studentID, &status)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, studentID, results[0].StudentID)
	assert.Equal(t, m.StatusSubmitted, results[0].Status)
}
func TestAchievementReference_UpdateStatus_ToVerified(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	ref := createTestAchievementReference()
	ref.Status = m.StatusSubmitted
	created, _ := repo.Create(ref)

	verifierID := uuid.New()
	result, err := repo.UpdateStatus(created.ID, verifierID, m.StatusVerified, nil)

	assert.NoError(t, err)
	assert.Equal(t, m.StatusVerified, result.Status)
	assert.NotNil(t, result.VerifiedAt)
	assert.NotNil(t, result.VerifiedBy)
	assert.Equal(t, verifierID, *result.VerifiedBy)
	assert.Nil(t, result.RejectionNote)
}

func TestAchievementReference_UpdateStatus_ToRejected(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	ref := createTestAchievementReference()
	ref.Status = m.StatusSubmitted
	created, _ := repo.Create(ref)

	verifierID := uuid.New()
	note := "Incomplete documentation"
	result, err := repo.UpdateStatus(created.ID, verifierID, m.StatusRejected, &note)

	assert.NoError(t, err)
	assert.Equal(t, m.StatusRejected, result.Status)
	assert.NotNil(t, result.VerifiedAt)
	assert.NotNil(t, result.VerifiedBy)
	assert.Equal(t, verifierID, *result.VerifiedBy)
	assert.NotNil(t, result.RejectionNote)
	assert.Equal(t, note, *result.RejectionNote)
}

func TestAchievementReference_UpdateStatus_NotSubmitted(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	ref := createTestAchievementReference()
	ref.Status = m.StatusDraft
	created, _ := repo.Create(ref)

	verifierID := uuid.New()
	_, err := repo.UpdateStatus(created.ID, verifierID, m.StatusVerified, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found or already processed")
}

func TestAchievementReference_UpdateStatus_NotFound(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	fakeID := uuid.New()
	verifierID := uuid.New()
	_, err := repo.UpdateStatus(fakeID, verifierID, m.StatusVerified, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found or already processed")
}

func TestAchievementReference_FindByStudentID_Success(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	studentID := uuid.New()
	ref1 := createTestAchievementReference()
	ref1.StudentID = studentID
	ref2 := createTestAchievementReference()
	ref2.StudentID = studentID

	_, _ = repo.Create(ref1)
	time.Sleep(10 * time.Millisecond)
	_, _ = repo.Create(ref2)

	results, err := repo.FindByStudentID(studentID)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.True(t, results[0].CreatedAt.After(results[1].CreatedAt) || results[0].CreatedAt.Equal(results[1].CreatedAt))
}

func TestAchievementReference_FindByStudentID_Empty(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	studentID := uuid.New()
	results, err := repo.FindByStudentID(studentID)

	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestAchievementReference_CreateHistory_Success(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	ref := createTestAchievementReference()
	created, _ := repo.Create(ref)

	history := m.AchievementHistory{
		AchievementRefID: created.ID,
		PreviousStatus:   m.StatusDraft,
		NewStatus:        m.StatusSubmitted,
		ChangedByUserID:  uuid.New(),
	}

	err := repo.CreateHistory(history)

	assert.NoError(t, err)
}

func TestAchievementReference_CreateHistory_WithNote(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	ref := createTestAchievementReference()
	created, _ := repo.Create(ref)

	note := "Changed status to submitted"
	history := m.AchievementHistory{
		AchievementRefID: created.ID,
		PreviousStatus:   m.StatusDraft,
		NewStatus:        m.StatusSubmitted,
		ChangedByUserID:  uuid.New(),
		Note:             &note,
	}

	err := repo.CreateHistory(history)

	assert.NoError(t, err)
}

func TestAchievementReference_FindHistoryByAchievementRefID_Success(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	ref := createTestAchievementReference()
	created, _ := repo.Create(ref)

	history1 := m.AchievementHistory{
		AchievementRefID: created.ID,
		PreviousStatus:   m.StatusDraft,
		NewStatus:        m.StatusSubmitted,
		ChangedByUserID:  uuid.New(),
	}
	history2 := m.AchievementHistory{
		AchievementRefID: created.ID,
		PreviousStatus:   m.StatusSubmitted,
		NewStatus:        m.StatusVerified,
		ChangedByUserID:  uuid.New(),
	}

	_ = repo.CreateHistory(history1)
	time.Sleep(10 * time.Millisecond)
	_ = repo.CreateHistory(history2)

	results, err := repo.FindHistoryByAchievementRefID(created.ID)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, m.StatusDraft, results[0].PreviousStatus)
	assert.Equal(t, m.StatusSubmitted, results[0].NewStatus)
	assert.Equal(t, m.StatusSubmitted, results[1].PreviousStatus)
	assert.Equal(t, m.StatusVerified, results[1].NewStatus)
}

func TestAchievementReference_FindHistoryByAchievementRefID_Empty(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	fakeID := uuid.New()
	results, err := repo.FindHistoryByAchievementRefID(fakeID)

	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestAchievementReference_FindHistoryByAchievementRefID_WithNote(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	ref := createTestAchievementReference()
	created, _ := repo.Create(ref)

	note := "Status changed by admin"
	history := m.AchievementHistory{
		AchievementRefID: created.ID,
		PreviousStatus:   m.StatusDraft,
		NewStatus:        m.StatusSubmitted,
		ChangedByUserID:  uuid.New(),
		Note:             &note,
	}

	_ = repo.CreateHistory(history)

	results, err := repo.FindHistoryByAchievementRefID(created.ID)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NotNil(t, results[0].Note)
	assert.Equal(t, note, *results[0].Note)
}

func TestAchievementReference_FullWorkflow(t *testing.T) {
	db := setupTestPostgresDB(t)
	defer cleanupTestPostgresDB(t, db)

	repo := NewAchievementReferenceRepository(db)

	ref := createTestAchievementReference()
	ref.Status = m.StatusDraft
	created, err := repo.Create(ref)
	assert.NoError(t, err)

	created.Status = m.StatusSubmitted
	now := time.Now()
	created.SubmittedAt = &now
	updated, err := repo.Update(created)
	assert.NoError(t, err)
	assert.Equal(t, m.StatusSubmitted, updated.Status)

	history1 := m.AchievementHistory{
		AchievementRefID: created.ID,
		PreviousStatus:   m.StatusDraft,
		NewStatus:        m.StatusSubmitted,
		ChangedByUserID:  created.StudentID,
	}
	err = repo.CreateHistory(history1)
	assert.NoError(t, err)

	verifierID := uuid.New()
	verified, err := repo.UpdateStatus(created.ID, verifierID, m.StatusVerified, nil)
	assert.NoError(t, err)
	assert.Equal(t, m.StatusVerified, verified.Status)

	history2 := m.AchievementHistory{
		AchievementRefID: created.ID,
		PreviousStatus:   m.StatusSubmitted,
		NewStatus:        m.StatusVerified,
		ChangedByUserID:  verifierID,
	}
	err = repo.CreateHistory(history2)
	assert.NoError(t, err)

	histories, err := repo.FindHistoryByAchievementRefID(created.ID)
	assert.NoError(t, err)
	assert.Len(t, histories, 2)
	assert.Equal(t, m.StatusDraft, histories[0].PreviousStatus)
	assert.Equal(t, m.StatusVerified, histories[1].NewStatus)

	studentRefs, err := repo.FindByStudentID(created.StudentID)
	assert.NoError(t, err)
	assert.Len(t, studentRefs, 1)
	assert.Equal(t, m.StatusVerified, studentRefs[0].Status)
}
