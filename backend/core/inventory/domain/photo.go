package domain

import (
	"errors"
	"time"
)

type PhotoPerspective string

const (
	PhotoPerspectiveFront         PhotoPerspective = "front"
	PhotoPerspectiveRear          PhotoPerspective = "rear"
	PhotoPerspectiveLeftSide      PhotoPerspective = "left_side"
	PhotoPerspectiveRightSide     PhotoPerspective = "right_side"
	PhotoPerspectiveInterior      PhotoPerspective = "interior"
	PhotoPerspectiveDashboard     PhotoPerspective = "dashboard"
	PhotoPerspectiveEngine        PhotoPerspective = "engine"
	PhotoPerspectiveDamage        PhotoPerspective = "damage"
)

type PhotoPurpose string

const (
	PhotoPurposeListing    PhotoPurpose = "listing"
	PhotoPurposeService    PhotoPurpose = "service"
	PhotoPurposeInspection PhotoPurpose = "inspection"
	PhotoPurposeDamage     PhotoPurpose = "damage"
	PhotoPurposeInternal   PhotoPurpose = "internal"
)

type VehiclePhoto struct {
	ID          uint
	VehicleID   uint
	URL         string
	Perspective PhotoPerspective
	Purpose     PhotoPurpose
	SortOrder   int
	IsPrimary   bool
	UploadedBy  uint
	CreatedAt   time.Time
}

func NewVehiclePhoto(vehicleID uint, url string, perspective PhotoPerspective, purpose PhotoPurpose, uploadedBy uint) (*VehiclePhoto, error) {
	if vehicleID == 0 {
		return nil, errors.New("vehicle is required")
	}
	if url == "" {
		return nil, errors.New("url is required")
	}
	if uploadedBy == 0 {
		return nil, errors.New("uploaded_by is required")
	}

	return &VehiclePhoto{
		VehicleID:   vehicleID,
		URL:         url,
		Perspective: perspective,
		Purpose:     purpose,
		UploadedBy:  uploadedBy,
		CreatedAt:   time.Now(),
	}, nil
}

func (p *VehiclePhoto) SetAsPrimary() {
	p.IsPrimary = true
}

func (p *VehiclePhoto) RemovePrimary() {
	p.IsPrimary = false
}

func (p *VehiclePhoto) SetSortOrder(order int) {
	p.SortOrder = order
}