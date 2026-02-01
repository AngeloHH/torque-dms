package services

import (
	"errors"
	"time"

	"torque-dms/core/inventory/domain"
	"torque-dms/core/inventory/ports/input"
	"torque-dms/core/inventory/ports/output"
)

type vehicleService struct {
	vehicleRepo output.VehicleRepository
	photoRepo   output.VehiclePhotoRepository
	locationRepo output.LocationRepository
}

func NewVehicleService(
	vehicleRepo output.VehicleRepository,
	photoRepo output.VehiclePhotoRepository,
	locationRepo output.LocationRepository,
) input.VehicleService {
	return &vehicleService{
		vehicleRepo:  vehicleRepo,
		photoRepo:    photoRepo,
		locationRepo: locationRepo,
	}
}

func (s *vehicleService) Create(inp input.CreateVehicleInput) (*domain.Vehicle, error) {
	// Verificar que VIN no exista
	exists, err := s.vehicleRepo.ExistsByVIN(inp.VIN)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("vehicle with this VIN already exists")
	}

	// Verificar que location exista
	if inp.LocationID != 0 {
		locationExists, err := s.locationRepo.Exists(inp.LocationID)
		if err != nil {
			return nil, err
		}
		if !locationExists {
			return nil, errors.New("location not found")
		}
	}

	// Crear vehicle
	vehicle, err := domain.NewVehicle(inp.StockNumber, inp.VIN, inp.Make, inp.Model, inp.Year)
	if err != nil {
		return nil, err
	}

	// Setear campos opcionales
	vehicle.Plate = inp.Plate
	vehicle.Trim = inp.Trim
	vehicle.Mileage = inp.Mileage
	vehicle.ExteriorColor = inp.ExteriorColor
	vehicle.InteriorColor = inp.InteriorColor
	vehicle.LocationID = inp.LocationID

	// Setear precios
	if inp.MSRP > 0 || inp.InvoicePrice > 0 || inp.AskingPrice > 0 {
		if err := vehicle.SetPricing(inp.MSRP, inp.InvoicePrice, inp.AskingPrice); err != nil {
			return nil, err
		}
	}

	// Setear condición
	if inp.Condition != "" {
		vehicle.SetCondition(domain.VehicleCondition(inp.Condition))
	}

	// Setear adquisición
	if inp.AcquisitionSource != "" {
		if err := vehicle.SetAcquisition(
			domain.AcquisitionSource(inp.AcquisitionSource),
			inp.AcquisitionCost,
			time.Now(),
		); err != nil {
			return nil, err
		}
	}

	// Guardar
	if err := s.vehicleRepo.Save(vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}

func (s *vehicleService) GetByID(id uint) (*domain.Vehicle, error) {
	return s.vehicleRepo.FindByID(id)
}

func (s *vehicleService) GetByVIN(vin string) (*domain.Vehicle, error) {
	return s.vehicleRepo.FindByVIN(vin)
}

func (s *vehicleService) Update(id uint, inp input.UpdateVehicleInput) (*domain.Vehicle, error) {
	vehicle, err := s.vehicleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("vehicle not found")
	}

	if vehicle.IsSold() {
		return nil, errors.New("cannot update sold vehicle")
	}

	// Actualizar campos
	if inp.Plate != nil {
		vehicle.Plate = *inp.Plate
	}
	if inp.Trim != nil {
		vehicle.Trim = *inp.Trim
	}
	if inp.Mileage != nil {
		vehicle.Mileage = *inp.Mileage
	}
	if inp.ExteriorColor != nil {
		vehicle.ExteriorColor = *inp.ExteriorColor
	}
	if inp.InteriorColor != nil {
		vehicle.InteriorColor = *inp.InteriorColor
	}

	// Actualizar precios
	msrp := vehicle.MSRP
	invoicePrice := vehicle.InvoicePrice
	askingPrice := vehicle.AskingPrice

	if inp.MSRP != nil {
		msrp = *inp.MSRP
	}
	if inp.InvoicePrice != nil {
		invoicePrice = *inp.InvoicePrice
	}
	if inp.AskingPrice != nil {
		askingPrice = *inp.AskingPrice
	}

	if err := vehicle.SetPricing(msrp, invoicePrice, askingPrice); err != nil {
		return nil, err
	}

	vehicle.ModifiedAt = time.Now()

	if err := s.vehicleRepo.Update(vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}

func (s *vehicleService) Delete(id uint) error {
	vehicle, err := s.vehicleRepo.FindByID(id)
	if err != nil {
		return errors.New("vehicle not found")
	}

	if vehicle.IsSold() {
		return errors.New("cannot delete sold vehicle")
	}

	// Eliminar fotos primero
	if err := s.photoRepo.DeleteByVehicleID(id); err != nil {
		return err
	}

	return s.vehicleRepo.Delete(id)
}

