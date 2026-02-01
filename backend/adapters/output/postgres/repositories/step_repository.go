package repositories

import (
	"gorm.io/gorm"
	"torque-dms/core/sales/domain"
	"torque-dms/core/sales/ports/output"
	"torque-dms/models"
)

// Preset Repository

type leadStepPresetRepository struct {
	db *gorm.DB
}

func NewLeadStepPresetRepository(db *gorm.DB) output.LeadStepPresetRepository {
	return &leadStepPresetRepository{db: db}
}

func (r *leadStepPresetRepository) Save(preset *domain.LeadStepPreset) error {
	model := toLeadStepPresetModel(preset)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	preset.ID = model.ID
	return nil
}

func (r *leadStepPresetRepository) Update(preset *domain.LeadStepPreset) error {
	model := toLeadStepPresetModel(preset)
	return r.db.Save(model).Error
}

func (r *leadStepPresetRepository) FindByID(id uint) (*domain.LeadStepPreset, error) {
	var model models.LeadStepPreset
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainLeadStepPreset(&model), nil
}

func (r *leadStepPresetRepository) FindByCode(code string) (*domain.LeadStepPreset, error) {
	var model models.LeadStepPreset
	result := r.db.Where("code = ?", code).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainLeadStepPreset(&model), nil
}

