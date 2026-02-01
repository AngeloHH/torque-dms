package repositories

import (
	"time"

	"gorm.io/gorm"
	"torque-dms/core/sales/domain"
	"torque-dms/core/sales/ports/output"
	"torque-dms/models"
)

type leadActivityRepository struct {
	db *gorm.DB
}

func NewLeadActivityRepository(db *gorm.DB) output.LeadActivityRepository {
	return &leadActivityRepository{db: db}
}

func (r *leadActivityRepository) Save(activity *domain.LeadActivity) error {
	model := toLeadActivityModel(activity)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	activity.ID = model.ID
	return nil
}

func (r *leadActivityRepository) Update(activity *domain.LeadActivity) error {
	model := toLeadActivityModel(activity)
	return r.db.Save(model).Error
}

func (r *leadActivityRepository) FindByID(id uint) (*domain.LeadActivity, error) {
	var model models.LeadActivity
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainLeadActivity(&model), nil
}

func (r *leadActivityRepository) FindByLeadID(leadID uint) ([]*domain.LeadActivity, error) {
	var modelList []models.LeadActivity
	result := r.db.Where("lead_id = ?", leadID).Order("created_at DESC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	activities := make([]*domain.LeadActivity, len(modelList))
	for i, model := range modelList {
		activities[i] = toDomainLeadActivity(&model)
	}
	return activities, nil
}

func (r *leadActivityRepository) FindScheduledByEntityID(entityID uint) ([]*domain.LeadActivity, error) {
	var modelList []models.LeadActivity
	result := r.db.Where("performed_by = ? AND scheduled_at IS NOT NULL AND completed_at IS NULL", entityID).
		Order("scheduled_at ASC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	activities := make([]*domain.LeadActivity, len(modelList))
	for i, model := range modelList {
		activities[i] = toDomainLeadActivity(&model)
	}
	return activities, nil
}

func (r *leadActivityRepository) FindOverdue() ([]*domain.LeadActivity, error) {
	var modelList []models.LeadActivity
	result := r.db.Where("scheduled_at < ? AND completed_at IS NULL", time.Now()).
		Order("scheduled_at ASC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	activities := make([]*domain.LeadActivity, len(modelList))
	for i, model := range modelList {
		activities[i] = toDomainLeadActivity(&model)
	}
	return activities, nil
}

func (r *leadActivityRepository) Delete(id uint) error {
	return r.db.Delete(&models.LeadActivity{}, id).Error
}

// Mappers

func toLeadActivityModel(a *domain.LeadActivity) *models.LeadActivity {
	return &models.LeadActivity{
		ID:          a.ID,
		LeadID:      a.LeadID,
		Type:        models.ActivityType(a.Type),
		Description: a.Description,
		Outcome:     a.Outcome,
		PhoneID:     a.PhoneID,
		Email:       a.Email,
		PerformedBy: a.PerformedBy,
		ScheduledAt: a.ScheduledAt,
		CompletedAt: a.CompletedAt,
		CreatedAt:   a.CreatedAt,
	}
}

func toDomainLeadActivity(m *models.LeadActivity) *domain.LeadActivity {
	return &domain.LeadActivity{
		ID:          m.ID,
		LeadID:      m.LeadID,
		Type:        domain.ActivityType(m.Type),
		Description: m.Description,
		Outcome:     m.Outcome,
		PhoneID:     m.PhoneID,
		Email:       m.Email,
		PerformedBy: m.PerformedBy,
		ScheduledAt: m.ScheduledAt,
		CompletedAt: m.CompletedAt,
		CreatedAt:   m.CreatedAt,
	}
}