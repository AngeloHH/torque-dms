package services

import (
	"errors"

	"torque-dms/core/sales/domain"
	"torque-dms/core/sales/ports/input"
	"torque-dms/core/sales/ports/output"
)

type stepService struct {
	presetRepo   output.LeadStepPresetRepository
	stepRepo     output.LeadStepRepository
	progressRepo output.LeadStepProgressRepository
	leadRepo     output.LeadRepository
}

func NewStepService(
	presetRepo output.LeadStepPresetRepository,
	stepRepo output.LeadStepRepository,
	progressRepo output.LeadStepProgressRepository,
	leadRepo output.LeadRepository,
) input.StepService {
	return &stepService{
		presetRepo:   presetRepo,
		stepRepo:     stepRepo,
		progressRepo: progressRepo,
		leadRepo:     leadRepo,
	}
}

// Presets

func (s *stepService) CreatePreset(inp input.CreatePresetInput) (*domain.LeadStepPreset, error) {
	preset, err := domain.NewLeadStepPreset(inp.Code, inp.Name, inp.CreatedBy)
	if err != nil {
		return nil, err
	}

	preset.Description = inp.Description

	if inp.IsPublic {
		preset.MakePublic()
	} else if inp.IsShared {
		preset.MakeShared()
	}

	if err := s.presetRepo.Save(preset); err != nil {
		return nil, err
	}

	return preset, nil
}

func (s *stepService) GetPreset(id uint) (*domain.LeadStepPreset, error) {
	return s.presetRepo.FindByID(id)
}

func (s *stepService) GetPresets() ([]*domain.LeadStepPreset, error) {
	return s.presetRepo.FindAll()
}

func (s *stepService) GetPublicPresets() ([]*domain.LeadStepPreset, error) {
	return s.presetRepo.FindPublic()
}

func (s *stepService) GetMyPresets(entityID uint) ([]*domain.LeadStepPreset, error) {
	return s.presetRepo.FindByCreatedBy(entityID)
}

func (s *stepService) DeletePreset(id uint) error {
	exists, err := s.presetRepo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("preset not found")
	}

	return s.presetRepo.Delete(id)
}

func (s *stepService) MakePresetPublic(id uint) error {
	preset, err := s.presetRepo.FindByID(id)
	if err != nil {
		return errors.New("preset not found")
	}

	preset.MakePublic()
	return s.presetRepo.Update(preset)
}

func (s *stepService) MakePresetShared(id uint) error {
	preset, err := s.presetRepo.FindByID(id)
	if err != nil {
		return errors.New("preset not found")
	}

	preset.MakeShared()
	return s.presetRepo.Update(preset)
}

func (s *stepService) MakePresetPrivate(id uint) error {
	preset, err := s.presetRepo.FindByID(id)
	if err != nil {
		return errors.New("preset not found")
	}

	preset.MakePrivate()
	return s.presetRepo.Update(preset)
}

// Steps

func (s *stepService) CreateStep(inp input.CreateStepInput) (*domain.LeadStep, error) {
	exists, err := s.presetRepo.Exists(inp.PresetID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("preset not found")
	}

	step, err := domain.NewLeadStep(inp.PresetID, inp.Code, inp.Name, inp.SortOrder)
	if err != nil {
		return nil, err
	}

	if inp.IsFinal {
		step.MarkAsFinal()
	}

	if err := s.stepRepo.Save(step); err != nil {
		return nil, err
	}

	return step, nil
}

func (s *stepService) GetSteps(presetID uint) ([]*domain.LeadStep, error) {
	return s.stepRepo.FindByPresetID(presetID)
}

func (s *stepService) GetActiveSteps(presetID uint) ([]*domain.LeadStep, error) {
	return s.stepRepo.FindActiveByPresetID(presetID)
}

func (s *stepService) DeactivateStep(id uint) error {
	step, err := s.stepRepo.FindByID(id)
	if err != nil {
		return errors.New("step not found")
	}

	step.Deactivate()
	return s.stepRepo.Update(step)
}

func (s *stepService) ActivateStep(id uint) error {
	step, err := s.stepRepo.FindByID(id)
	if err != nil {
		return errors.New("step not found")
	}

	step.Activate()
	return s.stepRepo.Update(step)
}

func (s *stepService) DeleteStep(id uint) error {
	exists, err := s.stepRepo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("step not found")
	}

	return s.stepRepo.Delete(id)
}

// Progress

func (s *stepService) InitializeProgress(leadID uint, presetID uint) error {
	// Verificar que lead exista
	exists, err := s.leadRepo.Exists(leadID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("lead not found")
	}

	// Obtener pasos activos del preset
	steps, err := s.stepRepo.FindActiveByPresetID(presetID)
	if err != nil {
		return err
	}

	// Crear progreso para cada paso
	for _, step := range steps {
		progress, err := domain.NewLeadStepProgress(leadID, step.ID)
		if err != nil {
			return err
		}

		if err := s.progressRepo.Save(progress); err != nil {
			return err
		}
	}

	return nil
}

func (s *stepService) GetProgress(leadID uint) ([]*domain.LeadStepProgress, error) {
	return s.progressRepo.FindByLeadID(leadID)
}

func (s *stepService) UpdateProgress(inp input.UpdateProgressInput) (*domain.LeadStepProgress, error) {
	progress, err := s.progressRepo.FindByLeadIDAndStepID(inp.LeadID, inp.StepID)
	if err != nil {
		return nil, errors.New("progress not found")
	}

	switch domain.StepStatus(inp.Status) {
	case domain.StepStatusCompleted:
		progress.Complete(inp.CompletedBy, inp.Notes)
	case domain.StepStatusSkipped:
		progress.Skip(inp.CompletedBy, inp.Notes)
	case domain.StepStatusFailed:
		progress.Fail(inp.CompletedBy, inp.Notes)
	default:
		return nil, errors.New("invalid status")
	}

	if err := s.progressRepo.Update(progress); err != nil {
		return nil, err
	}

	return progress, nil
}

func (s *stepService) CompleteStep(leadID uint, stepID uint, completedBy uint, notes string) error {
	progress, err := s.progressRepo.FindByLeadIDAndStepID(leadID, stepID)
	if err != nil {
		return errors.New("progress not found")
	}

	progress.Complete(completedBy, notes)
	return s.progressRepo.Update(progress)
}

func (s *stepService) SkipStep(leadID uint, stepID uint, completedBy uint, notes string) error {
	progress, err := s.progressRepo.FindByLeadIDAndStepID(leadID, stepID)
	if err != nil {
		return errors.New("progress not found")
	}

	progress.Skip(completedBy, notes)
	return s.progressRepo.Update(progress)
}

func (s *stepService) FailStep(leadID uint, stepID uint, completedBy uint, notes string) error {
	progress, err := s.progressRepo.FindByLeadIDAndStepID(leadID, stepID)
	if err != nil {
		return errors.New("progress not found")
	}

	progress.Fail(completedBy, notes)
	return s.progressRepo.Update(progress)
}