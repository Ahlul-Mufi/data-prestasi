package repositorypostgre

import (
	"database/sql"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	"github.com/google/uuid"
)

type StatisticsRepository interface {
	GetTotalAchievements(studentID *uuid.UUID) (int, error)
	GetAchievementsByType(studentID *uuid.UUID) (map[string]int, error)
	GetAchievementsByStatus(studentID *uuid.UUID) (map[string]int, error)
	GetTopStudents(limit int) ([]m.TopStudent, error)
	GetAchievementsByPeriod(studentID *uuid.UUID, periodType string) ([]m.PeriodStatistic, error)
	GetStudentInfo(studentID uuid.UUID) (m.Student, error)
	GetRecentAchievements(studentID uuid.UUID, limit int) ([]m.AchievementReference, error)
}

type statisticsRepository struct {
	db *sql.DB
}

func NewStatisticsRepository(db *sql.DB) StatisticsRepository {
	return &statisticsRepository{db}
}

func (r *statisticsRepository) GetTotalAchievements(studentID *uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM achievement_references WHERE 1=1`

	if studentID != nil {
		query += ` AND student_id = $1`
		err := r.db.QueryRow(query, studentID).Scan(&count)
		return count, err
	}

	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

func (r *statisticsRepository) GetAchievementsByType(studentID *uuid.UUID) (map[string]int, error) {
	result := make(map[string]int)

	query := `
		SELECT ar.mongo_achievement_id, ar.student_id 
		FROM achievement_references ar 
		WHERE 1=1
	`

	var rows *sql.Rows
	var err error

	if studentID != nil {
		query += ` AND ar.student_id = $1`
		rows, err = r.db.Query(query, studentID)
	} else {
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return result, err
	}
	defer rows.Close()

	result["academic"] = 0
	result["competition"] = 0
	result["organization"] = 0
	result["publication"] = 0
	result["certification"] = 0
	result["other"] = 0

	return result, nil
}

func (r *statisticsRepository) GetAchievementsByStatus(studentID *uuid.UUID) (map[string]int, error) {
	result := make(map[string]int)

	query := `
		SELECT status, COUNT(*) as count 
		FROM achievement_references 
		WHERE 1=1
	`

	var rows *sql.Rows
	var err error

	if studentID != nil {
		query += ` AND student_id = $1 GROUP BY status`
		rows, err = r.db.Query(query, studentID)
	} else {
		query += ` GROUP BY status`
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return result, err
		}
		result[status] = count
	}

	return result, nil
}

func (r *statisticsRepository) GetTopStudents(limit int) ([]m.TopStudent, error) {
	query := `
		SELECT 
			s.id,
			u.full_name,
			s.program_study,
			COUNT(ar.id) as achievement_count
		FROM students s
		JOIN users u ON u.id = s.user_id
		LEFT JOIN achievement_references ar ON ar.student_id = s.id
		WHERE ar.status = 'verified'
		GROUP BY s.id, u.full_name, s.program_study
		ORDER BY achievement_count DESC
		LIMIT $1
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topStudents []m.TopStudent
	for rows.Next() {
		var ts m.TopStudent
		if err := rows.Scan(&ts.StudentID, &ts.StudentName, &ts.ProgramStudy, &ts.AchievementCount); err != nil {
			return nil, err
		}
		ts.TotalPoints = 0
		topStudents = append(topStudents, ts)
	}

	return topStudents, nil
}

func (r *statisticsRepository) GetAchievementsByPeriod(studentID *uuid.UUID, periodType string) ([]m.PeriodStatistic, error) {
	var query string

	query = `
		SELECT 
			TO_CHAR(created_at, 'YYYY-MM') as period,
			COUNT(*) as count
		FROM achievement_references
		WHERE 1=1
	`

	if studentID != nil {
		query += ` AND student_id = $1`
	}

	query += ` GROUP BY period ORDER BY period DESC LIMIT 12`

	var rows *sql.Rows
	var err error

	if studentID != nil {
		rows, err = r.db.Query(query, studentID)
	} else {
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []m.PeriodStatistic
	for rows.Next() {
		var ps m.PeriodStatistic
		if err := rows.Scan(&ps.Period, &ps.Count); err != nil {
			return nil, err
		}
		stats = append(stats, ps)
	}

	return stats, nil
}

func (r *statisticsRepository) GetStudentInfo(studentID uuid.UUID) (m.Student, error) {
	var s m.Student
	var advisorID sql.NullString

	err := r.db.QueryRow(`
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students WHERE id=$1
	`, studentID).Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &advisorID, &s.CreatedAt)

	if err != nil {
		return m.Student{}, err
	}

	if advisorID.Valid {
		advUUID, _ := uuid.Parse(advisorID.String)
		s.AdvisorID = &advUUID
	}

	return s, nil
}

func (r *statisticsRepository) GetRecentAchievements(studentID uuid.UUID, limit int) ([]m.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status, submitted_at, 
		       verified_at, verified_by, rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, studentID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []m.AchievementReference
	for rows.Next() {
		var ref m.AchievementReference
		var submittedAt, verifiedAt sql.NullTime
		var verifiedBy sql.NullString
		var rejectionNote sql.NullString

		err := rows.Scan(
			&ref.ID, &ref.StudentID, &ref.MongoAchievementID, &ref.Status,
			&submittedAt, &verifiedAt, &verifiedBy, &rejectionNote,
			&ref.CreatedAt, &ref.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if submittedAt.Valid {
			ref.SubmittedAt = &submittedAt.Time
		}
		if verifiedAt.Valid {
			ref.VerifiedAt = &verifiedAt.Time
		}
		if verifiedBy.Valid {
			id, _ := uuid.Parse(verifiedBy.String)
			ref.VerifiedBy = &id
		}
		if rejectionNote.Valid {
			ref.RejectionNote = &rejectionNote.String
		}

		refs = append(refs, ref)
	}

	return refs, nil
}
