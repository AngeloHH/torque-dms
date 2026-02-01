package input

import "torque-dms/core/inventory/domain"

type CreateLocationInput struct {
	Name      string
	Type      string
	Address   string
	City      string
	State     string
	Zip       string
	CountryID uint
	Latitude  float64
	Longitude float64
	Capacity  int
}

type UpdateLocationInput struct {
	Name      *string
	Address   *string
	City      *string
	State     *string
	Zip       *string
	Latitude  *float64
	Longitude *float64
	Capacity  *int
}

type LocationService interface {
	Create(input CreateLocationInput) (*domain.Location, error)
	GetByID(id uint) (*domain.Location, error)
	Update(id uint, input UpdateLocationInput) (*domain.Location, error)
	Delete(id uint) error
	List() ([]*domain.Location, error)
	ListByType(locationType string) ([]*domain.Location, error)
	ListActive() ([]*domain.Location, error)
	Deactivate(id uint) error
	Activate(id uint) error
}