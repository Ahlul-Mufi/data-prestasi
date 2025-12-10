package repositorypostgre

import (
	"database/sql"

	"github.com/google/uuid"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
)

type PermissionRepository interface {
    GetAll() ([]m.Permission, error)
    GetByID(id uuid.UUID) (m.Permission, error)
    Create(m.Permission) (m.Permission, error)
    Update(id uuid.UUID, p m.Permission) (m.Permission, error)
    Delete(id uuid.UUID) error
}

type permissionRepository struct {
    db *sql.DB
}

func NewPermissionRepository(db *sql.DB) PermissionRepository {
    return &permissionRepository{db}
}

func (r *permissionRepository) GetAll() ([]m.Permission, error) {
    rows, err := r.db.Query(`SELECT id, name, resource, action, description FROM permissions`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var result []m.Permission
    for rows.Next() {
        var p m.Permission
        if err := rows.Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description); err != nil {
            return nil, err
        }
        result = append(result, p)
    }
    return result, nil
}

func (r *permissionRepository) GetByID(id uuid.UUID) (m.Permission, error) {
    var p m.Permission
    err := r.db.QueryRow(
        `SELECT id, name, resource, action, description FROM permissions WHERE id=$1`,
        id,
    ).Scan(&p.ID, &p.Name, &p.Resource, &p.Action, &p.Description)

    return p, err
}

func (r *permissionRepository) Create(p m.Permission) (m.Permission, error) {
    err := r.db.QueryRow(
        `INSERT INTO permissions (name, resource, action, description)
         VALUES ($1,$2,$3,$4)
         RETURNING id`,
        p.Name, p.Resource, p.Action, p.Description,
    ).Scan(&p.ID)

    return p, err
}

func (r *permissionRepository) Update(id uuid.UUID, p m.Permission) (m.Permission, error) {
    err := r.db.QueryRow(
        `UPDATE permissions
         SET name=$1, resource=$2, action=$3, description=$4
         WHERE id=$5
         RETURNING id`,
        p.Name, p.Resource, p.Action, p.Description, id,
    ).Scan(&p.ID)
    return p, err
}

func (r *permissionRepository) Delete(id uuid.UUID) error {
    _, err := r.db.Exec(`DELETE FROM permissions WHERE id=$1`, id)
    return err
}
