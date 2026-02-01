package handlers

import (
	"net/http"
	"strconv"

	"torque-dms/adapters/input/http/dto/request"
	"torque-dms/adapters/input/http/dto/response"
	"torque-dms/core/identity/domain"
	"torque-dms/core/identity/ports/input"

	"github.com/gin-gonic/gin"
)

type EntityHandler struct {
	entityService input.EntityService
}

func NewEntityHandler(entityService input.EntityService) *EntityHandler {
	return &EntityHandler{entityService: entityService}
}

func (h *EntityHandler) Create(c *gin.Context) {
	var req request.CreateEntityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity, err := h.entityService.Create(input.CreateEntityInput{
		Type:         req.Type,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		BusinessName: req.BusinessName,
		TaxID:        req.TaxID,
		Email:        req.Email,
		Phone:        req.Phone,
		Address:      req.Address,
		City:         req.City,
		State:        req.State,
		Zip:          req.Zip,
		CountryID:    req.CountryID,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toEntityResponse(entity))
}

func (h *EntityHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	entity, err := h.entityService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "entity not found"})
		return
	}

	c.JSON(http.StatusOK, toEntityResponse(entity))
}

func (h *EntityHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	entities, err := h.entityService.List(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.EntityResponse, len(entities))
	for i, entity := range entities {
		responseList[i] = *toEntityResponse(entity)
	}

	c.JSON(http.StatusOK, response.EntityListResponse{
		Entities: responseList,
		Total:    len(responseList),
	})
}

func (h *EntityHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req request.UpdateEntityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity, err := h.entityService.Update(uint(id), input.UpdateEntityInput{
		Field: req.Field,
		Value: req.Value,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toEntityResponse(entity))
}

func (h *EntityHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.entityService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "entity deleted successfully"})
}

func (h *EntityHandler) Suspend(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.entityService.Suspend(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "entity suspended successfully"})
}

func (h *EntityHandler) Activate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.entityService.Activate(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "entity activated successfully"})
}

// Helper
func toEntityResponse(e *domain.Entity) *response.EntityResponse {
	if e == nil {
		return nil
	}
	return &response.EntityResponse{
		ID:             e.ID,
		Type:           string(e.Type),
		FirstName:      e.FirstName,
		LastName:       e.LastName,
		BusinessName:   e.BusinessName,
		TaxID:          e.TaxID,
		Email:          e.Email,
		Address:        e.Address,
		City:           e.City,
		State:          e.State,
		Zip:            e.Zip,
		CountryID:      e.CountryID,
		IsSystemUser:   e.IsSystemUser,
		IsInternal:     e.IsInternal,
		ParentEntityID: e.ParentEntityID,
		Status:         string(e.Status),
		CreatedAt:      e.CreatedAt,
		ModifiedAt:     e.ModifiedAt,
	}
}