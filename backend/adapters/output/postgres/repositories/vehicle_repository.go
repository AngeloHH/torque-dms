package repositories

import (
	"gorm.io/gorm"
	"torque-dms/core/inventory/domain"
	"torque-dms/core/inventory/ports/output"
	"torque-dms/models"
)

type vehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) output.VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) Save(vehicle *domain.Vehicle) error {
	model := toVehicleModel(vehicle)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	vehicle.ID = model.ID
	return nil
}

func (r *vehicleRepository) Update(vehicle *domain.Vehicle) error {
	model := toVehicleModel(vehicle)
	return r.db.Save(model).Error
}

func (r *vehicleRepository) FindByID(id uint) (*domain.Vehicle, error) {
	var model models.Vehicle
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainVehicle(&model), nil
}

func (r *vehicleRepository) FindByVIN(vin string) (*domain.Vehicle, error) {
	var model models.Vehicle
	result := r.db.Where("vin = ?", vin).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainVehicle(&model), nil
}

func (r *vehicleRepository) FindByStockNumber(stockNumber string) (*domain.Vehicle, error) {
	var model models.Vehicle
	result := r.db.Where("stock_number = ?", stockNumber).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainVehicle(&model), nil
}

func (r *vehicleRepository) FindAll(limit int, offset int) ([]*domain.Vehicle, error) {
	var modelList []models.Vehicle
	result := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	vehicles := make([]*domain.Vehicle, len(modelList))
	for i, model := range modelList {
		vehicles[i] = toDomainVehicle(&model)
	}
	return vehicles, nil
}

func (r *vehicleRepository) FindByStatus(status domain.VehicleStatus, limit int, offset int) ([]*domain.Vehicle, error) {
	var modelList []models.Vehicle
	result := r.db.Where("status = ?", status).Limit(limit).Offset(offset).Order("created_at DESC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	vehicles := make([]*domain.Vehicle, len(modelList))
	for i, model := range modelList {
		vehicles[i] = toDomainVehicle(&model)
	}
	return vehicles, nil
}

func (r *vehicleRepository) FindByLocationID(locationID uint) ([]*domain.Vehicle, error) {
	var modelList []models.Vehicle
	result := r.db.Where("location_id = ?", locationID).Order("created_at DESC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	vehicles := make([]*domain.Vehicle, len(modelList))
	for i, model := range modelList {
		vehicles[i] = toDomainVehicle(&model)
	}
	return vehicles, nil
}

func (r *vehicleRepository) FindAvailable(limit int, offset int) ([]*domain.Vehicle, error) {
	var modelList []models.Vehicle
	result := r.db.Where("status = ?", domain.VehicleStatusReadyForSale).Limit(limit).Offset(offset).Order("created_at DESC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	vehicles := make([]*domain.Vehicle, len(modelList))
	for i, model := range modelList {
		vehicles[i] = toDomainVehicle(&model)
	}
	return vehicles, nil
}

func (r *vehicleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Vehicle{}, id).Error
}

func (r *vehicleRepository) Exists(id uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.Vehicle{}).Where("id = ?", id).Count(&count)
	return count > 0, result.Error
}

func (r *vehicleRepository) ExistsByVIN(vin string) (bool, error) {
	var count int64
	result := r.db.Model(&models.Vehicle{}).Where("vin = ?", vin).Count(&count)
	return count > 0, result.Error
}

// Mappers

func toVehicleModel(v *domain.Vehicle) *models.Vehicle {
	return &models.Vehicle{
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
		Condition:         models.VehicleCondition(v.Condition),
		Status:            models.VehicleStatus(v.Status),
		LotType:           models.LotType(v.LotType),
		LocationID:        v.LocationID,
		AcquisitionSource: models.AcquisitionSource(v.AcquisitionSource),
		AcquisitionDate:   v.AcquisitionDate,
		AcquisitionCost:   v.AcquisitionCost,
		Model3DID:         v.Model3DID,
		CreatedAt:         v.CreatedAt,
		ModifiedAt:        v.ModifiedAt,
	}
}

func toDomainVehicle(m *models.Vehicle) *domain.Vehicle {
	return &domain.Vehicle{
		ID:                m.ID,
		StockNumber:       m.StockNumber,
		VIN:               m.VIN,
		Plate:             m.Plate,
		Make:              m.Make,
		Model:             m.Model,
		Trim:              m.Trim,
		Year:              m.Year,
		Mileage:           m.Mileage,
		ExteriorColor:     m.ExteriorColor,
		InteriorColor:     m.InteriorColor,
		MSRP:              m.MSRP,
		InvoicePrice:      m.InvoicePrice,
		AskingPrice:       m.AskingPrice,
		Condition:         domain.VehicleCondition(m.Condition),
		Status:            domain.VehicleStatus(m.Status),
		LotType:           domain.LotType(m.LotType),
		LocationID:        m.LocationID,
		AcquisitionSource: domain.AcquisitionSource(m.AcquisitionSource),
		AcquisitionDate:   m.AcquisitionDate,
		AcquisitionCost:   m.AcquisitionCost,
		Model3DID:         m.Model3DID,
		CreatedAt:         m.CreatedAt,
		ModifiedAt:        m.ModifiedAt,
	}
}