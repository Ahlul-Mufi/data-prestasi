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
	GetByUserID(userID uuid.UUID) (m.Student, error)
	Create(student m.Student) (m.Student, error)
	Update(student m.Student) (m.Student, error)
	Delete(id uuid.UUID) error
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
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at, updated_at 
		FROM students
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []m.Student
	for rows.Next() {
		var s m.Student
		var advisorID sql.NullString

		err := rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &advisorID, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if advisorID.Valid {
			advUUID, _ := uuid.Parse(advisorID.String)
			s.AdvisorID = &advUUID
		} else {
			s.AdvisorID = nil
		}

		students = append(students, s)
	}
	return students, nil
}

func (r *studentRepository) GetByID(id uuid.UUID) (m.Student, error) {
	var s m.Student
	var advisorID sql.NullString

	err := r.db.QueryRow(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at, updated_at
		FROM students WHERE id=$1
	`, id).Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &advisorID, &s.CreatedAt, &s.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return m.Student{}, errors.New("student not found")
	}
	if err != nil {
		return m.Student{}, err
	}

	if advisorID.Valid {
		advUUID, _ := uuid.Parse(advisorID.String)
		s.AdvisorID = &advUUID
	} else {
		s.AdvisorID = nil
	}

	return s, nil
}

func (r *studentRepository) GetByUserID(userID uuid.UUID) (m.Student, error) {
	var s m.Student
	var advisorID sql.NullString

	err := r.db.QueryRow(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at, updated_at
		FROM students WHERE user_id=$1
	`, userID).Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &advisorID, &s.CreatedAt, &s.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return m.Student{}, errors.New("student not found")
	}
	if err != nil {
		return m.Student{}, err
	}

	if advisorID.Valid {
		advUUID, _ := uuid.Parse(advisorID.String)
		s.AdvisorID = &advUUID
	}

	return s, nil
}

func (r *studentRepository) Create(student m.Student) (m.Student, error) {
	err := r.db.QueryRow(`
		INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`, student.ID, student.UserID, student.StudentID, student.ProgramStudy, student.AcademicYear, student.AdvisorID).Scan(
		&student.CreatedAt, &student.UpdatedAt,
	)

	if err != nil {
		return student, err
	}

	return student, nil
}

func (r *studentRepository) Update(student m.Student) (m.Student, error) {
	result, err := r.db.Exec(`
		UPDATE students 
		SET student_id=$2, program_study=$3, academic_year=$4, advisor_id=$5, updated_at=NOW()
		WHERE id=$1
	`, student.ID, student.StudentID, student.ProgramStudy, student.AcademicYear, student.AdvisorID)

	if err != nil {
		return student, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return student, err
	}

	if rowsAffected == 0 {
		return student, errors.New("student not found")
	}

	return r.GetByID(student.ID)
}

func (r *studentRepository) Delete(id uuid.UUID) error {
	result, err := r.db.Exec("DELETE FROM students WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("student not found")
	}
	return nil
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
		return errors.New("student not found or update failed")
	}
	return nil
}
