package domain

import (
	"errors"
	"time"
)

type Lead struct {
	ID            uint
	EntityID      uint
	VehicleID     *uint
	InterestType  string
	InterestMake  string
	InterestModel string
	BudgetMin     float64
	BudgetMax     float64
	SourceID      uint
	SourceDetail  string
	PresetID      *uint
	CreatedAt     time.Time
	ModifiedAt    time.Time
}

func NewLead(entityID uint, sourceID uint) (*Lead, error) {
	if entityID == 0 {
		return nil, errors.New("entity is required")
	}
	if sourceID == 0 {
		return nil, errors.New("source is required")
	}

	return &Lead{
		EntityID:  entityID,
		SourceID:  sourceID,
		CreatedAt: time.Now(),
	}, nil
}

func (l *Lead) SetInterest(interestType string, make string, model string) {
	l.InterestType = interestType
	l.InterestMake = make
	l.InterestModel = model
	l.ModifiedAt = time.Now()
}

func (l *Lead) SetBudget(min float64, max float64) error {
	if min < 0 || max < 0 {
		return errors.New("budget cannot be negative")
	}
	if min > max && max > 0 {
		return errors.New("min budget cannot exceed max budget")
	}
	l.BudgetMin = min
	l.BudgetMax = max
	l.ModifiedAt = time.Now()
	return nil
}

func (l *Lead) SetVehicle(vehicleID uint) {
	l.VehicleID = &vehicleID
	l.ModifiedAt = time.Now()
}

func (l *Lead) RemoveVehicle() {
	l.VehicleID = nil
	l.ModifiedAt = time.Now()
}

func (l *Lead) SetPreset(presetID uint) {
	l.PresetID = &presetID
	l.ModifiedAt = time.Now()
}

func (l *Lead) IsHighValue() bool {
	return l.BudgetMax >= 50000
}

func (l *Lead) HasVehicleInterest() bool {
	return l.VehicleID != nil
}