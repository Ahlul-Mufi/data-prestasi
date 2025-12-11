package servicepostgre

import (
	"database/sql"
	"errors"

	m "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	repo "github.com/Ahlul-Mufi/data-prestasi/app/repository/postgre"
	helper "github.com/Ahlul-Mufi/data-prestasi/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AchievementReferenceService interface {
	Create(c *fiber.Ctx) error
	GetMyAchievements(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error


	GetPendingAchievements(c *fiber.Ctx) error
	Verify(c *fiber.Ctx) error
	Reject(c *fiber.Ctx) error

	GetAllAchievements(c *fiber.Ctx) error
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
		Status:             m.StatusSubmitted,
	}

	result, err := s.repo.Create(newReference)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create achievement reference", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusCreated, result)
}

func (s *achievementReferenceService) GetMyAchievements(c *fiber.Ctx) error {
	userID, err := s.getCallerUserID(c)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user ID in context", err.Error())
	}

	statusFilter := c.Query("status")
	var statusPtr *m.AchievementStatus
	if statusFilter != "" {
		s := m.AchievementStatus(statusFilter)
		statusPtr = &s
	}

	results, err := s.repo.GetFiltered(&userID, statusPtr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch achievements", err.Error())
	}
	return helper.SuccessResponse(c, fiber.StatusOK, results)
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
		if errors.Is(err, sql.ErrNoRows) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Achievement not found", "The specified achievement reference does not exist.")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch data", err.Error())
	}

	if oldRef.StudentID != userID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Forbidden", "You are not authorized to update this achievement.")
	}

	if oldRef.Status == m.StatusVerified {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Forbidden", "Cannot update a verified achievement.")
	}
	
	updatedRef := oldRef
	if req.MongoAchievementID != nil {
		updatedRef.MongoAchievementID = *req.MongoAchievementID
	}
	
	if req.Status != "" {
		if req.Status == m.StatusSubmitted {
			updatedRef.Status = m.StatusSubmitted
		} else {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Status", "Only 'submitted' status change is allowed for student update.")
		}
	}


	result, err := s.repo.Update(updatedRef)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update achievement reference", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, result)
}

func (s *achievementReferenceService) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	refID, err := uuid.Parse(idStr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format", err.Error())
	}

	userID, err := s.getCallerUserID(c)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Invalid user ID in context", err.Error())
	}
	
	oldRef, err := s.repo.GetByID(refID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return helper.ErrorResponse(c, fiber.StatusNotFound, "Achievement not found", "The specified achievement reference does not exist.")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch data", err.Error())
	}

	if oldRef.StudentID != userID {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Forbidden", "You are not authorized to delete this achievement.")
	}

	if oldRef.Status == m.StatusVerified {
		return helper.ErrorResponse(c, fiber.StatusForbidden, "Forbidden", "Cannot delete a verified achievement.")
	}

	err = s.repo.Delete(refID)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete achievement reference", err.Error())
	}

	return helper.SuccessResponse(c, fiber.StatusOK, fiber.Map{"message": "Achievement reference deleted successfully"})
}

func (s *achievementReferenceService) GetPendingAchievements(c *fiber.Ctx) error {
	status := m.StatusSubmitted
	results, err := s.repo.GetFiltered(nil, &status) 
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch pending achievements", err.Error())
	}
	return helper.SuccessResponse(c, fiber.StatusOK, results)
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

	result, err := s.repo.UpdateStatus(refID, verifierID, m.StatusVerified, nil) 
	if err != nil {
		if errors.Is(err, errors.New("achievement not found or already processed")) {
			return helper.ErrorResponse(c, fiber.StatusConflict, "Verification Failed", "Achievement not found, or its status is not 'submitted'.")
		}
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to verify achievement", err.Error())
	}

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

	return helper.SuccessResponse(c, fiber.StatusOK, result)
}

func (s *achievementReferenceService) GetAllAchievements(c *fiber.Ctx) error {
	studentIDFilter := c.Query("student_id")
	var userIDPtr *uuid.UUID
	if studentIDFilter != "" {
		id, err := uuid.Parse(studentIDFilter)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid student_id format", err.Error())
		}
		userIDPtr = &id
	}
	
	statusFilter := c.Query("status")
	var statusPtr *m.AchievementStatus
	if statusFilter != "" {
		s := m.AchievementStatus(statusFilter)
		statusPtr = &s
	}
	
	results, err := s.repo.GetFiltered(userIDPtr, statusPtr)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch all achievements", err.Error())
	}
	return helper.SuccessResponse(c, fiber.StatusOK, results)
}