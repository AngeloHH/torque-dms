package output

import "torque-dms/core/inventory/domain"

type LocationRepository interface {
	Save(location *domain.Location) error
	Update(location *domain.Location) error
	FindByID(id uint) (*domain.Location, error)
	FindByName(name string) (*domain.Location, error)
	FindAll() ([]*domain.Location, error)
	FindByType(locationType domain.LocationType) ([]*domain.Location, error)
	FindActive() ([]*domain.Location, error)
	Delete(id uint) error
	Exists(id uint) (bool, error)
}