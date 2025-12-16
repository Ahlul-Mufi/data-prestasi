package servicepostgre

import (
	"errors"
	"strings"
	"time"

	modelmongo "github.com/Ahlul-Mufi/data-prestasi/app/model/mongo"
	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repomongo "github.com/Ahlul-Mufi/data-prestasi/app/repository/mongo"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	helper "github.com/Ahlul-Mufi/data-prestasi/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementReferenceService interface {
	Create(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetMyAchievements(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Submit(c *fiber.Ctx) error

	GetPendingAchievements(c *fiber.Ctx) error
	GetAdviseeAchievements(c *fiber.Ctx) error
	Verify(c *fiber.Ctx) error
	Reject(c *fiber.Ctx) error

	GetAllAchievements(c *fiber.Ctx) error
	GetHistory(c *fiber.Ctx) error
}

type achievementReferenceService struct {
	postgresRepo repo.AchievementReferenceRepository
	mongoRepo    repomongo.AchievementRepository
	studentRepo  repo.StudentRepository
	lecturerRepo repo.LecturerRepository
}

func NewAchievementReferenceService(
	postgresRepo repo.AchievementReferenceRepository,
	mongoRepo repomongo.AchievementRepository,
	studentRepo repo.StudentRepository,
	lecturerRepo repo.LecturerRepository,
) AchievementReferenceService {
	return &achievementReferenceService{
		postgresRepo: postgresRepo,
		mongoRepo:    mongoRepo,
		studentRepo:  studentRepo,
		lecturerRepo: lecturerRepo,
	}
}

func (s *achievementReferenceService) getCallerUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	return uuid.Parse(userIDStr)
}

func (s *achievementReferenceService) getStudentProfile(userID uuid.UUID) (m.Student, error) {
	student, err := s.studentRepo.GetByUserID(userID)
	if err != nil {
		return m.Student{}, errors.New("student profile not found")
	}
	return student, nil
}

func (s *achievementReferenceService) Create(c *fiber.Ctx) error {
	userID, err := s.getCallerUserID(c)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user ID in context", err.Error())
	}

	student, err := s.getStudentProfile(userID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Forbidden", "Only students can create achievements")
	}

	var req modelmongo.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if req.Title == "" || req.Description == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Validation error", "Title and description are required")
	}

	achievement := modelmongo.Achievement{
		StudentID:       student.ID.String(),
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		Points:          req.Points,
		Attachments:     []modelmongo.Attachment{},
	}

	createdAchievement, err := s.mongoRepo.Create(achievement)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create achievement in MongoDB", err.Error())
	}

	newReference := m.AchievementReference{
		StudentID:          student.ID,
		MongoAchievementID: createdAchievement.ID.Hex(),
		Status:             m.StatusDraft,
	}

	result, err := s.postgresRepo.Create(newReference)
	if err != nil {
		_ = s.mongoRepo.SoftDelete(createdAchievement.ID)
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create achievement reference", err.Error())
	}

	history := m.AchievementHistory{
		AchievementRefID: result.ID,
		PreviousStatus:   m.AchievementStatus(""),
		NewStatus:        m.StatusDraft,
		ChangedByUserID:  userID,
		Note:             nil,
	}
	_ = s.postgresRepo.CreateHistory(history)

	response := modelmongo.AchievementWithReference{
		Achievement: createdAchievement,
		Status:      string(result.Status),
		ReferenceID: result.ID,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, response)
}

func (s *achievementReferenceService) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	refID, err := uuid.Parse(idStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", err.Error())
	}

	result, err := s.postgresRepo.GetByID(refID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "Achievement not found.")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch achievement", err.Error())
	}

	mongoID, err := primitive.ObjectIDFromHex(result.MongoAchievementID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid MongoDB ID", err.Error())
	}

	achievement, err := s.mongoRepo.GetByID(mongoID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Achievement details not found in MongoDB", err.Error())
	}

	response := modelmongo.AchievementWithReference{
		Achievement:   achievement,
		Status:        string(result.Status),
		SubmittedAt:   result.SubmittedAt,
		VerifiedAt:    result.VerifiedAt,
		VerifiedBy:    result.VerifiedBy,
		RejectionNote: result.RejectionNote,
		ReferenceID:   result.ID,
		CreatedAt:     result.CreatedAt,
		UpdatedAt:     result.UpdatedAt,
	}

	return helper.SuccessResponse(c, fiber.StatusOK, response)
}