func (s *vehicleService) List(limit int, offset int) ([]*domain.Vehicle, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.vehicleRepo.FindAll(limit, offset)
}

func (s *vehicleService) ListAvailable(limit int, offset int) ([]*domain.Vehicle, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.vehicleRepo.FindAvailable(limit, offset)
}

func (s *vehicleService) ListByStatus(status string, limit int, offset int) ([]*domain.Vehicle, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.vehicleRepo.FindByStatus(domain.VehicleStatus(status), limit, offset)
}

func (s *vehicleService) ListByLocation(locationID uint) ([]*domain.Vehicle, error) {
	return s.vehicleRepo.FindByLocationID(locationID)
}

// Status changes

func (s *vehicleService) MarkAsSold(id uint) error {
	vehicle, err := s.vehicleRepo.FindByID(id)
	if err != nil {
		return errors.New("vehicle not found")
	}

	if err := vehicle.MarkAsSold(); err != nil {
		return err
	}

	return s.vehicleRepo.Update(vehicle)
}

func (s *vehicleService) MarkAsReadyForSale(id uint) error {
	vehicle, err := s.vehicleRepo.FindByID(id)
	if err != nil {
		return errors.New("vehicle not found")
	}

	if err := vehicle.MarkAsReadyForSale(); err != nil {
		return err
	}

	return s.vehicleRepo.Update(vehicle)
}

func (s *vehicleService) SendToRecon(id uint) error {
	vehicle, err := s.vehicleRepo.FindByID(id)
	if err != nil {
		return errors.New("vehicle not found")
	}

	if err := vehicle.SendToRecon(); err != nil {
		return err
	}

	return s.vehicleRepo.Update(vehicle)
}

func (s *vehicleService) ChangeLocation(id uint, locationID uint) error {
	vehicle, err := s.vehicleRepo.FindByID(id)
	if err != nil {
		return errors.New("vehicle not found")
	}

	// Verificar que location exista
	locationExists, err := s.locationRepo.Exists(locationID)
	if err != nil {
		return err
	}
	if !locationExists {
		return errors.New("location not found")
	}

	vehicle.SetLocation(locationID)
	return s.vehicleRepo.Update(vehicle)
}

// Photos

func (s *vehicleService) AddPhoto(inp input.AddPhotoInput) (*domain.VehiclePhoto, error) {
	// Verificar que vehicle exista
	exists, err := s.vehicleRepo.Exists(inp.VehicleID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("vehicle not found")
	}

	photo, err := domain.NewVehiclePhoto(
		inp.VehicleID,
		inp.URL,
		domain.PhotoPerspective(inp.Perspective),
		domain.PhotoPurpose(inp.Purpose),
		inp.UploadedBy,
	)
	if err != nil {
		return nil, err
	}

	if inp.IsPrimary {
		// Quitar primary de otras fotos
		existingPhotos, err := s.photoRepo.FindByVehicleID(inp.VehicleID)
		if err != nil {
			return nil, err
		}
		for _, p := range existingPhotos {
			if p.IsPrimary {
				p.RemovePrimary()
				s.photoRepo.Update(p)
			}
		}
		photo.SetAsPrimary()
	}

	if err := s.photoRepo.Save(photo); err != nil {
		return nil, err
	}

	return photo, nil
}

func (s *vehicleService) GetPhotos(vehicleID uint) ([]*domain.VehiclePhoto, error) {
	return s.photoRepo.FindByVehicleID(vehicleID)
}

func (s *vehicleService) SetPrimaryPhoto(vehicleID uint, photoID uint) error {
	photo, err := s.photoRepo.FindByID(photoID)
	if err != nil {
		return errors.New("photo not found")
	}

	if photo.VehicleID != vehicleID {
		return errors.New("photo does not belong to this vehicle")
	}

	// Quitar primary de otras fotos
	existingPhotos, err := s.photoRepo.FindByVehicleID(vehicleID)
	if err != nil {
		return err
	}
	for _, p := range existingPhotos {
		if p.IsPrimary && p.ID != photoID {
			p.RemovePrimary()
			s.photoRepo.Update(p)
		}
	}

	photo.SetAsPrimary()
	return s.photoRepo.Update(photo)
}

func (s *vehicleService) DeletePhoto(photoID uint) error {
	exists, err := s.photoRepo.FindByID(photoID)
	if err != nil {
		return errors.New("photo not found")
	}
	_ = exists

	return s.photoRepo.Delete(photoID)
}