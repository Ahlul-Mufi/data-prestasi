package modelmongo

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementType string

const (
	TypeAcademic      AchievementType = "academic"
	TypeCompetition   AchievementType = "competition"
	TypeOrganization  AchievementType = "organization"
	TypePublication   AchievementType = "publication"
	TypeCertification AchievementType = "certification"
	TypeOther         AchievementType = "other"
)

type CompetitionLevel string

const (
	LevelInternational CompetitionLevel = "international"
	LevelNational      CompetitionLevel = "national"
	LevelRegional      CompetitionLevel = "regional"
	LevelLocal         CompetitionLevel = "local"
)

type PublicationType string

const (
	PubJournal    PublicationType = "journal"
	PubConference PublicationType = "conference"
	PubBook       PublicationType = "book"
)

type Achievement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID       string             `bson:"studentId" json:"student_id"`
	AchievementType AchievementType    `bson:"achievementType" json:"achievement_type"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	Details         AchievementDetails `bson:"details" json:"details"`
	Attachments     []Attachment       `bson:"attachments" json:"attachments"`
	Tags            []string           `bson:"tags" json:"tags"`
	Points          int                `bson:"points" json:"points"`
	IsDeleted       bool               `bson:"isDeleted" json:"is_deleted"`
	CreatedAt       time.Time          `bson:"createdAt" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updated_at"`
}

type AchievementDetails struct {
	CompetitionName  string            `bson:"competitionName,omitempty" json:"competition_name,omitempty"`
	CompetitionLevel *CompetitionLevel `bson:"competitionLevel,omitempty" json:"competition_level,omitempty"`
	Rank             *int              `bson:"rank,omitempty" json:"rank,omitempty"`
	MedalType        string            `bson:"medalType,omitempty" json:"medal_type,omitempty"`

	PublicationType  *PublicationType `bson:"publicationType,omitempty" json:"publication_type,omitempty"`
	PublicationTitle string           `bson:"publicationTitle,omitempty" json:"publication_title,omitempty"`
	Authors          []string         `bson:"authors,omitempty" json:"authors,omitempty"`
	Publisher        string           `bson:"publisher,omitempty" json:"publisher,omitempty"`
	ISSN             string           `bson:"issn,omitempty" json:"issn,omitempty"`

	OrganizationName string  `bson:"organizationName,omitempty" json:"organization_name,omitempty"`
	Position         string  `bson:"position,omitempty" json:"position,omitempty"`
	Period           *Period `bson:"period,omitempty" json:"period,omitempty"`

	CertificationName   string     `bson:"certificationName,omitempty" json:"certification_name,omitempty"`
	IssuedBy            string     `bson:"issuedBy,omitempty" json:"issued_by,omitempty"`
	CertificationNumber string     `bson:"certificationNumber,omitempty" json:"certification_number,omitempty"`
	ValidUntil          *time.Time `bson:"validUntil,omitempty" json:"valid_until,omitempty"`

	EventDate    *time.Time             `bson:"eventDate,omitempty" json:"event_date,omitempty"`
	Location     string                 `bson:"location,omitempty" json:"location,omitempty"`
	Organizer    string                 `bson:"organizer,omitempty" json:"organizer,omitempty"`
	Score        *float64               `bson:"score,omitempty" json:"score,omitempty"`
	CustomFields map[string]interface{} `bson:"customFields,omitempty" json:"custom_fields,omitempty"`
}

type Period struct {
	Start time.Time `bson:"start" json:"start"`
	End   time.Time `bson:"end" json:"end"`
}

type Attachment struct {
	FileName   string    `bson:"fileName" json:"file_name"`
	FileURL    string    `bson:"fileUrl" json:"file_url"`
	FileType   string    `bson:"fileType" json:"file_type"`
	FileSize   int64     `bson:"fileSize" json:"file_size"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploaded_at"`
}

type CreateAchievementRequest struct {
	AchievementType AchievementType    `json:"achievement_type" validate:"required"`
	Title           string             `json:"title" validate:"required"`
	Description     string             `json:"description" validate:"required"`
	Details         AchievementDetails `json:"details"`
	Tags            []string           `json:"tags"`
	Points          int                `json:"points"`
}

type UpdateAchievementRequest struct {
	AchievementType *AchievementType    `json:"achievement_type"`
	Title           *string             `json:"title"`
	Description     *string             `json:"description"`
	Details         *AchievementDetails `json:"details"`
	Tags            []string            `json:"tags"`
	Points          *int                `json:"points"`
}

type AchievementWithReference struct {
	Achievement   Achievement `json:"achievement"`
	Status        string      `json:"status"`
	SubmittedAt   *time.Time  `json:"submitted_at"`
	VerifiedAt    *time.Time  `json:"verified_at"`
	VerifiedBy    *uuid.UUID  `json:"verified_by"`
	RejectionNote *string     `json:"rejection_note"`
	ReferenceID   uuid.UUID   `json:"reference_id"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}
