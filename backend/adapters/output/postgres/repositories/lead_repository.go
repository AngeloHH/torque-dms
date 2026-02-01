package repositories

import (
	"gorm.io/gorm"
	"torque-dms/core/sales/domain"
	"torque-dms/core/sales/ports/output"
	"torque-dms/models"
)

type leadRepository struct {
	db *gorm.DB
}

func NewLeadRepository(db *gorm.DB) output.LeadRepository {
	return &leadRepository{db: db}
}

func (r *leadRepository) Save(lead *domain.Lead) error {
	model := toLeadModel(lead)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	lead.ID = model.ID
	return nil
}

func (r *leadRepository) Update(lead *domain.Lead) error {
	model := toLeadModel(lead)
	return r.db.Save(model).Error
}

func (r *leadRepository) FindByID(id uint) (*domain.Lead, error) {
	var model models.Lead
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainLead(&model), nil
}

func (r *leadRepository) FindByEntityID(entityID uint) ([]*domain.Lead, error) {
	var modelList []models.Lead
	result := r.db.Where("entity_id = ?", entityID).Order("created_at DESC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	leads := make([]*domain.Lead, len(modelList))
	for i, model := range modelList {
		leads[i] = toDomainLead(&model)
	}
	return leads, nil
}

func (r *leadRepository) FindAll(limit int, offset int) ([]*domain.Lead, error) {
	var modelList []models.Lead
	result := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	leads := make([]*domain.Lead, len(modelList))
	for i, model := range modelList {
		leads[i] = toDomainLead(&model)
	}
	return leads, nil
}

func (r *leadRepository) Delete(id uint) error {
	return r.db.Delete(&models.Lead{}, id).Error
}

func (r *leadRepository) Exists(id uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.Lead{}).Where("id = ?", id).Count(&count)
	return count > 0, result.Error
}

// Mappers

func toLeadModel(l *domain.Lead) *models.Lead {
	return &models.Lead{
		ID:            l.ID,
		EntityID:      l.EntityID,
		VehicleID:     l.VehicleID,
		InterestType:  models.VehicleCondition(l.InterestType),
		InterestMake:  l.InterestMake,
		InterestModel: l.InterestModel,
		BudgetMin:     l.BudgetMin,
		BudgetMax:     l.BudgetMax,
		SourceID:      l.SourceID,
		SourceDetail:  l.SourceDetail,
		PresetID:      l.PresetID,
		CreatedAt:     l.CreatedAt,
		ModifiedAt:    l.ModifiedAt,
	}
}

func toDomainLead(m *models.Lead) *domain.Lead {
	return &domain.Lead{
		ID:            m.ID,
		EntityID:      m.EntityID,
		VehicleID:     m.VehicleID,
		InterestType:  string(m.InterestType),
		InterestMake:  m.InterestMake,
		InterestModel: m.InterestModel,
		BudgetMin:     m.BudgetMin,
		BudgetMax:     m.BudgetMax,
		SourceID:      m.SourceID,
		SourceDetail:  m.SourceDetail,
		PresetID:      m.PresetID,
		CreatedAt:     m.CreatedAt,
		ModifiedAt:    m.ModifiedAt,
	}
}