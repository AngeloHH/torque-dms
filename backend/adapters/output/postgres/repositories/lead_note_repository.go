package repositories

import (
	"gorm.io/gorm"
	"torque-dms/core/sales/domain"
	"torque-dms/core/sales/ports/output"
	"torque-dms/models"
)

type leadNoteRepository struct {
	db *gorm.DB
}

func NewLeadNoteRepository(db *gorm.DB) output.LeadNoteRepository {
	return &leadNoteRepository{db: db}
}

func (r *leadNoteRepository) Save(note *domain.LeadNote) error {
	model := toLeadNoteModel(note)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	note.ID = model.ID
	return nil
}

func (r *leadNoteRepository) Update(note *domain.LeadNote) error {
	model := toLeadNoteModel(note)
	return r.db.Save(model).Error
}

func (r *leadNoteRepository) FindByID(id uint) (*domain.LeadNote, error) {
	var model models.LeadNote
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainLeadNote(&model), nil
}

func (r *leadNoteRepository) FindByLeadID(leadID uint) ([]*domain.LeadNote, error) {
	var modelList []models.LeadNote
	result := r.db.Where("lead_id = ?", leadID).Order("created_at DESC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	notes := make([]*domain.LeadNote, len(modelList))
	for i, model := range modelList {
		notes[i] = toDomainLeadNote(&model)
	}
	return notes, nil
}

func (r *leadNoteRepository) Delete(id uint) error {
	return r.db.Delete(&models.LeadNote{}, id).Error
}

// Mappers

func toLeadNoteModel(n *domain.LeadNote) *models.LeadNote {
	return &models.LeadNote{
		ID:         n.ID,
		LeadID:     n.LeadID,
		Content:    n.Content,
		CreatedBy:  n.CreatedBy,
		CreatedAt:  n.CreatedAt,
		ModifiedAt: n.ModifiedAt,
	}
}

func toDomainLeadNote(m *models.LeadNote) *domain.LeadNote {
	return &domain.LeadNote{
		ID:         m.ID,
		LeadID:     m.LeadID,
		Content:    m.Content,
		CreatedBy:  m.CreatedBy,
		CreatedAt:  m.CreatedAt,
		ModifiedAt: m.ModifiedAt,
	}
}