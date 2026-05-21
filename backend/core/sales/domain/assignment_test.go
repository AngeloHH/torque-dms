package domain

import "testing"

func TestNewLeadAssignment_Valid(t *testing.T) {
	a, err := NewLeadAssignment(1, 2, AssignmentRoleSalesperson, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.LeadID != 1 || a.EntityID != 2 || a.AssignedBy != 3 {
		t.Error("fields not set correctly")
	}
	if !a.Active {
		t.Error("new assignment should be active")
	}
	if a.IsPrimary {
		t.Error("new assignment should not be primary")
	}
}

func TestNewLeadAssignment_Invalid(t *testing.T) {
	if _, err := NewLeadAssignment(0, 1, AssignmentRoleSalesperson, 1); err == nil {
		t.Error("expected error for zero lead")
	}
	if _, err := NewLeadAssignment(1, 0, AssignmentRoleSalesperson, 1); err == nil {
		t.Error("expected error for zero entity")
	}
	if _, err := NewLeadAssignment(1, 1, AssignmentRoleSalesperson, 0); err == nil {
		t.Error("expected error for zero assigned_by")
	}
}

func TestLeadAssignment_Primary(t *testing.T) {
	a, _ := NewLeadAssignment(1, 1, AssignmentRoleSalesperson, 1)
	a.SetAsPrimary()
	if !a.IsPrimary {
		t.Error("should be primary")
	}
	a.RemovePrimary()
	if a.IsPrimary {
		t.Error("should not be primary")
	}
}

func TestLeadAssignment_ActiveToggle(t *testing.T) {
	a, _ := NewLeadAssignment(1, 1, AssignmentRoleSalesperson, 1)
	a.Deactivate()
	if a.Active {
		t.Error("should be inactive")
	}
	a.Activate()
	if !a.Active {
		t.Error("should be active")
	}
}