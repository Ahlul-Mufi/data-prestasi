package repositorypostgre

import (
	"database/sql"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	"github.com/google/uuid"
)

type UserRepository interface {
    FindByIdentity(identity string) (m.User, error)
    GetPermissions(roleID uuid.UUID) ([]string, error)
    FindByID(id uuid.UUID) (m.User, error)
}

type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepository{db}
}

func (r *userRepository) FindByIdentity(identity string) (m.User, error) {
    var u m.User
    err := r.db.QueryRow(`
        SELECT id, username, email, password_hash, full_name,
               role_id, is_active, created_at, updated_at
        FROM users
        WHERE username=$1 OR email=$1
    `, identity).Scan(
        &u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName,
        &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
    )
    return u, err
}

func (r *userRepository) GetPermissions(roleID uuid.UUID) ([]string, error) {
    rows, err := r.db.Query(`
        SELECT p.name 
        FROM role_permissions rp
        JOIN permissions p ON p.id = rp.permission_id
        WHERE rp.role_id = $1
    `, roleID)

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var perms []string
    for rows.Next() {
        var pm string
        rows.Scan(&pm)
        perms = append(perms, pm)
    }

    return perms, nil
}

func (r *userRepository) FindByID(id uuid.UUID) (m.User, error) {
    var u m.User
    err := r.db.QueryRow(`
        SELECT id, username, email, password_hash, full_name,
               role_id, is_active, created_at, updated_at
        FROM users
        WHERE id=$1
    `, id).Scan(
        &u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName,
        &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
    )
    return u, err
}
