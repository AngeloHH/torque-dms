package input

import "torque-dms/core/inventory/domain"

type CreateVehicleInput struct {
	StockNumber       string
	VIN               string
	Plate             string
	Make              string
	Model             string
	Trim              string
	Year              int
	Mileage           int
	ExteriorColor     string
	InteriorColor     string
	MSRP              float64
	InvoicePrice      float64
	AskingPrice       float64
	Condition         string
	LocationID        uint
	AcquisitionSource string
	AcquisitionCost   float64
}

type UpdateVehicleInput struct {
	Plate         *string
	Trim          *string
	Mileage       *int
	ExteriorColor *string
	InteriorColor *string
	MSRP          *float64
	InvoicePrice  *float64
	AskingPrice   *float64
}

type AddPhotoInput struct {
	VehicleID   uint
	URL         string
	Perspective string
	Purpose     string
	UploadedBy  uint
	IsPrimary   bool
}

type VehicleService interface {
	Create(input CreateVehicleInput) (*domain.Vehicle, error)
	GetByID(id uint) (*domain.Vehicle, error)
	GetByVIN(vin string) (*domain.Vehicle, error)
	Update(id uint, input UpdateVehicleInput) (*domain.Vehicle, error)
	Delete(id uint) error
	List(limit int, offset int) ([]*domain.Vehicle, error)
	ListAvailable(limit int, offset int) ([]*domain.Vehicle, error)
	ListByStatus(status string, limit int, offset int) ([]*domain.Vehicle, error)
	ListByLocation(locationID uint) ([]*domain.Vehicle, error)

	// Status changes
	MarkAsSold(id uint) error
	MarkAsReadyForSale(id uint) error
	SendToRecon(id uint) error
	ChangeLocation(id uint, locationID uint) error

	// Photos
	AddPhoto(input AddPhotoInput) (*domain.VehiclePhoto, error)
	GetPhotos(vehicleID uint) ([]*domain.VehiclePhoto, error)
	SetPrimaryPhoto(vehicleID uint, photoID uint) error
	DeletePhoto(photoID uint) error
}