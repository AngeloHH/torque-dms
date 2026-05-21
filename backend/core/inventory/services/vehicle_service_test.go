package services

import (
	"errors"
	"testing"

	"torque-dms/core/inventory/domain"
	"torque-dms/core/inventory/ports/input"
)

// --- Mocks ---

type mockVehicleRepo struct {
	saved       *domain.Vehicle
	saveErr     error
	findByID    *domain.Vehicle
	findErr     error
	existsVIN   bool
	existsID    bool
	vehicles    []*domain.Vehicle
	updated     *domain.Vehicle
}

func (m *mockVehicleRepo) Save(v *domain.Vehicle) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	v.ID = 1
	m.saved = v
	return nil
}
func (m *mockVehicleRepo) Update(v *domain.Vehicle) error { m.updated = v; return nil }
func (m *mockVehicleRepo) FindByID(id uint) (*domain.Vehicle, error) {
	if m.findByID != nil {
		return m.findByID, nil
	}
	return nil, errors.New("not found")
}
func (m *mockVehicleRepo) FindByVIN(vin string) (*domain.Vehicle, error)            { return nil, nil }
func (m *mockVehicleRepo) FindByStockNumber(sn string) (*domain.Vehicle, error)     { return nil, nil }
func (m *mockVehicleRepo) FindAll(l int, o int) ([]*domain.Vehicle, error)          { return m.vehicles, nil }
func (m *mockVehicleRepo) FindByStatus(s domain.VehicleStatus, l int, o int) ([]*domain.Vehicle, error) { return nil, nil }
func (m *mockVehicleRepo) FindByLocationID(id uint) ([]*domain.Vehicle, error)      { return m.vehicles, nil }
func (m *mockVehicleRepo) FindAvailable(l int, o int) ([]*domain.Vehicle, error)    { return nil, nil }
func (m *mockVehicleRepo) Delete(id uint) error                                     { return nil }
func (m *mockVehicleRepo) Exists(id uint) (bool, error)                             { return m.existsID, nil }
func (m *mockVehicleRepo) ExistsByVIN(vin string) (bool, error)                     { return m.existsVIN, nil }

type mockPhotoRepo struct {
	saved   *domain.VehiclePhoto
	photos  []*domain.VehiclePhoto
	findByID *domain.VehiclePhoto
}

func (m *mockPhotoRepo) Save(p *domain.VehiclePhoto) error    { p.ID = 1; m.saved = p; return nil }
func (m *mockPhotoRepo) Update(p *domain.VehiclePhoto) error  { return nil }
func (m *mockPhotoRepo) FindByID(id uint) (*domain.VehiclePhoto, error) {
	if m.findByID != nil {
		return m.findByID, nil
	}
	return nil, errors.New("not found")
}
func (m *mockPhotoRepo) FindByVehicleID(id uint) ([]*domain.VehiclePhoto, error) { return m.photos, nil }
func (m *mockPhotoRepo) FindPrimaryByVehicleID(id uint) (*domain.VehiclePhoto, error) { return nil, nil }
func (m *mockPhotoRepo) Delete(id uint) error                                    { return nil }
func (m *mockPhotoRepo) DeleteByVehicleID(id uint) error                         { return nil }

type mockLocationRepo struct {
	saved    *domain.Location
	findByID *domain.Location
	exists   bool
	locations []*domain.Location
}

func (m *mockLocationRepo) Save(l *domain.Location) error    { l.ID = 1; m.saved = l; return nil }
func (m *mockLocationRepo) Update(l *domain.Location) error  { return nil }
func (m *mockLocationRepo) FindByID(id uint) (*domain.Location, error) {
	if m.findByID != nil {
		return m.findByID, nil
	}
	return nil, errors.New("not found")
}
func (m *mockLocationRepo) FindByName(name string) (*domain.Location, error)                     { return nil, nil }
func (m *mockLocationRepo) FindAll() ([]*domain.Location, error)                                 { return m.locations, nil }
func (m *mockLocationRepo) FindByType(t domain.LocationType) ([]*domain.Location, error)         { return nil, nil }
func (m *mockLocationRepo) FindActive() ([]*domain.Location, error)                              { return nil, nil }
func (m *mockLocationRepo) Delete(id uint) error                                                 { return nil }
func (m *mockLocationRepo) Exists(id uint) (bool, error)                                         { return m.exists, nil }

