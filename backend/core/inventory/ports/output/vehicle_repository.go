package output

import "torque-dms/core/inventory/domain"

type VehicleRepository interface {
	Save(vehicle *domain.Vehicle) error
	Update(vehicle *domain.Vehicle) error
	FindByID(id uint) (*domain.Vehicle, error)
	FindByVIN(vin string) (*domain.Vehicle, error)
	FindByStockNumber(stockNumber string) (*domain.Vehicle, error)
	FindAll(limit int, offset int) ([]*domain.Vehicle, error)
	FindByStatus(status domain.VehicleStatus, limit int, offset int) ([]*domain.Vehicle, error)
	FindByLocationID(locationID uint) ([]*domain.Vehicle, error)
	FindAvailable(limit int, offset int) ([]*domain.Vehicle, error)
	Delete(id uint) error
	Exists(id uint) (bool, error)
	ExistsByVIN(vin string) (bool, error)
}

type VehiclePhotoRepository interface {
	Save(photo *domain.VehiclePhoto) error
	Update(photo *domain.VehiclePhoto) error
	FindByID(id uint) (*domain.VehiclePhoto, error)
	FindByVehicleID(vehicleID uint) ([]*domain.VehiclePhoto, error)
	FindPrimaryByVehicleID(vehicleID uint) (*domain.VehiclePhoto, error)
	Delete(id uint) error
	DeleteByVehicleID(vehicleID uint) error
}