func (s *achievementReferenceService) GetMyAchievements(c *fiber.Ctx) error {
	userID, err := s.getCallerUserID(c)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user ID in context", err.Error())
	}

	student, err := s.getStudentProfile(userID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Forbidden", "Only students can view achievements")
	}

	refs, err := s.postgresRepo.FindByStudentID(student.ID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch achievements", err.Error())
	}

	var mongoIDs []primitive.ObjectID
	for _, ref := range refs {
		objID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
		if err == nil {
			mongoIDs = append(mongoIDs, objID)
		}
	}

	achievements, err := s.mongoRepo.GetMultipleByIDs(mongoIDs)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch achievement details", err.Error())
	}

	achievementMap := make(map[string]modelmongo.Achievement)
	for _, ach := range achievements {
		achievementMap[ach.ID.Hex()] = ach
	}

	var results []modelmongo.AchievementWithReference
	for _, ref := range refs {
		if ach, exists := achievementMap[ref.MongoAchievementID]; exists {
			results = append(results, modelmongo.AchievementWithReference{
				Achievement:   ach,
				Status:        string(ref.Status),
				SubmittedAt:   ref.SubmittedAt,
				VerifiedAt:    ref.VerifiedAt,
				VerifiedBy:    ref.VerifiedBy,
				RejectionNote: ref.RejectionNote,
				ReferenceID:   ref.ID,
				CreatedAt:     ref.CreatedAt,
				UpdatedAt:     ref.UpdatedAt,
			})
		}
	}

	return helper.SuccessResponse(c, fiber.StatusOK, results)
}

func (s *achievementReferenceService) Submit(c *fiber.Ctx) error {
	idStr := c.Params("id")
	refID, err := uuid.Parse(idStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", err.Error())
	}

	submitterID, err := s.getCallerUserID(c)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user ID in context", err.Error())
	}

	oldRef, err := s.postgresRepo.GetByID(refID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "Achievement not found.")
	}

	student, err := s.getStudentProfile(submitterID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Forbidden", "Only students can submit")
	}

	if oldRef.StudentID != student.ID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Forbidden", "You are not authorized to submit this achievement.")
	}

	if oldRef.Status != m.StatusDraft && oldRef.Status != m.StatusRejected {
		return helper.ErrorResponse(c, fiber.StatusConflict, "Submission Failed", "Achievement status must be 'draft' or 'rejected' to be submitted.")
	}

	oldStatus := oldRef.Status

	updatedRef := oldRef
	updatedRef.Status = m.StatusSubmitted
	now := time.Now()
	updatedRef.SubmittedAt = &now

	result, err := s.postgresRepo.Update(updatedRef)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update achievement status", err.Error())
	}

	history := m.AchievementHistory{
		AchievementRefID: result.ID,
		PreviousStatus:   oldStatus,
		NewStatus:        m.StatusSubmitted,
		ChangedByUserID:  submitterID,
		Note:             nil,
	}
	_ = s.postgresRepo.CreateHistory(history)

	return helper.SuccessResponse(c, fiber.StatusOK, result)
}

