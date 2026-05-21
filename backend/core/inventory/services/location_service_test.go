package services

import (
	"errors"
	"testing"

	"torque-dms/core/inventory/domain"
	"torque-dms/core/inventory/ports/input"
)

// Re-use mockLocationRepo and mockVehicleRepo from vehicle_service_test.go

func TestLocationService_Create_Success(t *testing.T) {
	lRepo := &mockLocationRepo{}
	vRepo := &mockVehicleRepo{}
	svc := NewLocationService(lRepo, vRepo)

	loc, err := svc.Create(input.CreateLocationInput{
		Name:    "Main Lot",
		Type:    "sales_lot",
		Address: "123 Main St",
		City:    "Miami",
		State:   "FL",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.Name != "Main Lot" {
		t.Errorf("Name = %v, want Main Lot", loc.Name)
	}
}

func TestLocationService_Create_EmptyName(t *testing.T) {
	lRepo := &mockLocationRepo{}
	vRepo := &mockVehicleRepo{}
	svc := NewLocationService(lRepo, vRepo)

	_, err := svc.Create(input.CreateLocationInput{
		Name: "",
		Type: "sales_lot",
	})
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestLocationService_Create_WithCoordinates(t *testing.T) {
	lRepo := &mockLocationRepo{}
	vRepo := &mockVehicleRepo{}
	svc := NewLocationService(lRepo, vRepo)

	loc, err := svc.Create(input.CreateLocationInput{
		Name:      "Lot",
		Type:      "sales_lot",
		Latitude:  25.7617,
		Longitude: -80.1918,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.Latitude != 25.7617 {
		t.Errorf("Latitude = %v, want 25.7617", loc.Latitude)
	}
}

func TestLocationService_Create_InvalidCoordinates(t *testing.T) {
	lRepo := &mockLocationRepo{}
	vRepo := &mockVehicleRepo{}
	svc := NewLocationService(lRepo, vRepo)

	_, err := svc.Create(input.CreateLocationInput{
		Name:      "Lot",
		Type:      "sales_lot",
		Latitude:  999,
		Longitude: 0,
	})
	if err == nil {
		t.Error("expected error for invalid coordinates")
	}
}

func TestLocationService_Update_Success(t *testing.T) {
	existing, _ := domain.NewLocation("Old Name", domain.LocationTypeSalesLot)
	existing.ID = 1

	lRepo := &mockLocationRepo{findByID: existing}
	vRepo := &mockVehicleRepo{}
	svc := NewLocationService(lRepo, vRepo)

	newName := "New Name"
	loc, err := svc.Update(1, input.UpdateLocationInput{Name: &newName})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.Name != "New Name" {
		t.Errorf("Name = %v, want New Name", loc.Name)
	}
}

func TestLocationService_Update_NotFound(t *testing.T) {
	lRepo := &mockLocationRepo{}
	vRepo := &mockVehicleRepo{}
	svc := NewLocationService(lRepo, vRepo)

	newName := "Test"
	_, err := svc.Update(999, input.UpdateLocationInput{Name: &newName})
	if err == nil {
		t.Error("expected error for non-existent location")
	}
}

func TestLocationService_Delete_Success(t *testing.T) {
	lRepo := &mockLocationRepo{}
	vRepo := &mockVehicleRepo{vehicles: []*domain.Vehicle{}}
	svc := NewLocationService(lRepo, vRepo)

	err := svc.Delete(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLocationService_Delete_HasVehicles(t *testing.T) {
	v, _ := domain.NewVehicle("STK001", "12345678901234567", "Toyota", "Corolla", 2024)
	lRepo := &mockLocationRepo{}
	vRepo := &mockVehicleRepo{vehicles: []*domain.Vehicle{v}}
	svc := NewLocationService(lRepo, vRepo)

	err := svc.Delete(1)
	if err == nil {
		t.Error("expected error deleting location with vehicles")
	}
}

func TestLocationService_Deactivate(t *testing.T) {
	loc, _ := domain.NewLocation("Test", domain.LocationTypeSalesLot)
	loc.ID = 1

	lRepo := &mockLocationRepo{findByID: loc}
	vRepo := &mockVehicleRepo{}
	svc := NewLocationService(lRepo, vRepo)

	if err := svc.Deactivate(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLocationService_Deactivate_NotFound(t *testing.T) {
	lRepo := &mockLocationRepo{}
	vRepo := &mockVehicleRepo{}
	svc := NewLocationService(lRepo, vRepo)

	err := svc.Deactivate(999)
	if err == nil {
		t.Error("expected error for non-existent location")
	}
}

func TestLocationService_Activate(t *testing.T) {
	loc, _ := domain.NewLocation("Test", domain.LocationTypeSalesLot)
	loc.ID = 1
	loc.Deactivate()

	lRepo := &mockLocationRepo{findByID: loc}
	vRepo := &mockVehicleRepo{}
	svc := NewLocationService(lRepo, vRepo)

	if err := svc.Activate(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// mockVehicleRepo and mockLocationRepo need FindByLocationID to return error for unknown
type mockVehicleRepoWithErr struct {
	mockVehicleRepo
	findByLocErr error
}

func (m *mockVehicleRepoWithErr) FindByLocationID(id uint) ([]*domain.Vehicle, error) {
	if m.findByLocErr != nil {
		return nil, m.findByLocErr
	}
	return m.vehicles, nil
}

func TestLocationService_Delete_RepoError(t *testing.T) {
	lRepo := &mockLocationRepo{}
	vRepo := &mockVehicleRepoWithErr{
		mockVehicleRepo: mockVehicleRepo{},
		findByLocErr:    errors.New("db error"),
	}
	svc := NewLocationService(lRepo, vRepo)

	err := svc.Delete(1)
	if err == nil {
		t.Error("expected error when vehicle repo fails")
	}
}