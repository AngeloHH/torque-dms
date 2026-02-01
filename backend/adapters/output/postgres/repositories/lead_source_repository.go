package repositories

import (
	"gorm.io/gorm"
	"torque-dms/core/sales/domain"
	"torque-dms/core/sales/ports/output"
	"torque-dms/models"
)

type leadSourceRepository struct {
	db *gorm.DB
}

func NewLeadSourceRepository(db *gorm.DB) output.LeadSourceRepository {
	return &leadSourceRepository{db: db}
}

func (r *leadSourceRepository) Save(source *domain.LeadSource) error {
	model := toLeadSourceModel(source)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	source.ID = model.ID
	return nil
}

func (r *leadSourceRepository) Update(source *domain.LeadSource) error {
	model := toLeadSourceModel(source)
	return r.db.Save(model).Error
}

func (r *leadSourceRepository) FindByID(id uint) (*domain.LeadSource, error) {
	var model models.LeadSource
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainLeadSource(&model), nil
}

func (r *leadSourceRepository) FindByCode(code string) (*domain.LeadSource, error) {
	var model models.LeadSource
	result := r.db.Where("code = ?", code).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainLeadSource(&model), nil
}

func (r *leadSourceRepository) FindAll() ([]*domain.LeadSource, error) {
	var modelList []models.LeadSource
	result := r.db.Order("name ASC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	sources := make([]*domain.LeadSource, len(modelList))
	for i, model := range modelList {
		sources[i] = toDomainLeadSource(&model)
	}
	return sources, nil
}

func (r *leadSourceRepository) FindActive() ([]*domain.LeadSource, error) {
	var modelList []models.LeadSource
	result := r.db.Where("active = ?", true).Order("name ASC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	sources := make([]*domain.LeadSource, len(modelList))
	for i, model := range modelList {
		sources[i] = toDomainLeadSource(&model)
	}
	return sources, nil
}

func (r *leadSourceRepository) Delete(id uint) error {
	return r.db.Delete(&models.LeadSource{}, id).Error
}

func (r *leadSourceRepository) Exists(id uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.LeadSource{}).Where("id = ?", id).Count(&count)
	return count > 0, result.Error
}

// Mappers

func toLeadSourceModel(s *domain.LeadSource) *models.LeadSource {
	return &models.LeadSource{
		ID:         s.ID,
		Code:       s.Code,
		Name:       s.Name,
		IsExternal: s.IsExternal,
		Active:     s.Active,
		CreatedAt:  s.CreatedAt,
	}
}

func toDomainLeadSource(m *models.LeadSource) *domain.LeadSource {
	return &domain.LeadSource{
		ID:         m.ID,
		Code:       m.Code,
		Name:       m.Name,
		IsExternal: m.IsExternal,
		Active:     m.Active,
		CreatedAt:  m.CreatedAt,
	}
}