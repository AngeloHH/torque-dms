package domain

import "testing"

func TestNewLeadNote_Valid(t *testing.T) {
	n, err := NewLeadNote(1, "Client interested in financing", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.LeadID != 1 || n.Content != "Client interested in financing" || n.CreatedBy != 5 {
		t.Error("fields not set correctly")
	}
}

func TestNewLeadNote_Invalid(t *testing.T) {
	if _, err := NewLeadNote(0, "content", 1); err == nil {
		t.Error("expected error for zero lead")
	}
	if _, err := NewLeadNote(1, "", 1); err == nil {
		t.Error("expected error for empty content")
	}
	if _, err := NewLeadNote(1, "content", 0); err == nil {
		t.Error("expected error for zero created_by")
	}
}

func TestLeadNote_Update(t *testing.T) {
	n, _ := NewLeadNote(1, "original", 1)

	if err := n.Update("updated content"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.Content != "updated content" {
		t.Errorf("Content = %v, want updated content", n.Content)
	}
}

func TestLeadNote_Update_Empty(t *testing.T) {
	n, _ := NewLeadNote(1, "original", 1)
	if err := n.Update(""); err == nil {
		t.Error("expected error for empty content")
	}
}