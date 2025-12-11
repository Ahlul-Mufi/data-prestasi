package repositorypostgre

import (
	"database/sql"
	"fmt"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	"github.com/google/uuid"
)

type RolePermissionRepository interface {
	AddPermissionToRole(req m.RolePermission) (m.RolePermission, error)
	RemovePermissionFromRole(roleID, permissionID uuid.UUID) error
	GetPermissionsByRoleID(roleID uuid.UUID) ([]m.Permission, error)
}

type rolePermissionRepository struct {
	db *sql.DB
}

func NewRolePermissionRepository(db *sql.DB) RolePermissionRepository {
	return &rolePermissionRepository{db}
}

func (r *rolePermissionRepository) AddPermissionToRole(req m.RolePermission) (m.RolePermission, error) {
	query := `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) RETURNING role_id, permission_id`
	
	var existingRoleID uuid.UUID
	var existingPermID uuid.UUID
	checkQuery := `SELECT role_id, permission_id FROM role_permissions WHERE role_id=$1 AND permission_id=$2`
	err := r.db.QueryRow(checkQuery, req.RoleID, req.PermissionID).Scan(&existingRoleID, &existingPermID)
	
	if err == nil {
		return req, fmt.Errorf("permission already linked to role")
	}
	
	err = r.db.QueryRow(query, req.RoleID, req.PermissionID).Scan(&req.RoleID, &req.PermissionID)
	
	if err != nil {
		return m.RolePermission{}, err
	}
	return req, nil
}

func (r *rolePermissionRepository) RemovePermissionFromRole(roleID, permissionID uuid.UUID) error {
	result, err := r.db.Exec(`DELETE FROM role_permissions WHERE role_id=$1 AND permission_id=$2`, roleID, permissionID)
	if err != nil {
		return err
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}

func (r *rolePermissionRepository) GetPermissionsByRoleID(roleID uuid.UUID) ([]m.Permission, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.name, p.resource, p.action, p.description
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`, roleID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []m.Permission
	for rows.Next() {
		var p m.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	
	return perms, nil
}