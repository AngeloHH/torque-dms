package output

import "torque-dms/core/sales/domain"

type LeadStepPresetRepository interface {
	Save(preset *domain.LeadStepPreset) error
	Update(preset *domain.LeadStepPreset) error
	FindByID(id uint) (*domain.LeadStepPreset, error)
	FindByCode(code string) (*domain.LeadStepPreset, error)
	FindAll() ([]*domain.LeadStepPreset, error)
	FindPublic() ([]*domain.LeadStepPreset, error)
	FindByCreatedBy(entityID uint) ([]*domain.LeadStepPreset, error)
	Delete(id uint) error
	Exists(id uint) (bool, error)
}

type LeadStepRepository interface {
	Save(step *domain.LeadStep) error
	Update(step *domain.LeadStep) error
	FindByID(id uint) (*domain.LeadStep, error)
	FindByPresetID(presetID uint) ([]*domain.LeadStep, error)
	FindActiveByPresetID(presetID uint) ([]*domain.LeadStep, error)
	Delete(id uint) error
	Exists(id uint) (bool, error)
}

type LeadStepProgressRepository interface {
	Save(progress *domain.LeadStepProgress) error
	Update(progress *domain.LeadStepProgress) error
	FindByID(id uint) (*domain.LeadStepProgress, error)
	FindByLeadID(leadID uint) ([]*domain.LeadStepProgress, error)
	FindByLeadIDAndStepID(leadID uint, stepID uint) (*domain.LeadStepProgress, error)
	Delete(id uint) error
}