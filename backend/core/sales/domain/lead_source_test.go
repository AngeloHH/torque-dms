package domain

import "testing"

func TestNewLeadSource_Valid(t *testing.T) {
	s, err := NewLeadSource("web", "Website", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Code != "web" || s.Name != "Website" || !s.IsExternal {
		t.Error("fields not set correctly")
	}
	if !s.Active {
		t.Error("new source should be active")
	}
}

func TestNewLeadSource_Invalid(t *testing.T) {
	if _, err := NewLeadSource("", "Name", false); err == nil {
		t.Error("expected error for empty code")
	}
	if _, err := NewLeadSource("code", "", false); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestLeadSource_ActiveToggle(t *testing.T) {
	s, _ := NewLeadSource("web", "Website", false)
	s.Deactivate()
	if s.Active {
		t.Error("should be inactive")
	}
	s.Activate()
	if !s.Active {
		t.Error("should be active")
	}
}