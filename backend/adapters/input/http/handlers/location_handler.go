package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"torque-dms/adapters/input/http/dto/request"
	"torque-dms/adapters/input/http/dto/response"
	"torque-dms/core/inventory/domain"
	"torque-dms/core/inventory/ports/input"
)

type LocationHandler struct {
	locationService input.LocationService
}

func NewLocationHandler(locationService input.LocationService) *LocationHandler {
	return &LocationHandler{locationService: locationService}
}

func (h *LocationHandler) Create(c *gin.Context) {
	var req request.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location, err := h.locationService.Create(input.CreateLocationInput{
		Name:      req.Name,
		Type:      req.Type,
		Address:   req.Address,
		City:      req.City,
		State:     req.State,
		Zip:       req.Zip,
		CountryID: req.CountryID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Capacity:  req.Capacity,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toLocationResponse(location))
}

func (h *LocationHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	location, err := h.locationService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
		return
	}

	c.JSON(http.StatusOK, toLocationResponse(location))
}

func (h *LocationHandler) List(c *gin.Context) {
	locations, err := h.locationService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.LocationResponse, len(locations))
	for i, location := range locations {
		responseList[i] = *toLocationResponse(location)
	}

	c.JSON(http.StatusOK, response.LocationListResponse{
		Locations: responseList,
		Total:     len(responseList),
	})
}

func (h *LocationHandler) ListActive(c *gin.Context) {
	locations, err := h.locationService.ListActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.LocationResponse, len(locations))
	for i, location := range locations {
		responseList[i] = *toLocationResponse(location)
	}

	c.JSON(http.StatusOK, response.LocationListResponse{
		Locations: responseList,
		Total:     len(responseList),
	})
}

func (h *LocationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req request.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location, err := h.locationService.Update(uint(id), input.UpdateLocationInput{
		Name:      req.Name,
		Address:   req.Address,
		City:      req.City,
		State:     req.State,
		Zip:       req.Zip,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Capacity:  req.Capacity,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toLocationResponse(location))
}

func (h *LocationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.locationService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "location deleted successfully"})
}

func (h *LocationHandler) Deactivate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.locationService.Deactivate(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "location deactivated"})
}

func (h *LocationHandler) Activate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.locationService.Activate(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "location activated"})
}

// Helper

func toLocationResponse(l *domain.Location) *response.LocationResponse {
	return &response.LocationResponse{
		ID:        l.ID,
		Name:      l.Name,
		Type:      string(l.Type),
		Address:   l.Address,
		City:      l.City,
		State:     l.State,
		Zip:       l.Zip,
		CountryID: l.CountryID,
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		Capacity:  l.Capacity,
		Active:    l.Active,
		CreatedAt: l.CreatedAt,
	}
}