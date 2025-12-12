package repositorypostgre

import (
	"database/sql"
	"errors"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	"github.com/google/uuid"
)

type LecturerRepository interface {
	GetAll() ([]m.Lecturer, error)
	GetByID(id uuid.UUID) (m.Lecturer, error)
	GetAdvisees(lecturerID uuid.UUID) ([]m.Student, error)
}

type lecturerRepository struct {
	db *sql.DB
}

func NewLecturerRepository(db *sql.DB) LecturerRepository {
	return &lecturerRepository{db}
}

func (r *lecturerRepository) GetAll() ([]m.Lecturer, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, lecturer_id, department, created_at 
		FROM lecturers
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lecturers []m.Lecturer
	for rows.Next() {
		var l m.Lecturer
		err := rows.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt)
		if err != nil {
			return nil, err
		}
		lecturers = append(lecturers, l)
	}
	return lecturers, nil
}

func (r *lecturerRepository) GetByID(id uuid.UUID) (m.Lecturer, error) {
	var l m.Lecturer
	err := r.db.QueryRow(`
		SELECT id, user_id, lecturer_id, department, created_at 
		FROM lecturers WHERE id=$1
	`, id).Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt)
	
	if errors.Is(err, sql.ErrNoRows) {
		return m.Lecturer{}, errors.New("lecturer not found")
	}
	return l, err
}

func (r *lecturerRepository) GetAdvisees(lecturerID uuid.UUID) ([]m.Student, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at 
		FROM students 
		WHERE advisor_id = $1
	`, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var advisees []m.Student
	for rows.Next() {
		var s m.Student
		err := rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		advisees = append(advisees, s)
	}
	return advisees, nil
}