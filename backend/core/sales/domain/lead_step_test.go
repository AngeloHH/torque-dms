package domain

import "testing"

// --- LeadStepPreset ---

func TestNewLeadStepPreset_Valid(t *testing.T) {
	p, err := NewLeadStepPreset("standard", "Standard Process", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Code != "standard" || p.Name != "Standard Process" || p.CreatedBy != 1 {
		t.Error("fields not set correctly")
	}
	if p.IsPublic || p.IsShared {
		t.Error("new preset should be private")
	}
}

func TestNewLeadStepPreset_Invalid(t *testing.T) {
	if _, err := NewLeadStepPreset("", "Name", 1); err == nil {
		t.Error("expected error for empty code")
	}
	if _, err := NewLeadStepPreset("code", "", 1); err == nil {
		t.Error("expected error for empty name")
	}
	if _, err := NewLeadStepPreset("code", "Name", 0); err == nil {
		t.Error("expected error for zero created_by")
	}
}

func TestLeadStepPreset_Visibility(t *testing.T) {
	p, _ := NewLeadStepPreset("test", "Test", 1)

	p.MakePublic()
	if !p.IsPublic || p.IsShared {
		t.Error("MakePublic: should be public, not shared")
	}

	p.MakeShared()
	if p.IsPublic || !p.IsShared {
		t.Error("MakeShared: should be shared, not public")
	}

	p.MakePrivate()
	if p.IsPublic || p.IsShared {
		t.Error("MakePrivate: should be neither public nor shared")
	}
}

// --- LeadStep ---

func TestNewLeadStep_Valid(t *testing.T) {
	s, err := NewLeadStep(1, "contact", "First Contact", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.PresetID != 1 || s.Code != "contact" || s.Name != "First Contact" {
		t.Error("fields not set correctly")
	}
	if !s.Active {
		t.Error("new step should be active")
	}
	if s.IsFinal {
		t.Error("new step should not be final")
	}
}

func TestNewLeadStep_Invalid(t *testing.T) {
	if _, err := NewLeadStep(0, "code", "Name", 1); err == nil {
		t.Error("expected error for zero preset")
	}
	if _, err := NewLeadStep(1, "", "Name", 1); err == nil {
		t.Error("expected error for empty code")
	}
	if _, err := NewLeadStep(1, "code", "", 1); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestLeadStep_MarkAsFinal(t *testing.T) {
	s, _ := NewLeadStep(1, "close", "Close Deal", 10)
	s.MarkAsFinal()
	if !s.IsFinal {
		t.Error("should be final")
	}
}

func TestLeadStep_ActiveToggle(t *testing.T) {
	s, _ := NewLeadStep(1, "code", "Name", 1)
	s.Deactivate()
	if s.Active {
		t.Error("should be inactive")
	}
	s.Activate()
	if !s.Active {
		t.Error("should be active")
	}
}

// --- LeadStepProgress ---

func TestNewLeadStepProgress_Valid(t *testing.T) {
	p, err := NewLeadStepProgress(1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.LeadID != 1 || p.StepID != 2 {
		t.Error("fields not set correctly")
	}
	if p.Status != StepStatusPending {
		t.Errorf("Status = %v, want pending", p.Status)
	}
	if !p.IsPending() {
		t.Error("new progress should be pending")
	}
	if p.StartedAt == nil {
		t.Error("StartedAt should be set")
	}
}

func TestNewLeadStepProgress_Invalid(t *testing.T) {
	if _, err := NewLeadStepProgress(0, 1); err == nil {
		t.Error("expected error for zero lead")
	}
	if _, err := NewLeadStepProgress(1, 0); err == nil {
		t.Error("expected error for zero step")
	}
}

func TestLeadStepProgress_Complete(t *testing.T) {
	p, _ := NewLeadStepProgress(1, 1)
	p.Complete(5, "done")

	if !p.IsCompleted() {
		t.Error("should be completed")
	}
	if p.IsPending() {
		t.Error("should not be pending")
	}
	if p.CompletedAt == nil {
		t.Error("CompletedAt should be set")
	}
	if p.CompletedBy == nil || *p.CompletedBy != 5 {
		t.Error("CompletedBy not set correctly")
	}
	if p.Notes != "done" {
		t.Errorf("Notes = %v, want done", p.Notes)
	}
}

func TestLeadStepProgress_Skip(t *testing.T) {
	p, _ := NewLeadStepProgress(1, 1)
	p.Skip(5, "not applicable")

	if p.Status != StepStatusSkipped {
		t.Errorf("Status = %v, want skipped", p.Status)
	}
	if p.CompletedAt == nil {
		t.Error("CompletedAt should be set")
	}
}

func TestLeadStepProgress_Fail(t *testing.T) {
	p, _ := NewLeadStepProgress(1, 1)
	p.Fail(5, "client rejected")

	if p.Status != StepStatusFailed {
		t.Errorf("Status = %v, want failed", p.Status)
	}
}