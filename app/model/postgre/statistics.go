package modelpostgre

import "time"

type StatisticsResponse struct {
	TotalAchievements       int                         `json:"total_achievements"`
	AchievementsByType      map[string]int              `json:"achievements_by_type"`
	AchievementsByStatus    map[string]int              `json:"achievements_by_status"`
	AchievementsByPeriod    []PeriodStatistic           `json:"achievements_by_period"`
	TopStudents             []TopStudent                `json:"top_students"`
	CompetitionDistribution []CompetitionLevelStatistic `json:"competition_distribution"`
}

type PeriodStatistic struct {
	Period string `json:"period"`
	Count  int    `json:"count"`
}

type TopStudent struct {
	StudentID        string `json:"student_id"`
	StudentName      string `json:"student_name"`
	ProgramStudy     string `json:"program_study"`
	TotalPoints      int    `json:"total_points"`
	AchievementCount int    `json:"achievement_count"`
}

type CompetitionLevelStatistic struct {
	Level string `json:"level"`
	Count int    `json:"count"`
}

type StudentStatisticsResponse struct {
	StudentID            string               `json:"student_id"`
	StudentName          string               `json:"student_name"`
	ProgramStudy         string               `json:"program_study"`
	AcademicYear         string               `json:"academic_year"`
	TotalAchievements    int                  `json:"total_achievements"`
	TotalPoints          int                  `json:"total_points"`
	AchievementsByType   map[string]int       `json:"achievements_by_type"`
	AchievementsByStatus map[string]int       `json:"achievements_by_status"`
	RecentAchievements   []RecentAchievement  `json:"recent_achievements"`
	MonthlyTrend         []MonthlyAchievement `json:"monthly_trend"`
}

type RecentAchievement struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	AchievementType string    `json:"achievement_type"`
	Status          string    `json:"status"`
	Points          int       `json:"points"`
	CreatedAt       time.Time `json:"created_at"`
}

type MonthlyAchievement struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}
