package repositorypostgre

import (
	"database/sql"
	"errors"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	"github.com/google/uuid"
)

type StudentRepository interface {
	GetAll() ([]m.Student, error)
	GetByID(id uuid.UUID) (m.Student, error)
	UpdateAdvisor(id, advisorID uuid.UUID) error
}

type studentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentRepository{db}
}

func (r *studentRepository) GetAll() ([]m.Student, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at 
		FROM students
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []m.Student
	for rows.Next() {
		var s m.Student
		err := rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}

func (r *studentRepository) GetByID(id uuid.UUID) (m.Student, error) {
	var s m.Student
	err := r.db.QueryRow(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at 
		FROM students WHERE id=$1
	`, id).Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return m.Student{}, errors.New("student not found")
	}
	return s, err
}

func (r *studentRepository) UpdateAdvisor(id, advisorID uuid.UUID) error {
	result, err := r.db.Exec(`
		UPDATE students SET advisor_id = $2, updated_at = NOW() 
		WHERE id = $1
	`, id, advisorID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("student not found or advisor already set")
	}
	return nil
}