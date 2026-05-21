package domain

import "testing"

func TestNewLead_Valid(t *testing.T) {
	lead, err := NewLead(1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lead.EntityID != 1 {
		t.Errorf("EntityID = %v, want 1", lead.EntityID)
	}
	if lead.SourceID != 2 {
		t.Errorf("SourceID = %v, want 2", lead.SourceID)
	}
}

func TestNewLead_Invalid(t *testing.T) {
	if _, err := NewLead(0, 1); err == nil {
		t.Error("expected error for zero entity")
	}
	if _, err := NewLead(1, 0); err == nil {
		t.Error("expected error for zero source")
	}
}

func TestLead_SetInterest(t *testing.T) {
	lead, _ := NewLead(1, 1)
	lead.SetInterest("purchase", "Toyota", "Corolla")

	if lead.InterestType != "purchase" {
		t.Errorf("InterestType = %v, want purchase", lead.InterestType)
	}
	if lead.InterestMake != "Toyota" {
		t.Errorf("InterestMake = %v, want Toyota", lead.InterestMake)
	}
}

func TestLead_SetBudget_Valid(t *testing.T) {
	lead, _ := NewLead(1, 1)
	if err := lead.SetBudget(20000, 40000); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lead.BudgetMin != 20000 || lead.BudgetMax != 40000 {
		t.Errorf("budget = %v-%v, want 20000-40000", lead.BudgetMin, lead.BudgetMax)
	}
}

func TestLead_SetBudget_Invalid(t *testing.T) {
	lead, _ := NewLead(1, 1)

	if err := lead.SetBudget(-1, 1000); err == nil {
		t.Error("expected error for negative min")
	}
	if err := lead.SetBudget(0, -1); err == nil {
		t.Error("expected error for negative max")
	}
	if err := lead.SetBudget(50000, 20000); err == nil {
		t.Error("expected error when min > max")
	}
}

func TestLead_SetBudget_ZeroMax(t *testing.T) {
	lead, _ := NewLead(1, 1)
	if err := lead.SetBudget(20000, 0); err != nil {
		t.Fatalf("should allow zero max (no upper limit): %v", err)
	}
}

func TestLead_Vehicle(t *testing.T) {
	lead, _ := NewLead(1, 1)

	if lead.HasVehicleInterest() {
		t.Error("should not have vehicle interest initially")
	}

	lead.SetVehicle(42)
	if !lead.HasVehicleInterest() {
		t.Error("should have vehicle interest after SetVehicle")
	}
	if *lead.VehicleID != 42 {
		t.Errorf("VehicleID = %v, want 42", *lead.VehicleID)
	}

	lead.RemoveVehicle()
	if lead.HasVehicleInterest() {
		t.Error("should not have vehicle interest after RemoveVehicle")
	}
}

func TestLead_IsHighValue(t *testing.T) {
	lead, _ := NewLead(1, 1)
	lead.SetBudget(0, 49999)
	if lead.IsHighValue() {
		t.Error("49999 should not be high value")
	}

	lead.SetBudget(0, 50000)
	if !lead.IsHighValue() {
		t.Error("50000 should be high value")
	}
}

func TestLead_SetPreset(t *testing.T) {
	lead, _ := NewLead(1, 1)
	lead.SetPreset(5)
	if lead.PresetID == nil || *lead.PresetID != 5 {
		t.Errorf("PresetID = %v, want 5", lead.PresetID)
	}
}