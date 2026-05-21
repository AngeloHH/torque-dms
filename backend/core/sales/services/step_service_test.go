package services

import (
	"errors"
	"testing"

	"torque-dms/core/sales/domain"
	"torque-dms/core/sales/ports/input"
)

// --- Mocks ---

type mockPresetRepo struct {
	saved    *domain.LeadStepPreset
	findByID *domain.LeadStepPreset
	exists   bool
	presets  []*domain.LeadStepPreset
}

func (m *mockPresetRepo) Save(p *domain.LeadStepPreset) error {
	p.ID = 1
	m.saved = p
	return nil
}
func (m *mockPresetRepo) Update(p *domain.LeadStepPreset) error { return nil }
func (m *mockPresetRepo) FindByID(id uint) (*domain.LeadStepPreset, error) {
	if m.findByID != nil {
		return m.findByID, nil
	}
	return nil, errors.New("not found")
}
func (m *mockPresetRepo) FindByCode(code string) (*domain.LeadStepPreset, error) { return nil, nil }
func (m *mockPresetRepo) FindAll() ([]*domain.LeadStepPreset, error)             { return m.presets, nil }
func (m *mockPresetRepo) FindPublic() ([]*domain.LeadStepPreset, error)          { return nil, nil }
func (m *mockPresetRepo) FindByCreatedBy(id uint) ([]*domain.LeadStepPreset, error) {
	return nil, nil
}
func (m *mockPresetRepo) Delete(id uint) error        { return nil }
func (m *mockPresetRepo) Exists(id uint) (bool, error) { return m.exists, nil }

type mockStepRepo struct {
	saved    *domain.LeadStep
	findByID *domain.LeadStep
	exists   bool
	steps    []*domain.LeadStep
}

func (m *mockStepRepo) Save(s *domain.LeadStep) error {
	s.ID = 1
	m.saved = s
	return nil
}
func (m *mockStepRepo) Update(s *domain.LeadStep) error { return nil }
func (m *mockStepRepo) FindByID(id uint) (*domain.LeadStep, error) {
	if m.findByID != nil {
		return m.findByID, nil
	}
	return nil, errors.New("not found")
}
func (m *mockStepRepo) FindByPresetID(id uint) ([]*domain.LeadStep, error) { return m.steps, nil }
func (m *mockStepRepo) FindActiveByPresetID(id uint) ([]*domain.LeadStep, error) {
	return m.steps, nil
}
func (m *mockStepRepo) Delete(id uint) error        { return nil }
func (m *mockStepRepo) Exists(id uint) (bool, error) { return m.exists, nil }

type mockProgressRepo struct {
	saved            *domain.LeadStepProgress
	findByLeadAndStep *domain.LeadStepProgress
	progress         []*domain.LeadStepProgress
}

func (m *mockProgressRepo) Save(p *domain.LeadStepProgress) error {
	p.ID = 1
	m.saved = p
	return nil
}
func (m *mockProgressRepo) Update(p *domain.LeadStepProgress) error { return nil }
func (m *mockProgressRepo) FindByID(id uint) (*domain.LeadStepProgress, error) { return nil, nil }
func (m *mockProgressRepo) FindByLeadID(id uint) ([]*domain.LeadStepProgress, error) {
	return m.progress, nil
}
func (m *mockProgressRepo) FindByLeadIDAndStepID(leadID uint, stepID uint) (*domain.LeadStepProgress, error) {
	if m.findByLeadAndStep != nil {
		return m.findByLeadAndStep, nil
	}
	return nil, errors.New("not found")
}
func (m *mockProgressRepo) Delete(id uint) error { return nil }

// --- Preset Tests ---

