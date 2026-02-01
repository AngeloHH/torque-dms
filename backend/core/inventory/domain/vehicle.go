package domain

import (
	"errors"
	"time"
)

type VehicleCondition string

const (
	VehicleConditionNew       VehicleCondition = "new"
	VehicleConditionUsed      VehicleCondition = "used"
	VehicleConditionCertified VehicleCondition = "certified"
)

type VehicleStatus string

const (
	VehicleStatusInTransit    VehicleStatus = "in_transit"
	VehicleStatusInRecon      VehicleStatus = "in_recon"
	VehicleStatusReadyForSale VehicleStatus = "ready_for_sale"
	VehicleStatusSold         VehicleStatus = "sold"
	VehicleStatusWholesale    VehicleStatus = "wholesale"
)

type LotType string

const (
	LotTypeNew       LotType = "new"
	LotTypeUsed      LotType = "used"
	LotTypeCPO       LotType = "cpo"
	LotTypeWholesale LotType = "wholesale"
)

type AcquisitionSource string

const (
	AcquisitionSourceFactory        AcquisitionSource = "factory"
	AcquisitionSourceTradeIn        AcquisitionSource = "trade_in"
	AcquisitionSourceAuction        AcquisitionSource = "auction"
	AcquisitionSourceDealerTransfer AcquisitionSource = "dealer_transfer"
	AcquisitionSourceConsignment    AcquisitionSource = "consignment"
)

type Vehicle struct {
	ID                uint
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
	Condition         VehicleCondition
	Status            VehicleStatus
	LotType           LotType
	LocationID        uint
	AcquisitionSource AcquisitionSource
	AcquisitionDate   time.Time
	AcquisitionCost   float64
	Model3DID         *uint
	CreatedAt         time.Time
	ModifiedAt        time.Time
}

func NewVehicle(stockNumber string, vin string, make string, model string, year int) (*Vehicle, error) {
	if stockNumber == "" {
		return nil, errors.New("stock number is required")
	}
	if vin == "" {
		return nil, errors.New("VIN is required")
	}
	if len(vin) != 17 {
		return nil, errors.New("VIN must be 17 characters")
	}
	if make == "" {
		return nil, errors.New("make is required")
	}
	if model == "" {
		return nil, errors.New("model is required")
	}
	if year < 1900 || year > time.Now().Year()+2 {
		return nil, errors.New("invalid year")
	}

	return &Vehicle{
		StockNumber: stockNumber,
		VIN:         vin,
		Make:        make,
		Model:       model,
		Year:        year,
		Condition:   VehicleConditionUsed,
		Status:      VehicleStatusInTransit,
		LotType:     LotTypeUsed,
		CreatedAt:   time.Now(),
	}, nil
}

func (v *Vehicle) SetPricing(msrp float64, invoicePrice float64, askingPrice float64) error {
	if msrp < 0 || invoicePrice < 0 || askingPrice < 0 {
		return errors.New("prices cannot be negative")
	}
	v.MSRP = msrp
	v.InvoicePrice = invoicePrice
	v.AskingPrice = askingPrice
	v.ModifiedAt = time.Now()
	return nil
}

func (v *Vehicle) SetAcquisition(source AcquisitionSource, cost float64, date time.Time) error {
	if cost < 0 {
		return errors.New("acquisition cost cannot be negative")
	}
	v.AcquisitionSource = source
	v.AcquisitionCost = cost
	v.AcquisitionDate = date
	v.ModifiedAt = time.Now()
	return nil
}

func (v *Vehicle) SetCondition(condition VehicleCondition) {
	v.Condition = condition
	if condition == VehicleConditionNew {
		v.LotType = LotTypeNew
	}
	v.ModifiedAt = time.Now()
}

func (v *Vehicle) SetStatus(status VehicleStatus) {
	v.Status = status
	v.ModifiedAt = time.Now()
}

func (v *Vehicle) SetLocation(locationID uint) {
	v.LocationID = locationID
	v.ModifiedAt = time.Now()
}

func (v *Vehicle) MarkAsSold() error {
	if v.Status == VehicleStatusSold {
		return errors.New("vehicle is already sold")
	}
	v.Status = VehicleStatusSold
	v.ModifiedAt = time.Now()
	return nil
}

func (v *Vehicle) MarkAsReadyForSale() error {
	if v.Status == VehicleStatusSold {
		return errors.New("cannot change status of sold vehicle")
	}
	v.Status = VehicleStatusReadyForSale
	v.ModifiedAt = time.Now()
	return nil
}

func (v *Vehicle) SendToRecon() error {
	if v.Status == VehicleStatusSold {
		return errors.New("cannot send sold vehicle to recon")
	}
	v.Status = VehicleStatusInRecon
	v.ModifiedAt = time.Now()
	return nil
}

func (v *Vehicle) IsAvailable() bool {
	return v.Status == VehicleStatusReadyForSale
}

func (v *Vehicle) IsSold() bool {
	return v.Status == VehicleStatusSold
}

func (v *Vehicle) Profit() float64 {
	return v.AskingPrice - v.AcquisitionCost
}

func (v *Vehicle) Margin() float64 {
	if v.AskingPrice == 0 {
		return 0
	}
	return (v.Profit() / v.AskingPrice) * 100
}