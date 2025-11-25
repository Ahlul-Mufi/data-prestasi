package repositorypostgre

import (
	"database/sql"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	"github.com/google/uuid"
)

type RoleRepository interface {
    GetAll() ([]m.Role, error)
    GetByID(id uuid.UUID) (m.Role, error)
    Create(m.Role) (m.Role, error)
    Update(id uuid.UUID, r m.Role) (m.Role, error)
    Delete(id uuid.UUID) error
}

type roleRepository struct {
    db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
    return &roleRepository{db}
}

func (r *roleRepository) GetAll() ([]m.Role, error) {
    rows, err := r.db.Query(`SELECT id, name, description, created_at FROM roles`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var result []m.Role
    for rows.Next() {
        var role m.Role
        if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt); err != nil {
            return nil, err
        }
        result = append(result, role)
    }
    return result, nil
}

func (r *roleRepository) GetByID(id uuid.UUID) (m.Role, error) {
    var role m.Role
    err := r.db.QueryRow(
        `SELECT id, name, description, created_at FROM roles WHERE id=$1`,
        id,
    ).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)

    return role, err
}

func (r *roleRepository) Create(role m.Role) (m.Role, error) {
    err := r.db.QueryRow(
        `INSERT INTO roles (name, description) VALUES ($1, $2)
         RETURNING id, created_at`,
        role.Name, role.Description,
    ).Scan(&role.ID, &role.CreatedAt)

    return role, err
}

func (r *roleRepository) Update(id uuid.UUID, role m.Role) (m.Role, error) {
    err := r.db.QueryRow(
        `UPDATE roles SET name=$1, description=$2 WHERE id=$3
         RETURNING id, created_at`,
        role.Name, role.Description, id,
    ).Scan(&role.ID, &role.CreatedAt)

    return role, err
}

func (r *roleRepository) Delete(id uuid.UUID) error {
    _, err := r.db.Exec(`DELETE FROM roles WHERE id=$1`, id)
    return err
}
