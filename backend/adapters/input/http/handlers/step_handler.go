package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"torque-dms/adapters/input/http/dto/request"
	"torque-dms/adapters/input/http/dto/response"
	"torque-dms/core/sales/domain"
	"torque-dms/core/sales/ports/input"
)

type StepHandler struct {
	stepService input.StepService
}

func NewStepHandler(stepService input.StepService) *StepHandler {
	return &StepHandler{stepService: stepService}
}

// Presets

func (h *StepHandler) CreatePreset(c *gin.Context) {
	var req request.CreatePresetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBy, _ := c.Get("entity_id")

	preset, err := h.stepService.CreatePreset(input.CreatePresetInput{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsPublic:    req.IsPublic,
		IsShared:    req.IsShared,
		CreatedBy:   createdBy.(uint),
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toLeadStepPresetResponse(preset))
}

func (h *StepHandler) GetPreset(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	preset, err := h.stepService.GetPreset(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "preset not found"})
		return
	}

	c.JSON(http.StatusOK, toLeadStepPresetResponse(preset))
}

func (h *StepHandler) GetPresets(c *gin.Context) {
	presets, err := h.stepService.GetPresets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.LeadStepPresetResponse, len(presets))
	for i, p := range presets {
		responseList[i] = *toLeadStepPresetResponse(p)
	}

	c.JSON(http.StatusOK, response.LeadStepPresetListResponse{
		Presets: responseList,
		Total:   len(responseList),
	})
}

func (h *StepHandler) GetPublicPresets(c *gin.Context) {
	presets, err := h.stepService.GetPublicPresets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.LeadStepPresetResponse, len(presets))
	for i, p := range presets {
		responseList[i] = *toLeadStepPresetResponse(p)
	}

	c.JSON(http.StatusOK, response.LeadStepPresetListResponse{
		Presets: responseList,
		Total:   len(responseList),
	})
}

func (h *StepHandler) GetMyPresets(c *gin.Context) {
	entityID, _ := c.Get("entity_id")

	presets, err := h.stepService.GetMyPresets(entityID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.LeadStepPresetResponse, len(presets))
	for i, p := range presets {
		responseList[i] = *toLeadStepPresetResponse(p)
	}

	c.JSON(http.StatusOK, response.LeadStepPresetListResponse{
		Presets: responseList,
		Total:   len(responseList),
	})
}

func (h *StepHandler) DeletePreset(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.stepService.DeletePreset(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "preset deleted"})
}

func (h *StepHandler) MakePresetPublic(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.stepService.MakePresetPublic(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "preset made public"})
}

func (h *StepHandler) MakePresetShared(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.stepService.MakePresetShared(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "preset made shared"})
}

func (h *StepHandler) MakePresetPrivate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.stepService.MakePresetPrivate(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "preset made private"})
}

// Steps

func (h *StepHandler) CreateStep(c *gin.Context) {
	presetID, err := strconv.ParseUint(c.Param("presetId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid preset id"})
		return
	}

	var req request.CreateStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	step, err := h.stepService.CreateStep(input.CreateStepInput{
		PresetID:  uint(presetID),
		Code:      req.Code,
		Name:      req.Name,
		SortOrder: req.SortOrder,
		IsFinal:   req.IsFinal,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toLeadStepResponse(step))
}

func (h *StepHandler) GetSteps(c *gin.Context) {
	presetID, err := strconv.ParseUint(c.Param("presetId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid preset id"})
		return
	}

	steps, err := h.stepService.GetSteps(uint(presetID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.LeadStepResponse, len(steps))
	for i, s := range steps {
		responseList[i] = *toLeadStepResponse(s)
	}

	c.JSON(http.StatusOK, response.LeadStepListResponse{
		Steps: responseList,
		Total: len(responseList),
	})
}

func (h *StepHandler) DeactivateStep(c *gin.Context) {
	stepID, err := strconv.ParseUint(c.Param("stepId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step id"})
		return
	}

	if err := h.stepService.DeactivateStep(uint(stepID)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "step deactivated"})
}

func (h *StepHandler) ActivateStep(c *gin.Context) {
	stepID, err := strconv.ParseUint(c.Param("stepId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step id"})
		return
	}

	if err := h.stepService.ActivateStep(uint(stepID)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "step activated"})
}

func (h *StepHandler) DeleteStep(c *gin.Context) {
	stepID, err := strconv.ParseUint(c.Param("stepId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step id"})
		return
	}

	if err := h.stepService.DeleteStep(uint(stepID)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "step deleted"})
}

// Helpers

func toLeadStepPresetResponse(p *domain.LeadStepPreset) *response.LeadStepPresetResponse {
	return &response.LeadStepPresetResponse{
		ID:          p.ID,
		Code:        p.Code,
		Name:        p.Name,
		Description: p.Description,
		SortOrder:   p.SortOrder,
		IsPublic:    p.IsPublic,
		IsShared:    p.IsShared,
		CreatedBy:   p.CreatedBy,
		CreatedAt:   p.CreatedAt,
	}
}

func toLeadStepResponse(s *domain.LeadStep) *response.LeadStepResponse {
	return &response.LeadStepResponse{
		ID:        s.ID,
		PresetID:  s.PresetID,
		Code:      s.Code,
		Name:      s.Name,
		SortOrder: s.SortOrder,
		IsFinal:   s.IsFinal,
		Active:    s.Active,
		CreatedAt: s.CreatedAt,
	}
}