func (s *achievementReferenceService) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	refID, err := uuid.Parse(idStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", err.Error())
	}

	userID, err := s.getCallerUserID(c)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user ID in context", err.Error())
	}

	student, err := s.getStudentProfile(userID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Forbidden", "Only students can update")
	}

	oldRef, err := s.postgresRepo.GetByID(refID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "Achievement not found.")
	}

	if oldRef.StudentID != student.ID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Forbidden", "You are not authorized to update this achievement.")
	}

	if oldRef.Status != m.StatusDraft && oldRef.Status != m.StatusRejected {
		return helper.ErrorResponse(c, fiber.StatusConflict, "Update Failed", "Only draft or rejected achievements can be updated.")
	}

	var req modelmongo.UpdateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	mongoID, _ := primitive.ObjectIDFromHex(oldRef.MongoAchievementID)
	existing, err := s.mongoRepo.GetByID(mongoID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Achievement details not found", err.Error())
	}

	if req.AchievementType != nil {
		existing.AchievementType = *req.AchievementType
	}
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Details != nil {
		existing.Details = *req.Details
	}
	if req.Tags != nil {
		existing.Tags = req.Tags
	}
	if req.Points != nil {
		existing.Points = *req.Points
	}

	updated, err := s.mongoRepo.Update(mongoID, existing)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update achievement", err.Error())
	}

	if oldRef.Status == m.StatusRejected {
		oldRef.Status = m.StatusDraft
		oldRef.SubmittedAt = nil
		_, _ = s.postgresRepo.Update(oldRef)
	}

	response := modelmongo.AchievementWithReference{
		Achievement: updated,
		Status:      string(oldRef.Status),
		ReferenceID: oldRef.ID,
		CreatedAt:   oldRef.CreatedAt,
		UpdatedAt:   oldRef.UpdatedAt,
	}

	return helper.SuccessResponse(c, fiber.StatusOK, response)
}

func (s *achievementReferenceService) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	refID, err := uuid.Parse(idStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", err.Error())
	}

	userID, err := s.getCallerUserID(c)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user ID", err.Error())
	}

	student, err := s.getStudentProfile(userID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Forbidden", "Only students can delete")
	}

	oldRef, err := s.postgresRepo.GetByID(refID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "Achievement not found")
	}

	if oldRef.StudentID != student.ID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Forbidden", "You can only delete your own achievements")
	}

	if oldRef.Status != m.StatusDraft {
		return helper.ErrorResponse(c, fiber.StatusConflict, "Cannot delete", "Only draft achievements can be deleted")
	}

	mongoID, _ := primitive.ObjectIDFromHex(oldRef.MongoAchievementID)
	if err := s.mongoRepo.SoftDelete(mongoID); err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete from MongoDB", err.Error())
	}

	if err := s.postgresRepo.Delete(refID); err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete reference", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Achievement deleted successfully",
	})
}

func (s *achievementReferenceService) GetPendingAchievements(c *fiber.Ctx) error {
	status := m.StatusSubmitted
	refs, err := s.postgresRepo.GetFiltered(nil, &status)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch pending achievements", err.Error())
	}

	var mongoIDs []primitive.ObjectID
	for _, ref := range refs {
		objID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
		if err == nil {
			mongoIDs = append(mongoIDs, objID)
		}
	}

	achievements, err := s.mongoRepo.GetMultipleByIDs(mongoIDs)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch details", err.Error())
	}

	achievementMap := make(map[string]modelmongo.Achievement)
	for _, ach := range achievements {
		achievementMap[ach.ID.Hex()] = ach
	}

	var results []modelmongo.AchievementWithReference
	for _, ref := range refs {
		if ach, exists := achievementMap[ref.MongoAchievementID]; exists {
			results = append(results, modelmongo.AchievementWithReference{
				Achievement: ach,
				Status:      string(ref.Status),
				SubmittedAt: ref.SubmittedAt,
				ReferenceID: ref.ID,
				CreatedAt:   ref.CreatedAt,
				UpdatedAt:   ref.UpdatedAt,
			})
		}
	}

	return helper.SuccessResponse(c, fiber.StatusOK, results)
}

func (s *achievementReferenceService) GetAdviseeAchievements(c *fiber.Ctx) error {
	return helper.ErrorResponse(c, fiber.StatusNotImplemented, "Not implemented", "Under development")
}

