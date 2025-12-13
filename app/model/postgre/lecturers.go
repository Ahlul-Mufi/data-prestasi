package modelpostgre

import (
	"time"

	"github.com/google/uuid"
)

type Lecturer struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	LecturerID string    `json:"lecturer_id"`
	Department string    `json:"department"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateLecturerRequest struct {
	UserID     uuid.UUID `json:"user_id" validate:"required"`
	LecturerID string    `json:"lecturer_id" validate:"required"`
	Department string    `json:"department" validate:"required"`
}

type UpdateLecturerRequest struct {
	LecturerID string `json:"lecturer_id"`
	Department string `json:"department"`
}
