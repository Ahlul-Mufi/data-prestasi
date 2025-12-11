package repositorypostgre

import (
	"database/sql"
	"errors"
	"strings"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	"github.com/google/uuid"
)

type RolePermissionRepository interface {
	Add(rp m.RolePermission) (m.RolePermission, error)
	Remove(roleID, permissionID uuid.UUID) error
}

type rolePermissionRepository struct {
	db *sql.DB
}

func NewRolePermissionRepository(db *sql.DB) RolePermissionRepository {
	return &rolePermissionRepository{db}
}

func (r *rolePermissionRepository) Add(rp m.RolePermission) (m.RolePermission, error) {
	query := `
		INSERT INTO role_permissions (role_id, permission_id) 
		VALUES ($1, $2)
		ON CONFLICT (role_id, permission_id) DO NOTHING
	`
    _, err := r.db.Exec(query, rp.RoleID, rp.PermissionID)
    
    if err != nil && strings.Contains(err.Error(), "foreign key constraint") {
        return m.RolePermission{}, errors.New("role or permission not found")
    }
    if err != nil {
        return m.RolePermission{}, err
    }
    
	return rp, nil
}

func (r *rolePermissionRepository) Remove(roleID, permissionID uuid.UUID) error {
	result, err := r.db.Exec(`
		DELETE FROM role_permissions 
		WHERE role_id = $1 AND permission_id = $2
	`, roleID, permissionID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("role permission not found")
	}

	return nil
}