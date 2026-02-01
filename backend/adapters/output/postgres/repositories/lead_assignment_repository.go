package repositories

import (
	"gorm.io/gorm"
	"torque-dms/core/sales/domain"
	"torque-dms/core/sales/ports/output"
	"torque-dms/models"
)

type leadAssignmentRepository struct {
	db *gorm.DB
}

func NewLeadAssignmentRepository(db *gorm.DB) output.LeadAssignmentRepository {
	return &leadAssignmentRepository{db: db}
}

func (r *leadAssignmentRepository) Save(assignment *domain.LeadAssignment) error {
	model := toLeadAssignmentModel(assignment)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	assignment.ID = model.ID
	return nil
}

func (r *leadAssignmentRepository) Update(assignment *domain.LeadAssignment) error {
	model := toLeadAssignmentModel(assignment)
	return r.db.Save(model).Error
}

func (r *leadAssignmentRepository) FindByID(id uint) (*domain.LeadAssignment, error) {
	var model models.LeadAssignment
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainLeadAssignment(&model), nil
}

func (r *leadAssignmentRepository) FindByLeadID(leadID uint) ([]*domain.LeadAssignment, error) {
	var modelList []models.LeadAssignment
	result := r.db.Where("lead_id = ? AND active = ?", leadID, true).Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	assignments := make([]*domain.LeadAssignment, len(modelList))
	for i, model := range modelList {
		assignments[i] = toDomainLeadAssignment(&model)
	}
	return assignments, nil
}

func (r *leadAssignmentRepository) FindPrimaryByLeadID(leadID uint) (*domain.LeadAssignment, error) {
	var model models.LeadAssignment
	result := r.db.Where("lead_id = ? AND is_primary = ? AND active = ?", leadID, true, true).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainLeadAssignment(&model), nil
}

func (r *leadAssignmentRepository) FindByEntityID(entityID uint) ([]*domain.LeadAssignment, error) {
	var modelList []models.LeadAssignment
	result := r.db.Where("entity_id = ? AND active = ?", entityID, true).Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	assignments := make([]*domain.LeadAssignment, len(modelList))
	for i, model := range modelList {
		assignments[i] = toDomainLeadAssignment(&model)
	}
	return assignments, nil
}

func (r *leadAssignmentRepository) Delete(id uint) error {
	return r.db.Delete(&models.LeadAssignment{}, id).Error
}

// Mappers

func toLeadAssignmentModel(a *domain.LeadAssignment) *models.LeadAssignment {
	return &models.LeadAssignment{
		ID:         a.ID,
		LeadID:     a.LeadID,
		EntityID:   a.EntityID,
		Role:       string(a.Role),
		IsPrimary:  a.IsPrimary,
		AssignedBy: a.AssignedBy,
		Active:     a.Active,
		CreatedAt:  a.CreatedAt,
	}
}

func toDomainLeadAssignment(m *models.LeadAssignment) *domain.LeadAssignment {
	return &domain.LeadAssignment{
		ID:         m.ID,
		LeadID:     m.LeadID,
		EntityID:   m.EntityID,
		Role:       domain.AssignmentRole(m.Role),
		IsPrimary:  m.IsPrimary,
		AssignedBy: m.AssignedBy,
		Active:     m.Active,
		CreatedAt:  m.CreatedAt,
	}
}