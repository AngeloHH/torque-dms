package repositories

import (
	"gorm.io/gorm"
	"torque-dms/core/inventory/domain"
	"torque-dms/core/inventory/ports/output"
	"torque-dms/models"
)

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) output.LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) Save(location *domain.Location) error {
	model := toLocationModel(location)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	location.ID = model.ID
	return nil
}

func (r *locationRepository) Update(location *domain.Location) error {
	model := toLocationModel(location)
	return r.db.Save(model).Error
}

func (r *locationRepository) FindByID(id uint) (*domain.Location, error) {
	var model models.Location
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainLocation(&model), nil
}

func (r *locationRepository) FindByName(name string) (*domain.Location, error) {
	var model models.Location
	result := r.db.Where("name = ?", name).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainLocation(&model), nil
}

func (r *locationRepository) FindAll() ([]*domain.Location, error) {
	var modelList []models.Location
	result := r.db.Order("name ASC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	locations := make([]*domain.Location, len(modelList))
	for i, model := range modelList {
		locations[i] = toDomainLocation(&model)
	}
	return locations, nil
}

func (r *locationRepository) FindByType(locationType domain.LocationType) ([]*domain.Location, error) {
	var modelList []models.Location
	result := r.db.Where("type = ?", locationType).Order("name ASC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	locations := make([]*domain.Location, len(modelList))
	for i, model := range modelList {
		locations[i] = toDomainLocation(&model)
	}
	return locations, nil
}

func (r *locationRepository) FindActive() ([]*domain.Location, error) {
	var modelList []models.Location
	result := r.db.Where("active = ?", true).Order("name ASC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	locations := make([]*domain.Location, len(modelList))
	for i, model := range modelList {
		locations[i] = toDomainLocation(&model)
	}
	return locations, nil
}

func (r *locationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Location{}, id).Error
}

func (r *locationRepository) Exists(id uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.Location{}).Where("id = ?", id).Count(&count)
	return count > 0, result.Error
}

// Mappers

func toLocationModel(l *domain.Location) *models.Location {
	return &models.Location{
		ID:        l.ID,
		Name:      l.Name,
		Type:      models.LocationType(l.Type),
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

func toDomainLocation(m *models.Location) *domain.Location {
	return &domain.Location{
		ID:        m.ID,
		Name:      m.Name,
		Type:      domain.LocationType(m.Type),
		Address:   m.Address,
		City:      m.City,
		State:     m.State,
		Zip:       m.Zip,
		CountryID: m.CountryID,
		Latitude:  m.Latitude,
		Longitude: m.Longitude,
		Capacity:  m.Capacity,
		Active:    m.Active,
		CreatedAt: m.CreatedAt,
	}
}