// --- Vehicle Service Tests ---

func validInput() input.CreateVehicleInput {
	return input.CreateVehicleInput{
		StockNumber: "STK001",
		VIN:         "12345678901234567",
		Make:        "Toyota",
		Model:       "Corolla",
		Year:        2024,
	}
}

func TestVehicleService_Create_Success(t *testing.T) {
	vRepo := &mockVehicleRepo{existsVIN: false}
	pRepo := &mockPhotoRepo{}
	lRepo := &mockLocationRepo{}
	svc := NewVehicleService(vRepo, pRepo, lRepo)

	v, err := svc.Create(validInput())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.StockNumber != "STK001" {
		t.Errorf("StockNumber = %v, want STK001", v.StockNumber)
	}
}

func TestVehicleService_Create_DuplicateVIN(t *testing.T) {
	vRepo := &mockVehicleRepo{existsVIN: true}
	pRepo := &mockPhotoRepo{}
	lRepo := &mockLocationRepo{}
	svc := NewVehicleService(vRepo, pRepo, lRepo)

	_, err := svc.Create(validInput())
	if err == nil {
		t.Error("expected error for duplicate VIN")
	}
}

func TestVehicleService_Create_WithLocation(t *testing.T) {
	vRepo := &mockVehicleRepo{existsVIN: false}
	pRepo := &mockPhotoRepo{}
	lRepo := &mockLocationRepo{exists: true}
	svc := NewVehicleService(vRepo, pRepo, lRepo)

	inp := validInput()
	inp.LocationID = 1

	v, err := svc.Create(inp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.LocationID != 1 {
		t.Errorf("LocationID = %v, want 1", v.LocationID)
	}
}

func TestVehicleService_Create_InvalidLocation(t *testing.T) {
	vRepo := &mockVehicleRepo{existsVIN: false}
	pRepo := &mockPhotoRepo{}
	lRepo := &mockLocationRepo{exists: false}
	svc := NewVehicleService(vRepo, pRepo, lRepo)

	inp := validInput()
	inp.LocationID = 999

	_, err := svc.Create(inp)
	if err == nil {
		t.Error("expected error for non-existent location")
	}
}

func TestVehicleService_Create_WithPricing(t *testing.T) {
	vRepo := &mockVehicleRepo{existsVIN: false}
	pRepo := &mockPhotoRepo{}
	lRepo := &mockLocationRepo{}
	svc := NewVehicleService(vRepo, pRepo, lRepo)

	inp := validInput()
	inp.MSRP = 35000
	inp.InvoicePrice = 30000
	inp.AskingPrice = 33000

	v, err := svc.Create(inp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.MSRP != 35000 {
		t.Errorf("MSRP = %v, want 35000", v.MSRP)
	}
}

func TestVehicleService_Update_Success(t *testing.T) {
	existing, _ := domain.NewVehicle("STK001", "12345678901234567", "Toyota", "Corolla", 2024)
	existing.ID = 1

	vRepo := &mockVehicleRepo{findByID: existing}
	pRepo := &mockPhotoRepo{}
	lRepo := &mockLocationRepo{}
	svc := NewVehicleService(vRepo, pRepo, lRepo)

	plate := "ABC-1234"
	v, err := svc.Update(1, input.UpdateVehicleInput{Plate: &plate})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Plate != "ABC-1234" {
		t.Errorf("Plate = %v, want ABC-1234", v.Plate)
	}
}

func TestVehicleService_Update_SoldVehicle(t *testing.T) {
	existing, _ := domain.NewVehicle("STK001", "12345678901234567", "Toyota", "Corolla", 2024)
	existing.MarkAsSold()

	vRepo := &mockVehicleRepo{findByID: existing}
	pRepo := &mockPhotoRepo{}
	lRepo := &mockLocationRepo{}
	svc := NewVehicleService(vRepo, pRepo, lRepo)

	plate := "ABC-1234"
	_, err := svc.Update(1, input.UpdateVehicleInput{Plate: &plate})
	if err == nil {
		t.Error("expected error updating sold vehicle")
	}
}

func TestVehicleService_Delete_SoldVehicle(t *testing.T) {
	existing, _ := domain.NewVehicle("STK001", "12345678901234567", "Toyota", "Corolla", 2024)
	existing.MarkAsSold()

	vRepo := &mockVehicleRepo{findByID: existing}
	pRepo := &mockPhotoRepo{}
	lRepo := &mockLocationRepo{}
	svc := NewVehicleService(vRepo, pRepo, lRepo)

	err := svc.Delete(1)
	if err == nil {
		t.Error("expected error deleting sold vehicle")
	}
}

func TestVehicleService_MarkAsSold(t *testing.T) {
	existing, _ := domain.NewVehicle("STK001", "12345678901234567", "Toyota", "Corolla", 2024)
	existing.ID = 1

	vRepo := &mockVehicleRepo{findByID: existing}
	pRepo := &mockPhotoRepo{}
	lRepo := &mockLocationRepo{}
	svc := NewVehicleService(vRepo, pRepo, lRepo)

	if err := svc.MarkAsSold(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !vRepo.updated.IsSold() {
		t.Error("vehicle should be sold after MarkAsSold")
	}
}

func TestVehicleService_ChangeLocation(t *testing.T) {
	existing, _ := domain.NewVehicle("STK001", "12345678901234567", "Toyota", "Corolla", 2024)
	existing.ID = 1

	vRepo := &mockVehicleRepo{findByID: existing}
	pRepo := &mockPhotoRepo{}
	lRepo := &mockLocationRepo{exists: true}
	svc := NewVehicleService(vRepo, pRepo, lRepo)

	if err := svc.ChangeLocation(1, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vRepo.updated.LocationID != 5 {
		t.Errorf("LocationID = %v, want 5", vRepo.updated.LocationID)
	}
}

func TestVehicleService_ChangeLocation_NotFound(t *testing.T) {
	existing, _ := domain.NewVehicle("STK001", "12345678901234567", "Toyota", "Corolla", 2024)

	vRepo := &mockVehicleRepo{findByID: existing}
	pRepo := &mockPhotoRepo{}
	lRepo := &mockLocationRepo{exists: false}
	svc := NewVehicleService(vRepo, pRepo, lRepo)

	err := svc.ChangeLocation(1, 999)
	if err == nil {
		t.Error("expected error for non-existent location")
	}
}

func TestVehicleService_AddPhoto(t *testing.T) {
	vRepo := &mockVehicleRepo{existsID: true}
	pRepo := &mockPhotoRepo{}
	lRepo := &mockLocationRepo{}
	svc := NewVehicleService(vRepo, pRepo, lRepo)

	photo, err := svc.AddPhoto(input.AddPhotoInput{
		VehicleID:   1,
		URL:         "https://img.test/photo.jpg",
		Perspective: "front",
		Purpose:     "listing",
		UploadedBy:  5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if photo.VehicleID != 1 {
		t.Errorf("VehicleID = %v, want 1", photo.VehicleID)
	}
}

func TestVehicleService_AddPhoto_VehicleNotFound(t *testing.T) {
	vRepo := &mockVehicleRepo{existsID: false}
	pRepo := &mockPhotoRepo{}
	lRepo := &mockLocationRepo{}
	svc := NewVehicleService(vRepo, pRepo, lRepo)

	_, err := svc.AddPhoto(input.AddPhotoInput{
		VehicleID:   999,
		URL:         "https://img.test/photo.jpg",
		Perspective: "front",
		Purpose:     "listing",
		UploadedBy:  5,
	})
	if err == nil {
		t.Error("expected error for non-existent vehicle")
	}
}

func TestVehicleService_SetPrimaryPhoto(t *testing.T) {
	photo := &domain.VehiclePhoto{ID: 10, VehicleID: 1}

	vRepo := &mockVehicleRepo{}
	pRepo := &mockPhotoRepo{
		findByID: photo,
		photos:   []*domain.VehiclePhoto{photo},
	}
	lRepo := &mockLocationRepo{}
	svc := NewVehicleService(vRepo, pRepo, lRepo)

	err := svc.SetPrimaryPhoto(1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVehicleService_SetPrimaryPhoto_WrongVehicle(t *testing.T) {
	photo := &domain.VehiclePhoto{ID: 10, VehicleID: 2}

	vRepo := &mockVehicleRepo{}
	pRepo := &mockPhotoRepo{findByID: photo}
	lRepo := &mockLocationRepo{}
	svc := NewVehicleService(vRepo, pRepo, lRepo)

	err := svc.SetPrimaryPhoto(1, 10)
	if err == nil {
		t.Error("expected error when photo belongs to different vehicle")
	}
}