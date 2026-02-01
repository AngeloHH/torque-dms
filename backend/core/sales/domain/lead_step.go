package domain

import (
	"errors"
	"time"
)

type StepStatus string

const (
	StepStatusPending   StepStatus = "pending"
	StepStatusCompleted StepStatus = "completed"
	StepStatusSkipped   StepStatus = "skipped"
	StepStatusFailed    StepStatus = "failed"
)

type LeadStepPreset struct {
	ID          uint
	Code        string
	Name        string
	Description string
	SortOrder   int
	IsPublic    bool
	IsShared    bool
	CreatedBy   uint
	CreatedAt   time.Time
}

func NewLeadStepPreset(code string, name string, createdBy uint) (*LeadStepPreset, error) {
	if code == "" {
		return nil, errors.New("code is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if createdBy == 0 {
		return nil, errors.New("created_by is required")
	}

	return &LeadStepPreset{
		Code:      code,
		Name:      name,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}, nil
}

func (p *LeadStepPreset) MakePublic() {
	p.IsPublic = true
	p.IsShared = false
}

func (p *LeadStepPreset) MakeShared() {
	p.IsShared = true
	p.IsPublic = false
}

func (p *LeadStepPreset) MakePrivate() {
	p.IsPublic = false
	p.IsShared = false
}

type LeadStep struct {
	ID        uint
	PresetID  uint
	Code      string
	Name      string
	SortOrder int
	IsFinal   bool
	Active    bool
	CreatedAt time.Time
}

func NewLeadStep(presetID uint, code string, name string, sortOrder int) (*LeadStep, error) {
	if presetID == 0 {
		return nil, errors.New("preset is required")
	}
	if code == "" {
		return nil, errors.New("code is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	return &LeadStep{
		PresetID:  presetID,
		Code:      code,
		Name:      name,
		SortOrder: sortOrder,
		Active:    true,
		CreatedAt: time.Now(),
	}, nil
}

func (s *LeadStep) MarkAsFinal() {
	s.IsFinal = true
}

func (s *LeadStep) Deactivate() {
	s.Active = false
}

func (s *LeadStep) Activate() {
	s.Active = true
}

type LeadStepProgress struct {
	ID          uint
	LeadID      uint
	StepID      uint
	Status      StepStatus
	StartedAt   *time.Time
	CompletedAt *time.Time
	CompletedBy *uint
	Notes       string
	CreatedAt   time.Time
}

func NewLeadStepProgress(leadID uint, stepID uint) (*LeadStepProgress, error) {
	if leadID == 0 {
		return nil, errors.New("lead is required")
	}
	if stepID == 0 {
		return nil, errors.New("step is required")
	}

	now := time.Now()
	return &LeadStepProgress{
		LeadID:    leadID,
		StepID:    stepID,
		Status:    StepStatusPending,
		StartedAt: &now,
		CreatedAt: now,
	}, nil
}

func (p *LeadStepProgress) Complete(completedBy uint, notes string) {
	now := time.Now()
	p.Status = StepStatusCompleted
	p.CompletedAt = &now
	p.CompletedBy = &completedBy
	p.Notes = notes
}

func (p *LeadStepProgress) Skip(completedBy uint, notes string) {
	now := time.Now()
	p.Status = StepStatusSkipped
	p.CompletedAt = &now
	p.CompletedBy = &completedBy
	p.Notes = notes
}

func (p *LeadStepProgress) Fail(completedBy uint, notes string) {
	now := time.Now()
	p.Status = StepStatusFailed
	p.CompletedAt = &now
	p.CompletedBy = &completedBy
	p.Notes = notes
}

func (p *LeadStepProgress) IsCompleted() bool {
	return p.Status == StepStatusCompleted
}

func (p *LeadStepProgress) IsPending() bool {
	return p.Status == StepStatusPending
}