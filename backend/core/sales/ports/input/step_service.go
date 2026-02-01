package input

import "torque-dms/core/sales/domain"

type CreatePresetInput struct {
	Code        string
	Name        string
	Description string
	IsPublic    bool
	IsShared    bool
	CreatedBy   uint
}

type CreateStepInput struct {
	PresetID  uint
	Code      string
	Name      string
	SortOrder int
	IsFinal   bool
}

type UpdateProgressInput struct {
	LeadID      uint
	StepID      uint
	Status      string
	CompletedBy uint
	Notes       string
}

type StepService interface {
	// Presets
	CreatePreset(input CreatePresetInput) (*domain.LeadStepPreset, error)
	GetPreset(id uint) (*domain.LeadStepPreset, error)
	GetPresets() ([]*domain.LeadStepPreset, error)
	GetPublicPresets() ([]*domain.LeadStepPreset, error)
	GetMyPresets(entityID uint) ([]*domain.LeadStepPreset, error)
	DeletePreset(id uint) error
	MakePresetPublic(id uint) error
	MakePresetShared(id uint) error
	MakePresetPrivate(id uint) error

	// Steps
	CreateStep(input CreateStepInput) (*domain.LeadStep, error)
	GetSteps(presetID uint) ([]*domain.LeadStep, error)
	GetActiveSteps(presetID uint) ([]*domain.LeadStep, error)
	DeactivateStep(id uint) error
	ActivateStep(id uint) error
	DeleteStep(id uint) error

	// Progress
	InitializeProgress(leadID uint, presetID uint) error
	GetProgress(leadID uint) ([]*domain.LeadStepProgress, error)
	UpdateProgress(input UpdateProgressInput) (*domain.LeadStepProgress, error)
	CompleteStep(leadID uint, stepID uint, completedBy uint, notes string) error
	SkipStep(leadID uint, stepID uint, completedBy uint, notes string) error
	FailStep(leadID uint, stepID uint, completedBy uint, notes string) error
}