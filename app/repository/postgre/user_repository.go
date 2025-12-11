package repositorypostgre

import (
	"database/sql"
	"errors"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")
var ErrRoleNotFound = errors.New("role not found")
var ErrDuplicateUsernameOrEmail = errors.New("username or email already exists")

type UserRepository interface {
    FindByIdentity(identity string) (m.User, error)
    GetPermissions(roleID uuid.UUID) ([]string, error)
    FindByID(id uuid.UUID) (m.User, error)
    FindRoleIDByName(roleName string) (uuid.UUID, error)
    FindAll() ([]m.User, error)
    Create(user m.User, roleName string) (m.User, error)
    Update(user m.User, newRoleName string) (m.User, error)
    DeleteUser(id uuid.UUID) error
    UpdateRole(userID uuid.UUID, roleID uuid.UUID) error
}

type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepository{db}
}

func scanUser(row *sql.Row, u *m.User) error {
    return row.Scan(
        &u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName,
        &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
    )
}

func (r *userRepository) FindByIdentity(identity string) (m.User, error) {
    var u m.User
    row := r.db.QueryRow(`
        SELECT id, username, email, password_hash, full_name,
               role_id, is_active, created_at, updated_at
        FROM users
        WHERE username=$1 OR email=$1
    `, identity)
    
    err := scanUser(row, &u)
    if err == sql.ErrNoRows {
        return u, ErrUserNotFound
    }
    return u, err
}

func (r *userRepository) FindByID(id uuid.UUID) (m.User, error) {
    var u m.User
    row := r.db.QueryRow(`
        SELECT id, username, email, password_hash, full_name,
               role_id, is_active, created_at, updated_at
        FROM users
        WHERE id=$1
    `, id)
    
    err := scanUser(row, &u)
    if err == sql.ErrNoRows {
        return u, ErrUserNotFound
    }
    return u, err
}

func (r *userRepository) FindRoleIDByName(roleName string) (uuid.UUID, error) {
    var roleID uuid.UUID
    err := r.db.QueryRow("SELECT id FROM roles WHERE name = $1", roleName).Scan(&roleID)
    if err == sql.ErrNoRows {
        return uuid.Nil, ErrRoleNotFound
    }
    return roleID, err
}

func (r *userRepository) FindAll() ([]m.User, error) {
    rows, err := r.db.Query(`
        SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at
        FROM users
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []m.User
    for rows.Next() {
        var u m.User
        var roleID sql.NullString
        
        err := rows.Scan(
            &u.ID, &u.Username, &u.Email, &u.FullName,
            &roleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        
        if roleID.Valid {
            roleUUID, _ := uuid.Parse(roleID.String)
            u.RoleID = &roleUUID
        }
        
        users = append(users, u)
    }
    return users, nil
}

func (r *userRepository) Create(user m.User, roleName string) (m.User, error) {
    roleID, err := r.FindRoleIDByName(roleName)
    if err != nil {
        return user, ErrRoleNotFound
    }
    user.RoleID = &roleID
    
    err = r.db.QueryRow(`
        INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING created_at, updated_at
    `, user.ID, user.Username, user.Email, user.PasswordHash, user.FullName, user.RoleID, user.IsActive).Scan(
        &user.CreatedAt, &user.UpdatedAt,
    )
    
    if err != nil {
       if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` || 
          err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
           return user, ErrDuplicateUsernameOrEmail
       }
       return user, err
    }
    
    return user, nil
}

func (r *userRepository) Update(user m.User, newRoleName string) (m.User, error) {
    var roleID *uuid.UUID = user.RoleID
    
    if newRoleName != "" {
        newRoleID, err := r.FindRoleIDByName(newRoleName)
        if err != nil {
            return user, ErrRoleNotFound
        }
        roleID = &newRoleID
    }
    
    res, err := r.db.Exec(`
        UPDATE users
        SET username=$2, email=$3, password_hash=$4, 
            full_name=$5, role_id=$6, is_active=$7, updated_at=NOW()
        WHERE id=$1
    `, user.ID, user.Username, user.Email, user.PasswordHash, user.FullName, roleID, user.IsActive)
    
    if err != nil {
        return user, err
    }

    rowsAffected, _ := res.RowsAffected()
    if rowsAffected == 0 {
        return user, ErrUserNotFound
    }

    return r.FindByID(user.ID)
}

func (r *userRepository) DeleteUser(id uuid.UUID) error {
    res, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
    if err != nil {
        return err
    }

    rowsAffected, _ := res.RowsAffected()
    if rowsAffected == 0 {
        return ErrUserNotFound
    }
    return nil
}

func (r *userRepository) UpdateRole(userID uuid.UUID, roleID uuid.UUID) error {
    res, err := r.db.Exec(`
        UPDATE users
        SET role_id=$2, updated_at=NOW()
        WHERE id=$1
    `, userID, roleID)

    if err != nil {
        return err
    }

    rowsAffected, _ := res.RowsAffected()
    if rowsAffected == 0 {
        return ErrUserNotFound
    }
    return nil
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