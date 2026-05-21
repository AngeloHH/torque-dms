package domain

import (
	"testing"
	"time"
)

func TestNewLeadActivity_Valid(t *testing.T) {
	a, err := NewLeadActivity(1, ActivityTypeCallOutbound, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.LeadID != 1 || a.Type != ActivityTypeCallOutbound || a.PerformedBy != 5 {
		t.Error("fields not set correctly")
	}
}

func TestNewLeadActivity_Invalid(t *testing.T) {
	if _, err := NewLeadActivity(0, ActivityTypeCallOutbound, 1); err == nil {
		t.Error("expected error for zero lead")
	}
	if _, err := NewLeadActivity(1, ActivityTypeCallOutbound, 0); err == nil {
		t.Error("expected error for zero performed_by")
	}
}

func TestLeadActivity_SetFields(t *testing.T) {
	a, _ := NewLeadActivity(1, ActivityTypeCallOutbound, 1)

	a.SetDescription("Follow up call")
	if a.Description != "Follow up call" {
		t.Errorf("Description = %v, want Follow up call", a.Description)
	}

	a.SetOutcome("Left voicemail")
	if a.Outcome != "Left voicemail" {
		t.Errorf("Outcome = %v, want Left voicemail", a.Outcome)
	}

	a.SetPhone(10)
	if a.PhoneID == nil || *a.PhoneID != 10 {
		t.Error("PhoneID not set correctly")
	}

	a.SetEmail("client@test.com")
	if a.Email != "client@test.com" {
		t.Errorf("Email = %v, want client@test.com", a.Email)
	}
}

func TestLeadActivity_ScheduleAndComplete(t *testing.T) {
	a, _ := NewLeadActivity(1, ActivityTypeAppointmentScheduled, 1)

	if a.IsScheduled() {
		t.Error("should not be scheduled initially")
	}
	if a.IsCompleted() {
		t.Error("should not be completed initially")
	}

	future := time.Now().Add(1 * time.Hour)
	a.Schedule(future)
	if !a.IsScheduled() {
		t.Error("should be scheduled after Schedule()")
	}

	a.Complete()
	if !a.IsCompleted() {
		t.Error("should be completed after Complete()")
	}
	if a.IsScheduled() {
		t.Error("completed activity should not be considered scheduled")
	}
}

func TestLeadActivity_IsOverdue(t *testing.T) {
	a, _ := NewLeadActivity(1, ActivityTypeCallOutbound, 1)

	if a.IsOverdue() {
		t.Error("unscheduled activity should not be overdue")
	}

	past := time.Now().Add(-1 * time.Hour)
	a.Schedule(past)
	if !a.IsOverdue() {
		t.Error("past scheduled activity should be overdue")
	}

	a.Complete()
	if a.IsOverdue() {
		t.Error("completed activity should not be overdue")
	}
}