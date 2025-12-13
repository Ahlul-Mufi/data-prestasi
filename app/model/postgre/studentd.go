package modelpostgre

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	StudentID    string     `json:"student_id"`
	ProgramStudy string     `json:"program_study"`
	AcademicYear string     `json:"academic_year"`
	AdvisorID    *uuid.UUID `json:"advisor_id"`
	CreatedAt    time.Time  `json:"created_at"`
}

type CreateStudentRequest struct {
	UserID       uuid.UUID  `json:"user_id" validate:"required"`
	StudentID    string     `json:"student_id" validate:"required"`
	ProgramStudy string     `json:"program_study" validate:"required"`
	AcademicYear string     `json:"academic_year" validate:"required"`
	AdvisorID    *uuid.UUID `json:"advisor_id"`
}

type UpdateStudentRequest struct {
	StudentID    string     `json:"student_id"`
	ProgramStudy string     `json:"program_study"`
	AcademicYear string     `json:"academic_year"`
	AdvisorID    *uuid.UUID `json:"advisor_id"`
}

type UpdateAdvisorRequest struct {
	AdvisorID uuid.UUID `json:"advisor_id" validate:"required"`
}
