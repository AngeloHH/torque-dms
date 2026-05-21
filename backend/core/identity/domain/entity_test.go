package domain

import (
	"testing"
)

func TestNewEntity_Valid(t *testing.T) {
	tests := []struct {
		name  string
		typ   EntityType
		phone string
		email string
	}{
		{"with email only", EntityTypePerson, "", "test@example.com"},
		{"with phone only", EntityTypePerson, "+1 7875551234", ""},
		{"with both", EntityTypeCompany, "+1 7875551234", "test@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := NewEntity(tt.typ, tt.phone, tt.email)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if entity.Type != tt.typ {
				t.Errorf("type = %v, want %v", entity.Type, tt.typ)
			}
			if entity.Status != EntityStatusActive {
				t.Errorf("status = %v, want active", entity.Status)
			}
		})
	}
}

func TestNewEntity_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		email string
	}{
		{"both empty", "", ""},
		{"invalid email", "", "notanemail"},
		{"invalid phone", "12345", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEntity(EntityTypePerson, tt.phone, tt.email)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestEntity_SetField(t *testing.T) {
	entity, _ := NewEntity(EntityTypePerson, "", "test@example.com")

	if err := entity.SetField("first_name", "Juan"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entity.FirstName != "Juan" {
		t.Errorf("FirstName = %v, want Juan", entity.FirstName)
	}

	if err := entity.SetField("last_name", "Perez"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entity.LastName != "Perez" {
		t.Errorf("LastName = %v, want Perez", entity.LastName)
	}

	if err := entity.SetField("invalid_field", "value"); err == nil {
		t.Error("expected error for invalid field")
	}
}

func TestEntity_SetField_InvalidEmail(t *testing.T) {
	entity, _ := NewEntity(EntityTypePerson, "", "test@example.com")
	err := entity.SetField("email", "notvalid")
	if err == nil {
		t.Error("expected error for invalid email")
	}
}

func TestEntity_StatusTransitions(t *testing.T) {
	entity, _ := NewEntity(EntityTypePerson, "", "test@example.com")

	// Active -> Suspend
	if err := entity.Suspend(); err != nil {
		t.Fatalf("Suspend() error: %v", err)
	}
	if entity.Status != EntityStatusSuspended {
		t.Errorf("status = %v, want suspended", entity.Status)
	}

	// Already suspended
	if err := entity.Suspend(); err == nil {
		t.Error("expected error suspending already suspended entity")
	}

	// Suspended -> Active
	if err := entity.Activate(); err != nil {
		t.Fatalf("Activate() error: %v", err)
	}

	// Already active
	if err := entity.Activate(); err == nil {
		t.Error("expected error activating already active entity")
	}

	// Active -> Inactive
	if err := entity.Deactivate(); err != nil {
		t.Fatalf("Deactivate() error: %v", err)
	}
	if entity.Status != EntityStatusInactive {
		t.Errorf("status = %v, want inactive", entity.Status)
	}

	// Already inactive
	if err := entity.Deactivate(); err == nil {
		t.Error("expected error deactivating already inactive entity")
	}
}

func TestEntity_CanLogin(t *testing.T) {
	entity, _ := NewEntity(EntityTypePerson, "", "test@example.com")

	if entity.CanLogin() {
		t.Error("non-system user should not be able to login")
	}

	entity.SetAsSystemUser()
	if !entity.CanLogin() {
		t.Error("active system user should be able to login")
	}

	entity.Suspend()
	if entity.CanLogin() {
		t.Error("suspended system user should not be able to login")
	}
}

func TestEntity_TypeChecks(t *testing.T) {
	person, _ := NewEntity(EntityTypePerson, "", "a@b.com")
	if !person.IsPerson() {
		t.Error("expected IsPerson() = true")
	}
	if person.IsCompany() {
		t.Error("expected IsCompany() = false")
	}

	company, _ := NewEntity(EntityTypeCompany, "", "a@b.com")
	if !company.IsCompany() {
		t.Error("expected IsCompany() = true")
	}
}

func TestEntity_EmailNormalized(t *testing.T) {
	entity, _ := NewEntity(EntityTypePerson, "", "Test@EXAMPLE.com")
	if entity.Email != "test@example.com" {
		t.Errorf("email = %v, want test@example.com", entity.Email)
	}
}