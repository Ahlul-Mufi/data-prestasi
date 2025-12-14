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
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	StudentID       string             `bson:"studentId"`
	AchievementType AchievementType    `bson:"achievementType"`
	Title           string             `bson:"title"`
	Description     string             `bson:"description"`
	Details         AchievementDetails `bson:"details"`
	Attachments     []Attachment       `bson:"attachments"`
	Tags            []string           `bson:"tags"`
	Points          int                `bson:"points"`
	IsDeleted       bool               `bson:"isDeleted"`
	CreatedAt       time.Time          `bson:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt"`
}

type AchievementDetails struct {
	CompetitionName  string            `bson:"competitionName,omitempty"`
	CompetitionLevel *CompetitionLevel `bson:"competitionLevel,omitempty"`
	Rank             *int              `bson:"rank,omitempty"`
	MedalType        string            `bson:"medalType,omitempty"`

	PublicationType  *PublicationType `bson:"publicationType,omitempty"`
	PublicationTitle string           `bson:"publicationTitle,omitempty"`
	Authors          []string         `bson:"authors,omitempty"`
	Publisher        string           `bson:"publisher,omitempty"`
	ISSN             string           `bson:"issn,omitempty"`

	OrganizationName string  `bson:"organizationName,omitempty"`
	Position         string  `bson:"position,omitempty"`
	Period           *Period `bson:"period,omitempty"`

	CertificationName   string     `bson:"certificationName,omitempty"`
	IssuedBy            string     `bson:"issuedBy,omitempty"`
	CertificationNumber string     `bson:"certificationNumber,omitempty"`
	ValidUntil          *time.Time `bson:"validUntil,omitempty"`

	EventDate    *time.Time             `bson:"eventDate,omitempty"`
	Location     string                 `bson:"location,omitempty"`
	Organizer    string                 `bson:"organizer,omitempty"`
	Score        *float64               `bson:"score,omitempty"`
	CustomFields map[string]interface{} `bson:"customFields,omitempty"`
}

type Period struct {
	Start time.Time `bson:"start"`
	End   time.Time `bson:"end"`
}

type Attachment struct {
	FileName   string    `bson:"fileName"`
	FileURL    string    `bson:"fileUrl"`
	FileType   string    `bson:"fileType"`
	FileSize   int64     `bson:"fileSize"`
	UploadedAt time.Time `bson:"uploadedAt"`
}

type CreateAchievementRequest struct {
	AchievementType AchievementType    `json:"achievement_type"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
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