func TestStepService_CreatePreset(t *testing.T) {
	pr := &mockPresetRepo{}
	sr := &mockStepRepo{}
	pgr := &mockProgressRepo{}
	lr := &mockLeadRepo{exists: true}
	svc := NewStepService(pr, sr, pgr, lr)

	preset, err := svc.CreatePreset(input.CreatePresetInput{
		Code: "standard", Name: "Standard Process", CreatedBy: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if preset.Code != "standard" {
		t.Errorf("Code = %v, want standard", preset.Code)
	}
}

func TestStepService_CreatePreset_Public(t *testing.T) {
	pr := &mockPresetRepo{}
	svc := NewStepService(pr, &mockStepRepo{}, &mockProgressRepo{}, &mockLeadRepo{})

	preset, err := svc.CreatePreset(input.CreatePresetInput{
		Code: "public", Name: "Public", CreatedBy: 1, IsPublic: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !preset.IsPublic {
		t.Error("should be public")
	}
}

func TestStepService_DeletePreset_NotFound(t *testing.T) {
	pr := &mockPresetRepo{exists: false}
	svc := NewStepService(pr, &mockStepRepo{}, &mockProgressRepo{}, &mockLeadRepo{})

	if err := svc.DeletePreset(999); err == nil {
		t.Error("expected error for non-existent preset")
	}
}

func TestStepService_MakePresetPublic(t *testing.T) {
	preset, _ := domain.NewLeadStepPreset("test", "Test", 1)
	pr := &mockPresetRepo{findByID: preset}
	svc := NewStepService(pr, &mockStepRepo{}, &mockProgressRepo{}, &mockLeadRepo{})

	if err := svc.MakePresetPublic(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepService_MakePresetShared(t *testing.T) {
	preset, _ := domain.NewLeadStepPreset("test", "Test", 1)
	pr := &mockPresetRepo{findByID: preset}
	svc := NewStepService(pr, &mockStepRepo{}, &mockProgressRepo{}, &mockLeadRepo{})

	if err := svc.MakePresetShared(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepService_MakePresetPrivate(t *testing.T) {
	preset, _ := domain.NewLeadStepPreset("test", "Test", 1)
	pr := &mockPresetRepo{findByID: preset}
	svc := NewStepService(pr, &mockStepRepo{}, &mockProgressRepo{}, &mockLeadRepo{})

	if err := svc.MakePresetPrivate(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepService_MakePreset_NotFound(t *testing.T) {
	pr := &mockPresetRepo{}
	svc := NewStepService(pr, &mockStepRepo{}, &mockProgressRepo{}, &mockLeadRepo{})

	if err := svc.MakePresetPublic(999); err == nil {
		t.Error("expected error")
	}
	if err := svc.MakePresetShared(999); err == nil {
		t.Error("expected error")
	}
	if err := svc.MakePresetPrivate(999); err == nil {
		t.Error("expected error")
	}
}

// --- Step Tests ---

func TestStepService_CreateStep(t *testing.T) {
	pr := &mockPresetRepo{exists: true}
	sr := &mockStepRepo{}
	svc := NewStepService(pr, sr, &mockProgressRepo{}, &mockLeadRepo{})

	step, err := svc.CreateStep(input.CreateStepInput{
		PresetID: 1, Code: "contact", Name: "First Contact", SortOrder: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if step.Code != "contact" {
		t.Errorf("Code = %v, want contact", step.Code)
	}
}

func TestStepService_CreateStep_Final(t *testing.T) {
	pr := &mockPresetRepo{exists: true}
	sr := &mockStepRepo{}
	svc := NewStepService(pr, sr, &mockProgressRepo{}, &mockLeadRepo{})

	step, err := svc.CreateStep(input.CreateStepInput{
		PresetID: 1, Code: "close", Name: "Close", SortOrder: 10, IsFinal: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !step.IsFinal {
		t.Error("step should be final")
	}
}

func TestStepService_CreateStep_PresetNotFound(t *testing.T) {
	pr := &mockPresetRepo{exists: false}
	svc := NewStepService(pr, &mockStepRepo{}, &mockProgressRepo{}, &mockLeadRepo{})

	_, err := svc.CreateStep(input.CreateStepInput{
		PresetID: 999, Code: "test", Name: "Test", SortOrder: 1,
	})
	if err == nil {
		t.Error("expected error for non-existent preset")
	}
}

func TestStepService_DeactivateStep(t *testing.T) {
	step, _ := domain.NewLeadStep(1, "contact", "Contact", 1)
	sr := &mockStepRepo{findByID: step}
	svc := NewStepService(&mockPresetRepo{}, sr, &mockProgressRepo{}, &mockLeadRepo{})

	if err := svc.DeactivateStep(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepService_ActivateStep(t *testing.T) {
	step, _ := domain.NewLeadStep(1, "contact", "Contact", 1)
	step.Deactivate()
	sr := &mockStepRepo{findByID: step}
	svc := NewStepService(&mockPresetRepo{}, sr, &mockProgressRepo{}, &mockLeadRepo{})

	if err := svc.ActivateStep(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepService_DeleteStep_NotFound(t *testing.T) {
	sr := &mockStepRepo{exists: false}
	svc := NewStepService(&mockPresetRepo{}, sr, &mockProgressRepo{}, &mockLeadRepo{})

	if err := svc.DeleteStep(999); err == nil {
		t.Error("expected error for non-existent step")
	}
}

// --- Progress Tests ---

func TestStepService_InitializeProgress(t *testing.T) {
	step1, _ := domain.NewLeadStep(1, "s1", "Step 1", 1)
	step1.ID = 10
	step2, _ := domain.NewLeadStep(1, "s2", "Step 2", 2)
	step2.ID = 20

	lr := &mockLeadRepo{exists: true}
	sr := &mockStepRepo{steps: []*domain.LeadStep{step1, step2}}
	pgr := &mockProgressRepo{}
	svc := NewStepService(&mockPresetRepo{}, sr, pgr, lr)

	if err := svc.InitializeProgress(1, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepService_InitializeProgress_LeadNotFound(t *testing.T) {
	lr := &mockLeadRepo{exists: false}
	svc := NewStepService(&mockPresetRepo{}, &mockStepRepo{}, &mockProgressRepo{}, lr)

	if err := svc.InitializeProgress(999, 1); err == nil {
		t.Error("expected error for non-existent lead")
	}
}

func TestStepService_CompleteStep(t *testing.T) {
	progress, _ := domain.NewLeadStepProgress(1, 1)
	pgr := &mockProgressRepo{findByLeadAndStep: progress}
	svc := NewStepService(&mockPresetRepo{}, &mockStepRepo{}, pgr, &mockLeadRepo{})

	if err := svc.CompleteStep(1, 1, 5, "done"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepService_SkipStep(t *testing.T) {
	progress, _ := domain.NewLeadStepProgress(1, 1)
	pgr := &mockProgressRepo{findByLeadAndStep: progress}
	svc := NewStepService(&mockPresetRepo{}, &mockStepRepo{}, pgr, &mockLeadRepo{})

	if err := svc.SkipStep(1, 1, 5, "not needed"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepService_FailStep(t *testing.T) {
	progress, _ := domain.NewLeadStepProgress(1, 1)
	pgr := &mockProgressRepo{findByLeadAndStep: progress}
	svc := NewStepService(&mockPresetRepo{}, &mockStepRepo{}, pgr, &mockLeadRepo{})

	if err := svc.FailStep(1, 1, 5, "rejected"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepService_UpdateProgress_Completed(t *testing.T) {
	progress, _ := domain.NewLeadStepProgress(1, 1)
	pgr := &mockProgressRepo{findByLeadAndStep: progress}
	svc := NewStepService(&mockPresetRepo{}, &mockStepRepo{}, pgr, &mockLeadRepo{})

	result, err := svc.UpdateProgress(input.UpdateProgressInput{
		LeadID: 1, StepID: 1, Status: "completed", CompletedBy: 5, Notes: "ok",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsCompleted() {
		t.Error("should be completed")
	}
}

func TestStepService_UpdateProgress_InvalidStatus(t *testing.T) {
	progress, _ := domain.NewLeadStepProgress(1, 1)
	pgr := &mockProgressRepo{findByLeadAndStep: progress}
	svc := NewStepService(&mockPresetRepo{}, &mockStepRepo{}, pgr, &mockLeadRepo{})

	_, err := svc.UpdateProgress(input.UpdateProgressInput{
		LeadID: 1, StepID: 1, Status: "invalid", CompletedBy: 1,
	})
	if err == nil {
		t.Error("expected error for invalid status")
	}
}

func TestStepService_UpdateProgress_NotFound(t *testing.T) {
	pgr := &mockProgressRepo{}
	svc := NewStepService(&mockPresetRepo{}, &mockStepRepo{}, pgr, &mockLeadRepo{})

	_, err := svc.UpdateProgress(input.UpdateProgressInput{
		LeadID: 999, StepID: 999, Status: "completed", CompletedBy: 1,
	})
	if err == nil {
		t.Error("expected error for non-existent progress")
	}
}