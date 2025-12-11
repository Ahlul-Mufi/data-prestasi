package modelpostgre

import (
	"time"

	"github.com/google/uuid"
)

type AchievementStatus string

const (
    StatusDraft     AchievementStatus = "draft"
    StatusSubmitted AchievementStatus = "submitted"
    StatusVerified  AchievementStatus = "verified"
    StatusRejected  AchievementStatus = "rejected"
)

type AchievementReference struct {
    ID                 uuid.UUID         `json:"id"`
    StudentID          uuid.UUID         `json:"student_id"`
    MongoAchievementID string            `json:"mongo_achievement_id"` 
    Status             AchievementStatus `json:"status"`
    SubmittedAt        *time.Time        `json:"submitted_at"`
    VerifiedAt         *time.Time        `json:"verified_at"`
    VerifiedBy         *uuid.UUID        `json:"verified_by"`
    RejectionNote      *string           `json:"rejection_note"`
    CreatedAt          time.Time         `json:"created_at"`
    UpdatedAt          time.Time         `json:"updated_at"`
}

type CreateAchievementRefRequest struct {
    MongoAchievementID string `json:"mongo_achievement_id" validate:"required"`
}

type VerificationRequest struct {
    RejectionNote *string `json:"rejection_note"`
}

type UpdateAchievementRefRequest struct {
    MongoAchievementID *string           `json:"mongo_achievement_id"`
    Status             AchievementStatus `json:"status"` 
}