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

type VehicleHandler struct {
	vehicleService input.VehicleService
}

func NewVehicleHandler(vehicleService input.VehicleService) *VehicleHandler {
	return &VehicleHandler{vehicleService: vehicleService}
}

func (h *VehicleHandler) Create(c *gin.Context) {
	var req request.CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vehicle, err := h.vehicleService.Create(input.CreateVehicleInput{
		StockNumber:       req.StockNumber,
		VIN:               req.VIN,
		Plate:             req.Plate,
		Make:              req.Make,
		Model:             req.Model,
		Trim:              req.Trim,
		Year:              req.Year,
		Mileage:           req.Mileage,
		ExteriorColor:     req.ExteriorColor,
		InteriorColor:     req.InteriorColor,
		MSRP:              req.MSRP,
		InvoicePrice:      req.InvoicePrice,
		AskingPrice:       req.AskingPrice,
		Condition:         req.Condition,
		LocationID:        req.LocationID,
		AcquisitionSource: req.AcquisitionSource,
		AcquisitionCost:   req.AcquisitionCost,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toVehicleResponse(vehicle))
}

func (h *VehicleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	vehicle, err := h.vehicleService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
		return
	}

	c.JSON(http.StatusOK, toVehicleResponse(vehicle))
}

func (h *VehicleHandler) GetByVIN(c *gin.Context) {
	vin := c.Param("vin")

	vehicle, err := h.vehicleService.GetByVIN(vin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
		return
	}

	c.JSON(http.StatusOK, toVehicleResponse(vehicle))
}

func (h *VehicleHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	vehicles, err := h.vehicleService.List(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.VehicleResponse, len(vehicles))
	for i, vehicle := range vehicles {
		responseList[i] = *toVehicleResponse(vehicle)
	}

	c.JSON(http.StatusOK, response.VehicleListResponse{
		Vehicles: responseList,
		Total:    len(responseList),
	})
}

func (h *VehicleHandler) ListAvailable(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	vehicles, err := h.vehicleService.ListAvailable(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.VehicleResponse, len(vehicles))
	for i, vehicle := range vehicles {
		responseList[i] = *toVehicleResponse(vehicle)
	}

	c.JSON(http.StatusOK, response.VehicleListResponse{
		Vehicles: responseList,
		Total:    len(responseList),
	})
}

func (h *VehicleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req request.UpdateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vehicle, err := h.vehicleService.Update(uint(id), input.UpdateVehicleInput{
		Plate:         req.Plate,
		Trim:          req.Trim,
		Mileage:       req.Mileage,
		ExteriorColor: req.ExteriorColor,
		InteriorColor: req.InteriorColor,
		MSRP:          req.MSRP,
		InvoicePrice:  req.InvoicePrice,
		AskingPrice:   req.AskingPrice,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toVehicleResponse(vehicle))
}

func (h *VehicleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.vehicleService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vehicle deleted successfully"})
}

func (h *VehicleHandler) MarkAsSold(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.vehicleService.MarkAsSold(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vehicle marked as sold"})
}

func (h *VehicleHandler) MarkAsReadyForSale(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.vehicleService.MarkAsReadyForSale(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vehicle marked as ready for sale"})
}

func (h *VehicleHandler) SendToRecon(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.vehicleService.SendToRecon(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vehicle sent to recon"})
}

func (h *VehicleHandler) ChangeLocation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req request.ChangeLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.vehicleService.ChangeLocation(uint(id), req.LocationID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vehicle location changed"})
}

// Photos

func (h *VehicleHandler) AddPhoto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req request.AddPhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadedBy, _ := c.Get("entity_id")

	photo, err := h.vehicleService.AddPhoto(input.AddPhotoInput{
		VehicleID:   uint(id),
		URL:         req.URL,
		Perspective: req.Perspective,
		Purpose:     req.Purpose,
		UploadedBy:  uploadedBy.(uint),
		IsPrimary:   req.IsPrimary,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toPhotoResponse(photo))
}

func (h *VehicleHandler) GetPhotos(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	photos, err := h.vehicleService.GetPhotos(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.VehiclePhotoResponse, len(photos))
	for i, photo := range photos {
		responseList[i] = *toPhotoResponse(photo)
	}

	c.JSON(http.StatusOK, response.VehiclePhotosResponse{
		Photos: responseList,
		Total:  len(responseList),
	})
}

func (h *VehicleHandler) SetPrimaryPhoto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req request.SetPrimaryPhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.vehicleService.SetPrimaryPhoto(uint(id), req.PhotoID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "primary photo set"})
}

func (h *VehicleHandler) DeletePhoto(c *gin.Context) {
	photoID, err := strconv.ParseUint(c.Param("photoId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid photo id"})
		return
	}

	if err := h.vehicleService.DeletePhoto(uint(photoID)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "photo deleted"})
}

// Helpers

func toVehicleResponse(v *domain.Vehicle) *response.VehicleResponse {
	return &response.VehicleResponse{
		ID:                v.ID,
		StockNumber:       v.StockNumber,
		VIN:               v.VIN,
		Plate:             v.Plate,
		Make:              v.Make,
		Model:             v.Model,
		Trim:              v.Trim,
		Year:              v.Year,
		Mileage:           v.Mileage,
		ExteriorColor:     v.ExteriorColor,
		InteriorColor:     v.InteriorColor,
		MSRP:              v.MSRP,
		InvoicePrice:      v.InvoicePrice,
		AskingPrice:       v.AskingPrice,
		Condition:         string(v.Condition),
		Status:            string(v.Status),
		LotType:           string(v.LotType),
		LocationID:        v.LocationID,
		AcquisitionSource: string(v.AcquisitionSource),
		AcquisitionDate:   v.AcquisitionDate,
		AcquisitionCost:   v.AcquisitionCost,
		Profit:            v.Profit(),
		Margin:            v.Margin(),
		CreatedAt:         v.CreatedAt,
		ModifiedAt:        v.ModifiedAt,
	}
}

func toPhotoResponse(p *domain.VehiclePhoto) *response.VehiclePhotoResponse {
	return &response.VehiclePhotoResponse{
		ID:          p.ID,
		VehicleID:   p.VehicleID,
		URL:         p.URL,
		Perspective: string(p.Perspective),
		Purpose:     string(p.Purpose),
		SortOrder:   p.SortOrder,
		IsPrimary:   p.IsPrimary,
		CreatedAt:   p.CreatedAt,
	}
}