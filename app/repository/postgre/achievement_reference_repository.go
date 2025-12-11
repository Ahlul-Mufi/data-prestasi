package repositorypostgre

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	"github.com/google/uuid"
)

type AchievementReferenceRepository interface {
	Create(a m.AchievementReference) (m.AchievementReference, error)
	GetByID(id uuid.UUID) (m.AchievementReference, error)
	Update(a m.AchievementReference) (m.AchievementReference, error)
	Delete(id uuid.UUID) error
	GetFiltered(userID *uuid.UUID, status *m.AchievementStatus) ([]m.AchievementReference, error)
	UpdateStatus(id, verifierID uuid.UUID, newStatus m.AchievementStatus, rejectionNote *string) (m.AchievementReference, error)
}

type achievementReferenceRepository struct {
	db *sql.DB
}

func NewAchievementReferenceRepository(db *sql.DB) AchievementReferenceRepository {
	return &achievementReferenceRepository{db}
}

func scanAchievementReference(row *sql.Row) (m.AchievementReference, error) {
	var ref m.AchievementReference
	var submittedAt sql.NullTime
	var verifiedAt sql.NullTime
	var verifiedByID sql.NullString
	var rejectionNote sql.NullString

	err := row.Scan(
		&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &submittedAt,
		&verifiedAt, &verifiedByID, &rejectionNote, &ref.CreatedAt, &ref.UpdatedAt,
	)

	if err != nil {
		return m.AchievementReference{}, err
	}

	if submittedAt.Valid {
		ref.SubmittedAt = &submittedAt.Time
	} else {
		ref.SubmittedAt = nil
	}
	if verifiedAt.Valid {
		ref.VerifiedAt = &verifiedAt.Time
	} else {
		ref.VerifiedAt = nil
	}
	if verifiedByID.Valid {
		id, _ := uuid.Parse(verifiedByID.String)
		ref.VerifiedBy = &id
	} else {
		ref.VerifiedBy = nil
	}
	if rejectionNote.Valid {
		ref.RejectionNote = &rejectionNote.String
	} else {
		ref.RejectionNote = nil
	}

	return ref, nil
}


func (r *achievementReferenceRepository) Create(a m.AchievementReference) (m.AchievementReference, error) {
	a.ID = uuid.New()
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	var submittedAtPtr *time.Time
	if a.Status == m.StatusSubmitted {
		submittedAtPtr = &now
	} else {
		submittedAtPtr = nil
	}

	query := `
        INSERT INTO achievement_references 
        (id, student_id, mongo_achievement_id, status, submitted_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
	`
	ref, err := scanAchievementReference(r.db.QueryRow(query,
		a.ID, a.StudentID, a.MongoAchievementID, a.Status, submittedAtPtr, a.CreatedAt, a.UpdatedAt,
	))

	if err != nil {
		return m.AchievementReference{}, err
	}
	return ref, nil
}

func (r *achievementReferenceRepository) GetByID(id uuid.UUID) (m.AchievementReference, error) {
	query := `
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references
        WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	ref, err := scanAchievementReference(row)

	if errors.Is(err, sql.ErrNoRows) {
		return m.AchievementReference{}, errors.New("achievement reference not found")
	}
	if err != nil {
		return m.AchievementReference{}, err
	}
	return ref, nil
}

func (r *achievementReferenceRepository) Update(a m.AchievementReference) (m.AchievementReference, error) {
	a.UpdatedAt = time.Now()

	var submittedAtDB interface{}
	if a.SubmittedAt != nil {
		submittedAtDB = *a.SubmittedAt
	}

	query := `
        UPDATE achievement_references 
        SET mongo_achievement_id = $1, status = $2, submitted_at = $3, updated_at = $4
        WHERE id = $5 
        RETURNING id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
	`
	row := r.db.QueryRow(query,
		a.MongoAchievementID, a.Status, submittedAtDB, a.UpdatedAt, a.ID,
	)
	
	ref, err := scanAchievementReference(row)
	
	if errors.Is(err, sql.ErrNoRows) {
		return m.AchievementReference{}, errors.New("achievement reference not found or cannot be updated")
	}
	if err != nil {
		return m.AchievementReference{}, err
	}
	return ref, nil
}

func (r *achievementReferenceRepository) Delete(id uuid.UUID) error {
	result, err := r.db.Exec(`DELETE FROM achievement_references WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("achievement reference not found or already deleted")
	}

	return nil
}

func (r *achievementReferenceRepository) GetFiltered(userID *uuid.UUID, status *m.AchievementStatus) ([]m.AchievementReference, error) {
	baseQuery := `
        SELECT id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
        FROM achievement_references
        WHERE 1=1
	`
	var args []interface{}
	argCount := 1

	if userID != nil {
		baseQuery += fmt.Sprintf(" AND student_id = $%d", argCount)
		args = append(args, userID)
		argCount++
	}

	if status != nil && *status != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	rows, err := r.db.Query(baseQuery + " ORDER BY created_at DESC", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []m.AchievementReference
	for rows.Next() {
		var ref m.AchievementReference
		var submittedAt, verifiedAt sql.NullTime
		var verifiedByID sql.NullString
		var rejectionNote sql.NullString

		if err := rows.Scan(
			&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status, &submittedAt,
			&verifiedAt, &verifiedByID, &rejectionNote, &ref.CreatedAt, &ref.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if submittedAt.Valid {
			ref.SubmittedAt = &submittedAt.Time
		}
		if verifiedAt.Valid {
			ref.VerifiedAt = &verifiedAt.Time
		}
		if verifiedByID.Valid {
			id, _ := uuid.Parse(verifiedByID.String)
			ref.VerifiedBy = &id
		}
		if rejectionNote.Valid {
			ref.RejectionNote = &rejectionNote.String
		}

		refs = append(refs, ref)
	}
	return refs, nil
}

func (r *achievementReferenceRepository) UpdateStatus(id, verifierID uuid.UUID, newStatus m.AchievementStatus, rejectionNote *string) (m.AchievementReference, error) {
	now := time.Now()

	var rejectionNoteDB interface{}
	if rejectionNote != nil {
		rejectionNoteDB = *rejectionNote
	} else {
		rejectionNoteDB = nil
	}

	var verifiedAtDB interface{} = nil
	if newStatus == m.StatusVerified || newStatus == m.StatusRejected {
		verifiedAtDB = now
	}

	query := `
        UPDATE achievement_references 
        SET status = $1, verified_by = $2, verified_at = $3, rejection_note = $4, updated_at = $5
        WHERE id = $6 AND status = $7
        RETURNING id, student_id, mongo_achievement_id, status, submitted_at, verified_at, verified_by, rejection_note, created_at, updated_at
	`
	row := r.db.QueryRow(query,
		newStatus, verifierID, verifiedAtDB, rejectionNoteDB, now, id, m.StatusSubmitted,
	)
	ref, err := scanAchievementReference(row)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m.AchievementReference{}, errors.New("achievement not found or already processed")
		}
		return m.AchievementReference{}, err
	}
	return ref, nil
}