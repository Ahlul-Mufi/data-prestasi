package servicepostgre

import (
	"errors"
	"time"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	helper "github.com/Ahlul-Mufi/data-prestasi/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AchievementReferenceService interface {
	Create(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetMyAchievements(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Submit(c *fiber.Ctx) error

	GetPendingAchievements(c *fiber.Ctx) error
	Verify(c *fiber.Ctx) error
	Reject(c *fiber.Ctx) error

	GetAllAchievements(c *fiber.Ctx) error
	GetHistory(c *fiber.Ctx) error
}

type achievementReferenceService struct {
	repo repo.AchievementReferenceRepository
}

func NewAchievementReferenceService(r repo.AchievementReferenceRepository) AchievementReferenceService {
	return &achievementReferenceService{r}
}

func (s *achievementReferenceService) getCallerUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	return uuid.Parse(userIDStr)
}

func (s *achievementReferenceService) Create(c *fiber.Ctx) error {
	userID, err := s.getCallerUserID(c)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user ID in context", err.Error())
	}

	var req m.CreateAchievementRefRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if req.MongoAchievementID == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", "MongoAchievementID is required")
	}

	newReference := m.AchievementReference{
		StudentID:          userID,
		MongoAchievementID: req.MongoAchievementID,
		Status:             m.StatusDraft,
	}

	result, err := s.repo.Create(newReference)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create achievement reference", err.Error())
	}

	history := m.AchievementHistory{
		AchievementRefID: result.ID,
		PreviousStatus:   m.AchievementStatus(""),
		NewStatus:        m.StatusDraft,
		ChangedByUserID:  userID,
		Note:             nil,
	}
	_ = s.repo.CreateHistory(history)

	return helper.SuccessResponse(c, fiber.StatusCreated, result)
}

func (s *achievementReferenceService) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	refID, err := uuid.Parse(idStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", err.Error())
	}

	result, err := s.repo.GetByID(refID)
	if err != nil {
		if errors.Is(err, errors.New("achievement reference not found")) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "Achievement not found.")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch achievement", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, result)
}

func (s *achievementReferenceService) GetMyAchievements(c *fiber.Ctx) error {
	userID, err := s.getCallerUserID(c)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user ID in context", err.Error())
	}

	refs, err := s.repo.FindByStudentID(userID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch achievements", err.Error())
	}
	return helper.SuccessResponse(c, fiber.StatusOK, refs)
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

	oldRef, err := s.repo.GetByID(refID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "Achievement not found.")
	}

	if oldRef.StudentID != submitterID {
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

	result, err := s.repo.Update(updatedRef)
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
	_ = s.repo.CreateHistory(history)

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

	var req m.UpdateAchievementRefRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	oldRef, err := s.repo.GetByID(refID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "Achievement not found.")
	}

	if oldRef.StudentID != userID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Forbidden", "You are not authorized to update this achievement.")
	}

	if oldRef.Status == m.StatusVerified {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Forbidden", "Cannot update a verified achievement.")
	}

	oldStatus := oldRef.Status
	updatedRef := oldRef

	if req.MongoAchievementID != nil && *req.MongoAchievementID != "" {
		updatedRef.MongoAchievementID = *req.MongoAchievementID
		if oldRef.Status != m.StatusDraft {
			updatedRef.Status = m.StatusDraft
			updatedRef.SubmittedAt = nil
		}
	}

	if req.Status != "" {
		if req.Status == m.StatusDraft {
			updatedRef.Status = m.StatusDraft
			updatedRef.SubmittedAt = nil
		} else if req.Status != oldStatus {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Status", "Only 'draft' status change is allowed for student update (or use /submit).")
		}
	}

	result, err := s.repo.Update(updatedRef)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update achievement reference", err.Error())
	}

	if oldStatus != updatedRef.Status {
		history := m.AchievementHistory{
			AchievementRefID: result.ID,
			PreviousStatus:   oldStatus,
			NewStatus:        updatedRef.Status,
			ChangedByUserID:  userID,
			Note:             nil,
		}
		_ = s.repo.CreateHistory(history)
	}

	return helper.SuccessResponse(c, fiber.StatusOK, result)
}

func (s *achievementReferenceService) Delete(c *fiber.Ctx) error {
	return errors.New("not implemented yet")
}

func (s *achievementReferenceService) GetPendingAchievements(c *fiber.Ctx) error {
	return errors.New("not implemented yet")
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

	oldRef, err := s.repo.GetByID(refID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "Achievement not found.")
	}
	oldStatus := oldRef.Status

	result, err := s.repo.UpdateStatus(refID, verifierID, m.StatusVerified, nil)
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
	_ = s.repo.CreateHistory(history)

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

	oldRef, err := s.repo.GetByID(refID)
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

	result, err := s.repo.UpdateStatus(refID, verifierID, m.StatusRejected, req.RejectionNote)
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
	_ = s.repo.CreateHistory(history)

	return helper.SuccessResponse(c, fiber.StatusOK, result)
}

func (s *achievementReferenceService) GetAllAchievements(c *fiber.Ctx) error {
	return errors.New("not implemented yet")
}

func (s *achievementReferenceService) GetHistory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	refID, err := uuid.Parse(idStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", err.Error())
	}

	_, err = s.repo.GetByID(refID)
	if err != nil {
		if errors.Is(err, errors.New("achievement reference not found")) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Not Found", "Achievement not found.")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch achievement reference", err.Error())
	}

	histories, err := s.repo.FindHistoryByAchievementRefID(refID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch achievement history", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, histories)
}
