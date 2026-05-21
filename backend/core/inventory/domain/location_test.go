package domain

import "testing"

func TestNewLocation_Valid(t *testing.T) {
	loc, err := NewLocation("Main Lot", LocationTypeSalesLot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.Name != "Main Lot" {
		t.Errorf("Name = %v, want Main Lot", loc.Name)
	}
	if !loc.Active {
		t.Error("new location should be active")
	}
	if !loc.IsSalesLot() {
		t.Error("should be sales lot")
	}
}

func TestNewLocation_EmptyName(t *testing.T) {
	_, err := NewLocation("", LocationTypeSalesLot)
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestLocation_SetAddress(t *testing.T) {
	loc, _ := NewLocation("Test Lot", LocationTypeSalesLot)
	loc.SetAddress("123 Main St", "Miami", "FL", "33101", 1)

	if loc.Address != "123 Main St" {
		t.Errorf("Address = %v, want 123 Main St", loc.Address)
	}
	if loc.City != "Miami" {
		t.Errorf("City = %v, want Miami", loc.City)
	}
}

func TestLocation_SetCoordinates_Valid(t *testing.T) {
	loc, _ := NewLocation("Test", LocationTypeSalesLot)
	if err := loc.SetCoordinates(25.7617, -80.1918); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.Latitude != 25.7617 {
		t.Errorf("Latitude = %v, want 25.7617", loc.Latitude)
	}
}

func TestLocation_SetCoordinates_Invalid(t *testing.T) {
	loc, _ := NewLocation("Test", LocationTypeSalesLot)

	if err := loc.SetCoordinates(91, 0); err == nil {
		t.Error("expected error for latitude > 90")
	}
	if err := loc.SetCoordinates(-91, 0); err == nil {
		t.Error("expected error for latitude < -90")
	}
	if err := loc.SetCoordinates(0, 181); err == nil {
		t.Error("expected error for longitude > 180")
	}
	if err := loc.SetCoordinates(0, -181); err == nil {
		t.Error("expected error for longitude < -180")
	}
}

func TestLocation_SetCapacity(t *testing.T) {
	loc, _ := NewLocation("Test", LocationTypeSalesLot)
	if err := loc.SetCapacity(100); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.Capacity != 100 {
		t.Errorf("Capacity = %v, want 100", loc.Capacity)
	}
}

func TestLocation_SetCapacity_Negative(t *testing.T) {
	loc, _ := NewLocation("Test", LocationTypeSalesLot)
	if err := loc.SetCapacity(-1); err == nil {
		t.Error("expected error for negative capacity")
	}
}

func TestLocation_ActivateDeactivate(t *testing.T) {
	loc, _ := NewLocation("Test", LocationTypeSalesLot)

	loc.Deactivate()
	if loc.Active {
		t.Error("should be inactive after Deactivate")
	}

	loc.Activate()
	if !loc.Active {
		t.Error("should be active after Activate")
	}
}

func TestLocation_IsSalesLot(t *testing.T) {
	sales, _ := NewLocation("Sales", LocationTypeSalesLot)
	if !sales.IsSalesLot() {
		t.Error("expected IsSalesLot() = true")
	}

	storage, _ := NewLocation("Storage", LocationTypeStorage)
	if storage.IsSalesLot() {
		t.Error("expected IsSalesLot() = false for storage")
	}
}