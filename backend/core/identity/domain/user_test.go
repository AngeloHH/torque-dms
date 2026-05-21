package domain

import (
	"testing"

	sharedDomain "torque-dms/core/shared/domain"
)

func init() {
	sharedDomain.LoadValidationRules("../../../settings/validation_rules.yml")
}

func TestNewUserAccount_Valid(t *testing.T) {
	user, err := NewUserAccount(1, "juanperez", "Strong@123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.EntityID != 1 {
		t.Errorf("EntityID = %v, want 1", user.EntityID)
	}
	if user.Username != "juanperez" {
		t.Errorf("Username = %v, want juanperez", user.Username)
	}
	if user.Status != EntityStatusActive {
		t.Errorf("Status = %v, want active", user.Status)
	}
	if user.PasswordHash == "Strong@123" {
		t.Error("password should be hashed, not stored in plain text")
	}
}

func TestNewUserAccount_InvalidUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
	}{
		{"blacklisted admin", "admin"},
		{"blacklisted root", "root"},
		{"too short", "ab"},
		{"special chars", "user@name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUserAccount(1, tt.username, "Strong@123")
			if err == nil {
				t.Error("expected error for invalid username")
			}
		})
	}
}

func TestNewUserAccount_InvalidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"too short", "Ab@1"},
		{"no uppercase", "abcdefg@1"},
		{"no special", "Abcdefg12"},
		{"blacklisted", "password"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUserAccount(1, "validuser", tt.password)
			if err == nil {
				t.Error("expected error for invalid password")
			}
		})
	}
}

func TestNewUserAccount_ZeroEntityID(t *testing.T) {
	_, err := NewUserAccount(0, "validuser", "Strong@123")
	if err == nil {
		t.Error("expected error for zero entity ID")
	}
}

func TestUserAccount_CheckPassword(t *testing.T) {
	user, _ := NewUserAccount(1, "juanperez", "Strong@123")

	if !user.CheckPassword("Strong@123") {
		t.Error("correct password should return true")
	}
	if user.CheckPassword("WrongPass@1") {
		t.Error("wrong password should return false")
	}
}

func TestUserAccount_ChangePassword(t *testing.T) {
	user, _ := NewUserAccount(1, "juanperez", "Strong@123")

	// Wrong old password
	if err := user.ChangePassword("wrong", "NewStrong@1"); err == nil {
		t.Error("expected error for wrong old password")
	}

	// New password too short
	if err := user.ChangePassword("Strong@123", "short"); err == nil {
		t.Error("expected error for short new password")
	}

	// Valid change
	if err := user.ChangePassword("Strong@123", "NewStrong@1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !user.CheckPassword("NewStrong@1") {
		t.Error("new password should work after change")
	}
	if user.CheckPassword("Strong@123") {
		t.Error("old password should not work after change")
	}
}

func TestUserAccount_StatusMethods(t *testing.T) {
	user, _ := NewUserAccount(1, "juanperez", "Strong@123")

	if !user.IsActive() {
		t.Error("new user should be active")
	}

	user.Suspend()
	if user.IsActive() {
		t.Error("suspended user should not be active")
	}
	if user.Status != EntityStatusSuspended {
		t.Errorf("status = %v, want suspended", user.Status)
	}

	user.Activate()
	if !user.IsActive() {
		t.Error("activated user should be active")
	}
}

func TestUserAccount_RecordLogin(t *testing.T) {
	user, _ := NewUserAccount(1, "juanperez", "Strong@123")
	if !user.LastLogin.IsZero() {
		t.Error("LastLogin should be zero initially")
	}

	user.RecordLogin()
	if user.LastLogin.IsZero() {
		t.Error("LastLogin should be set after RecordLogin")
	}
}