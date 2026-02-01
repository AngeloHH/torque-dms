package domain

import (
	"testing"
)

func TestMain(m *testing.M) {
	// Se ejecuta antes de todos los tests
	err := LoadValidationRules("../../../settings/validation_rules.yml")
	if err != nil {
		panic("Failed to load validation rules: " + err.Error())
	}
	m.Run()
}

func TestValidate_EntityID(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"zero is invalid", "0", true},
		{"positive is valid", "5", false},
		{"large number is valid", "12345", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate("entity_id", tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidate_Username(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"admin is blacklisted", "admin", true},
		{"root is blacklisted", "root", true},
		{"normal username is valid", "juan", false},
		{"too short", "ab", true},
		{"with underscore is valid", "juan_perez", false},
		{"with special chars is invalid", "juan@perez", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate("username", tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidate_Password(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"too short", "Abc@1", true},
		{"no uppercase", "abcdefg@1", true},
		{"no lowercase", "ABCDEFG@1", true},
		{"no number", "Abcdefg@a", true},
		{"no special char", "Abcdefg12", true},
		{"valid password", "Strong@123", false},
		{"blacklisted password", "password", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate("password", tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidate_Email(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid email", "juan@gmail.com", false},
		{"valid with subdomain", "juan@mail.empresa.com", false},
		{"missing @", "juangmail.com", true},
		{"missing domain", "juan@", true},
		{"missing tld", "juan@gmail", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate("email", tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidate_Phone(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid US", "+1 7875551234", false},
		{"valid MX", "+52 5512345678", false},
		{"missing plus", "1 7875551234", true},
		{"missing space", "+17875551234", true},
		{"too short number", "+1 12345", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate("phone", tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidate_UnknownField(t *testing.T) {
	err := Validate("unknown_field", "any_value")
	if err != nil {
		t.Errorf("Unknown field should return nil, got %v", err)
	}
}