package repositories

import (
	"gorm.io/gorm"
	"torque-dms/core/inventory/domain"
	"torque-dms/core/inventory/ports/output"
	"torque-dms/models"
)

type vehiclePhotoRepository struct {
	db *gorm.DB
}

func NewVehiclePhotoRepository(db *gorm.DB) output.VehiclePhotoRepository {
	return &vehiclePhotoRepository{db: db}
}

func (r *vehiclePhotoRepository) Save(photo *domain.VehiclePhoto) error {
	model := toPhotoModel(photo)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	photo.ID = model.ID
	return nil
}

func (r *vehiclePhotoRepository) Update(photo *domain.VehiclePhoto) error {
	model := toPhotoModel(photo)
	return r.db.Save(model).Error
}

func (r *vehiclePhotoRepository) FindByID(id uint) (*domain.VehiclePhoto, error) {
	var model models.VehiclePhoto
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainPhoto(&model), nil
}

func (r *vehiclePhotoRepository) FindByVehicleID(vehicleID uint) ([]*domain.VehiclePhoto, error) {
	var modelList []models.VehiclePhoto
	result := r.db.Where("vehicle_id = ?", vehicleID).Order("sort_order ASC, created_at ASC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	photos := make([]*domain.VehiclePhoto, len(modelList))
	for i, model := range modelList {
		photos[i] = toDomainPhoto(&model)
	}
	return photos, nil
}

func (r *vehiclePhotoRepository) FindPrimaryByVehicleID(vehicleID uint) (*domain.VehiclePhoto, error) {
	var model models.VehiclePhoto
	result := r.db.Where("vehicle_id = ? AND is_primary = ?", vehicleID, true).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainPhoto(&model), nil
}

func (r *vehiclePhotoRepository) Delete(id uint) error {
	return r.db.Delete(&models.VehiclePhoto{}, id).Error
}

func (r *vehiclePhotoRepository) DeleteByVehicleID(vehicleID uint) error {
	return r.db.Where("vehicle_id = ?", vehicleID).Delete(&models.VehiclePhoto{}).Error
}

// Mappers

func toPhotoModel(p *domain.VehiclePhoto) *models.VehiclePhoto {
	return &models.VehiclePhoto{
		ID:          p.ID,
		VehicleID:   p.VehicleID,
		URL:         p.URL,
		Perspective: models.PhotoPerspective(p.Perspective),
		Purpose:     models.PhotoPurpose(p.Purpose),
		SortOrder:   p.SortOrder,
		IsPrimary:   p.IsPrimary,
		UploadedBy:  p.UploadedBy,
		CreatedAt:   p.CreatedAt,
	}
}

func toDomainPhoto(m *models.VehiclePhoto) *domain.VehiclePhoto {
	return &domain.VehiclePhoto{
		ID:          m.ID,
		VehicleID:   m.VehicleID,
		URL:         m.URL,
		Perspective: domain.PhotoPerspective(m.Perspective),
		Purpose:     domain.PhotoPurpose(m.Purpose),
		SortOrder:   m.SortOrder,
		IsPrimary:   m.IsPrimary,
		UploadedBy:  m.UploadedBy,
		CreatedAt:   m.CreatedAt,
	}
}