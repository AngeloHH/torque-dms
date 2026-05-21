package domain

import "testing"

func TestNewVehiclePhoto_Valid(t *testing.T) {
	p, err := NewVehiclePhoto(1, "https://img.test/photo.jpg", PhotoPerspectiveFront, PhotoPurposeListing, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.VehicleID != 1 {
		t.Errorf("VehicleID = %v, want 1", p.VehicleID)
	}
	if p.IsPrimary {
		t.Error("new photo should not be primary")
	}
}

func TestNewVehiclePhoto_Invalid(t *testing.T) {
	tests := []struct {
		name       string
		vehicleID  uint
		url        string
		uploadedBy uint
	}{
		{"zero vehicle", 0, "https://img.test/p.jpg", 1},
		{"empty url", 1, "", 1},
		{"zero uploader", 1, "https://img.test/p.jpg", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewVehiclePhoto(tt.vehicleID, tt.url, PhotoPerspectiveFront, PhotoPurposeListing, tt.uploadedBy)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestVehiclePhoto_Primary(t *testing.T) {
	p, _ := NewVehiclePhoto(1, "https://img.test/p.jpg", PhotoPerspectiveFront, PhotoPurposeListing, 1)

	p.SetAsPrimary()
	if !p.IsPrimary {
		t.Error("should be primary after SetAsPrimary")
	}

	p.RemovePrimary()
	if p.IsPrimary {
		t.Error("should not be primary after RemovePrimary")
	}
}

func TestVehiclePhoto_SortOrder(t *testing.T) {
	p, _ := NewVehiclePhoto(1, "https://img.test/p.jpg", PhotoPerspectiveFront, PhotoPurposeListing, 1)
	p.SetSortOrder(5)
	if p.SortOrder != 5 {
		t.Errorf("SortOrder = %v, want 5", p.SortOrder)
	}
}