func (s *achievementReferenceService) Verify(c *fiber.Ctx) error {
	idStr := c.Params("id")
	refID, err := uuid.Parse(idStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", err.Error())
	}

	verifierID, err := s.getCallerUserID(c)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid verifier ID in context", err.Error())
	}

	oldRef, err := s.postgresRepo.GetByID(refID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "Achievement not found.")
	}
	oldStatus := oldRef.Status

	result, err := s.postgresRepo.UpdateStatus(refID, verifierID, m.StatusVerified, nil)
	if err != nil {
		if errors.Is(err, errors.New("achievement not found or already processed")) {
			return helper.ErrorResponse(c, fiber.StatusConflict, "Verification Failed", "Achievement not found, or its status is not 'submitted'.")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to verify achievement", err.Error())
	}

	history := m.AchievementHistory{
		AchievementRefID: result.ID,
		PreviousStatus:   oldStatus,
		NewStatus:        m.StatusVerified,
		ChangedByUserID:  verifierID,
		Note:             nil,
	}
	_ = s.postgresRepo.CreateHistory(history)

	return helper.SuccessResponse(c, fiber.StatusOK, result)
}

func (s *achievementReferenceService) Reject(c *fiber.Ctx) error {
	idStr := c.Params("id")
	refID, err := uuid.Parse(idStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", err.Error())
	}

	verifierID, err := s.getCallerUserID(c)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid verifier ID in context", err.Error())
	}

	oldRef, err := s.postgresRepo.GetByID(refID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "Achievement not found.")
	}
	oldStatus := oldRef.Status

	var req m.VerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if req.RejectionNote == nil || *req.RejectionNote == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Missing Rejection Note", "Rejection Note must be provided when rejecting an achievement.")
	}

	result, err := s.postgresRepo.UpdateStatus(refID, verifierID, m.StatusRejected, req.RejectionNote)
	if err != nil {
		if errors.Is(err, errors.New("achievement not found or already processed")) {
			return helper.ErrorResponse(c, fiber.StatusConflict, "Rejection Failed", "Achievement not found, or its status is not 'submitted'.")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to reject achievement", err.Error())
	}

	history := m.AchievementHistory{
		AchievementRefID: result.ID,
		PreviousStatus:   oldStatus,
		NewStatus:        m.StatusRejected,
		ChangedByUserID:  verifierID,
		Note:             req.RejectionNote,
	}
	_ = s.postgresRepo.CreateHistory(history)

	return helper.SuccessResponse(c, fiber.StatusOK, result)
}

func (s *achievementReferenceService) GetAllAchievements(c *fiber.Ctx) error {
	statusQuery := c.Query("status")

	var statusFilter *m.AchievementStatus
	if statusQuery != "" {
		s := m.AchievementStatus(statusQuery)
		statusFilter = &s
	}

	refs, err := s.postgresRepo.GetFiltered(nil, statusFilter)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch achievements", err.Error())
	}

	var mongoIDs []primitive.ObjectID
	for _, ref := range refs {
		objID, err := primitive.ObjectIDFromHex(ref.MongoAchievementID)
		if err == nil {
			mongoIDs = append(mongoIDs, objID)
		}
	}

	achievements, err := s.mongoRepo.GetMultipleByIDs(mongoIDs)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch details", err.Error())
	}

	achievementMap := make(map[string]modelmongo.Achievement)
	for _, ach := range achievements {
		achievementMap[ach.ID.Hex()] = ach
	}

	var results []modelmongo.AchievementWithReference
	for _, ref := range refs {
		if ach, exists := achievementMap[ref.MongoAchievementID]; exists {
			results = append(results, modelmongo.AchievementWithReference{
				Achievement:   ach,
				Status:        string(ref.Status),
				SubmittedAt:   ref.SubmittedAt,
				VerifiedAt:    ref.VerifiedAt,
				VerifiedBy:    ref.VerifiedBy,
				RejectionNote: ref.RejectionNote,
				ReferenceID:   ref.ID,
				CreatedAt:     ref.CreatedAt,
				UpdatedAt:     ref.UpdatedAt,
			})
		}
	}

	return helper.SuccessResponse(c, fiber.StatusOK, results)
}

func (s *achievementReferenceService) GetHistory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	refID, err := uuid.Parse(idStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", err.Error())
	}

	_, err = s.postgresRepo.GetByID(refID)
	if err != nil {
		if errors.Is(err, errors.New("achievement reference not found")) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "Achievement not found.")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch achievement reference", err.Error())
	}

	histories, err := s.postgresRepo.FindHistoryByAchievementRefID(refID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch achievement history", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, histories)
}