func (r *leadStepPresetRepository) FindAll() ([]*domain.LeadStepPreset, error) {
	var modelList []models.LeadStepPreset
	result := r.db.Order("sort_order ASC, name ASC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	presets := make([]*domain.LeadStepPreset, len(modelList))
	for i, model := range modelList {
		presets[i] = toDomainLeadStepPreset(&model)
	}
	return presets, nil
}

func (r *leadStepPresetRepository) FindPublic() ([]*domain.LeadStepPreset, error) {
	var modelList []models.LeadStepPreset
	result := r.db.Where("is_public = ?", true).Order("sort_order ASC, name ASC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	presets := make([]*domain.LeadStepPreset, len(modelList))
	for i, model := range modelList {
		presets[i] = toDomainLeadStepPreset(&model)
	}
	return presets, nil
}

func (r *leadStepPresetRepository) FindByCreatedBy(entityID uint) ([]*domain.LeadStepPreset, error) {
	var modelList []models.LeadStepPreset
	result := r.db.Where("created_by = ?", entityID).Order("sort_order ASC, name ASC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	presets := make([]*domain.LeadStepPreset, len(modelList))
	for i, model := range modelList {
		presets[i] = toDomainLeadStepPreset(&model)
	}
	return presets, nil
}

func (r *leadStepPresetRepository) Delete(id uint) error {
	return r.db.Delete(&models.LeadStepPreset{}, id).Error
}

func (r *leadStepPresetRepository) Exists(id uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.LeadStepPreset{}).Where("id = ?", id).Count(&count)
	return count > 0, result.Error
}

// Step Repository

type leadStepRepository struct {
	db *gorm.DB
}

func NewLeadStepRepository(db *gorm.DB) output.LeadStepRepository {
	return &leadStepRepository{db: db}
}

func (r *leadStepRepository) Save(step *domain.LeadStep) error {
	model := toLeadStepModel(step)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	step.ID = model.ID
	return nil
}

func (r *leadStepRepository) Update(step *domain.LeadStep) error {
	model := toLeadStepModel(step)
	return r.db.Save(model).Error
}

func (r *leadStepRepository) FindByID(id uint) (*domain.LeadStep, error) {
	var model models.LeadStep
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainLeadStep(&model), nil
}

func (r *leadStepRepository) FindByPresetID(presetID uint) ([]*domain.LeadStep, error) {
	var modelList []models.LeadStep
	result := r.db.Where("preset_id = ?", presetID).Order("sort_order ASC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	steps := make([]*domain.LeadStep, len(modelList))
	for i, model := range modelList {
		steps[i] = toDomainLeadStep(&model)
	}
	return steps, nil
}

func (r *leadStepRepository) FindActiveByPresetID(presetID uint) ([]*domain.LeadStep, error) {
	var modelList []models.LeadStep
	result := r.db.Where("preset_id = ? AND active = ?", presetID, true).Order("sort_order ASC").Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	steps := make([]*domain.LeadStep, len(modelList))
	for i, model := range modelList {
		steps[i] = toDomainLeadStep(&model)
	}
	return steps, nil
}

func (r *leadStepRepository) Delete(id uint) error {
	return r.db.Delete(&models.LeadStep{}, id).Error
}

func (r *leadStepRepository) Exists(id uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.LeadStep{}).Where("id = ?", id).Count(&count)
	return count > 0, result.Error
}

// Progress Repository

type leadStepProgressRepository struct {
	db *gorm.DB
}

func NewLeadStepProgressRepository(db *gorm.DB) output.LeadStepProgressRepository {
	return &leadStepProgressRepository{db: db}
}

func (r *leadStepProgressRepository) Save(progress *domain.LeadStepProgress) error {
	model := toLeadStepProgressModel(progress)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	progress.ID = model.ID
	return nil
}

func (r *leadStepProgressRepository) Update(progress *domain.LeadStepProgress) error {
	model := toLeadStepProgressModel(progress)
	return r.db.Save(model).Error
}

func (r *leadStepProgressRepository) FindByID(id uint) (*domain.LeadStepProgress, error) {
	var model models.LeadStepProgress
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainLeadStepProgress(&model), nil
}

func (r *leadStepProgressRepository) FindByLeadID(leadID uint) ([]*domain.LeadStepProgress, error) {
	var modelList []models.LeadStepProgress
	result := r.db.Where("lead_id = ?", leadID).Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	progressList := make([]*domain.LeadStepProgress, len(modelList))
	for i, model := range modelList {
		progressList[i] = toDomainLeadStepProgress(&model)
	}
	return progressList, nil
}

func (r *leadStepProgressRepository) FindByLeadIDAndStepID(leadID uint, stepID uint) (*domain.LeadStepProgress, error) {
	var model models.LeadStepProgress
	result := r.db.Where("lead_id = ? AND step_id = ?", leadID, stepID).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainLeadStepProgress(&model), nil
}

func (r *leadStepProgressRepository) Delete(id uint) error {
	return r.db.Delete(&models.LeadStepProgress{}, id).Error
}

// Mappers

func toLeadStepPresetModel(p *domain.LeadStepPreset) *models.LeadStepPreset {
	return &models.LeadStepPreset{
		ID:          p.ID,
		Code:        p.Code,
		Name:        p.Name,
		Description: p.Description,
		SortOrder:   p.SortOrder,
		IsPublic:    p.IsPublic,
		IsShared:    p.IsShared,
		CreatedBy:   p.CreatedBy,
		CreatedAt:   p.CreatedAt,
	}
}

func toDomainLeadStepPreset(m *models.LeadStepPreset) *domain.LeadStepPreset {
	return &domain.LeadStepPreset{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		SortOrder:   m.SortOrder,
		IsPublic:    m.IsPublic,
		IsShared:    m.IsShared,
		CreatedBy:   m.CreatedBy,
		CreatedAt:   m.CreatedAt,
	}
}

func toLeadStepModel(s *domain.LeadStep) *models.LeadStep {
	return &models.LeadStep{
		ID:        s.ID,
		PresetID:  s.PresetID,
		Code:      s.Code,
		Name:      s.Name,
		SortOrder: s.SortOrder,
		IsFinal:   s.IsFinal,
		Active:    s.Active,
		CreatedAt: s.CreatedAt,
	}
}

func toDomainLeadStep(m *models.LeadStep) *domain.LeadStep {
	return &domain.LeadStep{
		ID:        m.ID,
		PresetID:  m.PresetID,
		Code:      m.Code,
		Name:      m.Name,
		SortOrder: m.SortOrder,
		IsFinal:   m.IsFinal,
		Active:    m.Active,
		CreatedAt: m.CreatedAt,
	}
}

func toLeadStepProgressModel(p *domain.LeadStepProgress) *models.LeadStepProgress {
	return &models.LeadStepProgress{
		ID:          p.ID,
		LeadID:      p.LeadID,
		StepID:      p.StepID,
		Status:      models.StepStatus(p.Status),
		StartedAt:   p.StartedAt,
		CompletedAt: p.CompletedAt,
		CompletedBy: p.CompletedBy,
		Notes:       p.Notes,
		CreatedAt:   p.CreatedAt,
	}
}

func toDomainLeadStepProgress(m *models.LeadStepProgress) *domain.LeadStepProgress {
	return &domain.LeadStepProgress{
		ID:          m.ID,
		LeadID:      m.LeadID,
		StepID:      m.StepID,
		Status:      domain.StepStatus(m.Status),
		StartedAt:   m.StartedAt,
		CompletedAt: m.CompletedAt,
		CompletedBy: m.CompletedBy,
		Notes:       m.Notes,
		CreatedAt:   m.CreatedAt,
	}
}