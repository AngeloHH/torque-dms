package services

import (
	"errors"

	"torque-dms/core/inventory/domain"
	"torque-dms/core/inventory/ports/input"
	"torque-dms/core/inventory/ports/output"
)

type locationService struct {
	locationRepo output.LocationRepository
	vehicleRepo  output.VehicleRepository
}

func NewLocationService(
	locationRepo output.LocationRepository,
	vehicleRepo output.VehicleRepository,
) input.LocationService {
	return &locationService{
		locationRepo: locationRepo,
		vehicleRepo:  vehicleRepo,
	}
}

func (s *locationService) Create(inp input.CreateLocationInput) (*domain.Location, error) {
	location, err := domain.NewLocation(inp.Name, domain.LocationType(inp.Type))
	if err != nil {
		return nil, err
	}

	location.SetAddress(inp.Address, inp.City, inp.State, inp.Zip, inp.CountryID)

	if inp.Latitude != 0 || inp.Longitude != 0 {
		if err := location.SetCoordinates(inp.Latitude, inp.Longitude); err != nil {
			return nil, err
		}
	}

	if inp.Capacity > 0 {
		if err := location.SetCapacity(inp.Capacity); err != nil {
			return nil, err
		}
	}

	if err := s.locationRepo.Save(location); err != nil {
		return nil, err
	}

	return location, nil
}

func (s *locationService) GetByID(id uint) (*domain.Location, error) {
	return s.locationRepo.FindByID(id)
}

func (s *locationService) Update(id uint, inp input.UpdateLocationInput) (*domain.Location, error) {
	location, err := s.locationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("location not found")
	}

	if inp.Name != nil {
		location.Name = *inp.Name
	}
	if inp.Address != nil {
		location.Address = *inp.Address
	}
	if inp.City != nil {
		location.City = *inp.City
	}
	if inp.State != nil {
		location.State = *inp.State
	}
	if inp.Zip != nil {
		location.Zip = *inp.Zip
	}
	if inp.Latitude != nil && inp.Longitude != nil {
		if err := location.SetCoordinates(*inp.Latitude, *inp.Longitude); err != nil {
			return nil, err
		}
	}
	if inp.Capacity != nil {
		if err := location.SetCapacity(*inp.Capacity); err != nil {
			return nil, err
		}
	}

	if err := s.locationRepo.Update(location); err != nil {
		return nil, err
	}

	return location, nil
}

func (s *locationService) Delete(id uint) error {
	// Verificar que no haya vehículos en esta ubicación
	vehicles, err := s.vehicleRepo.FindByLocationID(id)
	if err != nil {
		return err
	}
	if len(vehicles) > 0 {
		return errors.New("cannot delete location with vehicles")
	}

	return s.locationRepo.Delete(id)
}

func (s *locationService) List() ([]*domain.Location, error) {
	return s.locationRepo.FindAll()
}

func (s *locationService) ListByType(locationType string) ([]*domain.Location, error) {
	return s.locationRepo.FindByType(domain.LocationType(locationType))
}

func (s *locationService) ListActive() ([]*domain.Location, error) {
	return s.locationRepo.FindActive()
}

func (s *locationService) Deactivate(id uint) error {
	location, err := s.locationRepo.FindByID(id)
	if err != nil {
		return errors.New("location not found")
	}

	location.Deactivate()
	return s.locationRepo.Update(location)
}

func (s *locationService) Activate(id uint) error {
	location, err := s.locationRepo.FindByID(id)
	if err != nil {
		return errors.New("location not found")
	}

	location.Activate()
	return s.locationRepo.Update(location)
}