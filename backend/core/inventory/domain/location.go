package domain

import (
	"errors"
	"time"
)

type LocationType string

const (
	LocationTypeSalesLot  LocationType = "sales_lot"
	LocationTypeStorage   LocationType = "storage"
	LocationTypeService   LocationType = "service"
	LocationTypeOffsite   LocationType = "offsite"
	LocationTypeInTransit LocationType = "in_transit"
)

type Location struct {
	ID        uint
	Name      string
	Type      LocationType
	Address   string
	City      string
	State     string
	Zip       string
	CountryID uint
	Latitude  float64
	Longitude float64
	Capacity  int
	Active    bool
	CreatedAt time.Time
}

func NewLocation(name string, locationType LocationType) (*Location, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	return &Location{
		Name:      name,
		Type:      locationType,
		Active:    true,
		CreatedAt: time.Now(),
	}, nil
}

func (l *Location) SetAddress(address string, city string, state string, zip string, countryID uint) {
	l.Address = address
	l.City = city
	l.State = state
	l.Zip = zip
	l.CountryID = countryID
}

func (l *Location) SetCoordinates(latitude float64, longitude float64) error {
	if latitude < -90 || latitude > 90 {
		return errors.New("invalid latitude")
	}
	if longitude < -180 || longitude > 180 {
		return errors.New("invalid longitude")
	}
	l.Latitude = latitude
	l.Longitude = longitude
	return nil
}

func (l *Location) SetCapacity(capacity int) error {
	if capacity < 0 {
		return errors.New("capacity cannot be negative")
	}
	l.Capacity = capacity
	return nil
}

func (l *Location) Deactivate() {
	l.Active = false
}

func (l *Location) Activate() {
	l.Active = true
}

func (l *Location) IsSalesLot() bool {
	return l.Type == LocationTypeSalesLot
}