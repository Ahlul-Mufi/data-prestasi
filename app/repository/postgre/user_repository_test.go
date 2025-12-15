package repositorypostgre

import (
	"database/sql"
	"testing"
	"time"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	return db, mock
}

func TestFindByIdentity_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	userID := uuid.New()
	roleID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "password_hash", "full_name",
		"role_id", "is_active", "created_at", "updated_at",
	}).AddRow(
		userID, "testuser", "test@example.com", "hashedpassword", "Test User",
		roleID, true, now, now,
	)

	mock.ExpectQuery(`SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at FROM users WHERE username=\$1 OR email=\$1`).
		WithArgs("testuser").
		WillReturnRows(rows)

	user, err := repo.FindByIdentity("testuser")

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByIdentity_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectQuery(`SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at FROM users WHERE username=\$1 OR email=\$1`).
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	user, err := repo.FindByIdentity("nonexistent")

	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
	assert.Equal(t, uuid.Nil, user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByID_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	userID := uuid.New()
	roleID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "password_hash", "full_name",
		"role_id", "is_active", "created_at", "updated_at",
	}).AddRow(
		userID, "testuser", "test@example.com", "hashedpassword", "Test User",
		roleID, true, now, now,
	)

	mock.ExpectQuery(`SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at FROM users WHERE id=\$1`).
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := repo.FindByID(userID)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	userID := uuid.New()

	mock.ExpectQuery(`SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at FROM users WHERE id=\$1`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.FindByID(userID)

	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
	assert.Equal(t, uuid.Nil, user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindRoleIDByName_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	roleID := uuid.New()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(roleID)

	mock.ExpectQuery(`SELECT id FROM roles WHERE name = \$1`).
		WithArgs("admin").
		WillReturnRows(rows)

	result, err := repo.FindRoleIDByName("admin")

	assert.NoError(t, err)
	assert.Equal(t, roleID, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindRoleIDByName_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectQuery(`SELECT id FROM roles WHERE name = \$1`).
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	result, err := repo.FindRoleIDByName("nonexistent")

	assert.Error(t, err)
	assert.Equal(t, ErrRoleNotFound, err)
	assert.Equal(t, uuid.Nil, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAll_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	userID1 := uuid.New()
	userID2 := uuid.New()
	roleID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "full_name", "role_id", "is_active", "created_at", "updated_at",
	}).AddRow(
		userID1, "user1", "user1@example.com", "User One", roleID.String(), true, now, now,
	).AddRow(
		userID2, "user2", "user2@example.com", "User Two", roleID.String(), true, now, now,
	)

	mock.ExpectQuery(`SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at FROM users`).
		WillReturnRows(rows)

	users, err := repo.FindAll()

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "user1", users[0].Username)
	assert.Equal(t, "user2", users[1].Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	userID := uuid.New()
	roleID := uuid.New()
	now := time.Now()

	user := m.User{
		ID:           userID,
		Username:     "newuser",
		Email:        "newuser@example.com",
		PasswordHash: "hashedpassword",
		FullName:     "New User",
		IsActive:     true,
	}

	// Mock role lookup
	roleRows := sqlmock.NewRows([]string{"id"}).AddRow(roleID)
	mock.ExpectQuery(`SELECT id FROM roles WHERE name = \$1`).
		WithArgs("user").
		WillReturnRows(roleRows)

	// Mock insert
	insertRows := sqlmock.NewRows([]string{"created_at", "updated_at"}).AddRow(now, now)
	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(userID, "newuser", "newuser@example.com", "hashedpassword", "New User", roleID, true).
		WillReturnRows(insertRows)

	createdUser, err := repo.Create(user, "user")

	assert.NoError(t, err)
	assert.Equal(t, userID, createdUser.ID)
	assert.Equal(t, "newuser", createdUser.Username)
	assert.NotNil(t, createdUser.RoleID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_RoleNotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	userID := uuid.New()

	user := m.User{
		ID:           userID,
		Username:     "newuser",
		Email:        "newuser@example.com",
		PasswordHash: "hashedpassword",
		FullName:     "New User",
		IsActive:     true,
	}

	mock.ExpectQuery(`SELECT id FROM roles WHERE name = \$1`).
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	createdUser, err := repo.Create(user, "nonexistent")

	assert.Error(t, err)
	assert.Equal(t, ErrRoleNotFound, err)
	assert.Equal(t, userID, createdUser.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	userID := uuid.New()
	roleID := uuid.New()
	now := time.Now()

	user := m.User{
		ID:           userID,
		Username:     "updateduser",
		Email:        "updated@example.com",
		PasswordHash: "newhash",
		FullName:     "Updated User",
		RoleID:       &roleID,
		IsActive:     true,
	}

	mock.ExpectExec(`UPDATE users SET username=\$2, email=\$3, password_hash=\$4, full_name=\$5, role_id=\$6, is_active=\$7, updated_at=NOW\(\) WHERE id=\$1`).
		WithArgs(userID, "updateduser", "updated@example.com", "newhash", "Updated User", roleID, true).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Mock FindByID for return value
	findRows := sqlmock.NewRows([]string{
		"id", "username", "email", "password_hash", "full_name",
		"role_id", "is_active", "created_at", "updated_at",
	}).AddRow(
		userID, "updateduser", "updated@example.com", "newhash", "Updated User",
		roleID, true, now, now,
	)
	mock.ExpectQuery(`SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at FROM users WHERE id=\$1`).
		WithArgs(userID).
		WillReturnRows(findRows)

	updatedUser, err := repo.Update(user, "")

	assert.NoError(t, err)
	assert.Equal(t, "updateduser", updatedUser.Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	userID := uuid.New()

	mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteUser(userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	userID := uuid.New()

	mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteUser(userID)

	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateRole_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	userID := uuid.New()
	roleID := uuid.New()

	mock.ExpectExec(`UPDATE users SET role_id=\$2, updated_at=NOW\(\) WHERE id=\$1`).
		WithArgs(userID, roleID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateRole(userID, roleID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPermissions_Success(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	roleID := uuid.New()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("user:read").
		AddRow("user:write").
		AddRow("user:delete")

	mock.ExpectQuery(`SELECT p.name FROM role_permissions rp JOIN permissions p ON p.id = rp.permission_id WHERE rp.role_id = \$1`).
		WithArgs(roleID).
		WillReturnRows(rows)

	permissions, err := repo.GetPermissions(roleID)

	assert.NoError(t, err)
	assert.Len(t, permissions, 3)
	assert.Contains(t, permissions, "user:read")
	assert.Contains(t, permissions, "user:write")
	assert.Contains(t, permissions, "user:delete")
	assert.NoError(t, mock.ExpectationsWereMet())
}
