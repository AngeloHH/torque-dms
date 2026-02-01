package models

import "time"

type VehicleModel3D struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `json:"name"`
	BodyType     BodyType  `json:"body_type"`
	FileURL      string    `json:"file_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Active       bool      `gorm:"default:true" json:"active"`
	CreatedAt    time.Time `json:"created_at"`
}

type VehicleModelZone struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Model3DID uint           `json:"model_3d_id"`
	Model3D   VehicleModel3D `gorm:"foreignKey:Model3DID" json:"-"`
	Code      string         `json:"code"`
	Name      string         `json:"name"`
	MeshID    string         `json:"mesh_id"`
	CreatedAt time.Time      `json:"created_at"`
}

type Vehicle struct {
	ID                uint              `gorm:"primaryKey" json:"id"`
	StockNumber       string            `gorm:"unique" json:"stock_number"`
	VIN               string            `gorm:"unique" json:"vin"`
	Plate             string            `json:"plate"`
	Make              string            `json:"make"`
	Model             string            `json:"model"`
	Trim              string            `json:"trim"`
	Year              int               `json:"year"`
	Mileage           int               `json:"mileage"`
	ExteriorColor     string            `json:"exterior_color"`
	InteriorColor     string            `json:"interior_color"`
	MSRP              float64           `json:"msrp"`
	InvoicePrice      float64           `json:"invoice_price"`
	AskingPrice       float64           `json:"asking_price"`
	Condition         VehicleCondition  `json:"condition"`
	Status            VehicleStatus     `json:"status"`
	LotType           LotType           `json:"lot_type"`
	LocationID        uint              `json:"location_id"`
	Location          Location          `gorm:"foreignKey:LocationID" json:"location"`
	AcquisitionSource AcquisitionSource `json:"acquisition_source"`
	AcquisitionDate   time.Time         `json:"acquisition_date"`
	AcquisitionCost   float64           `json:"acquisition_cost"`
	Model3DID         *uint             `json:"model_3d_id"`
	Model3D           *VehicleModel3D   `gorm:"foreignKey:Model3DID" json:"model_3d,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	ModifiedAt        time.Time         `json:"modified_at"`
}

type VehicleLocationHistory struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	VehicleID      uint      `json:"vehicle_id"`
	Vehicle        Vehicle   `gorm:"foreignKey:VehicleID" json:"-"`
	FromLocationID uint      `json:"from_location_id"`
	FromLocation   Location  `gorm:"foreignKey:FromLocationID" json:"from_location"`
	ToLocationID   uint      `json:"to_location_id"`
	ToLocation     Location  `gorm:"foreignKey:ToLocationID" json:"to_location"`
	MovedBy        uint      `json:"moved_by"`
	Reason         string    `json:"reason"`
	CreatedAt      time.Time `json:"created_at"`
}

type VehicleTracking struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	VehicleID  uint      `json:"vehicle_id"`
	Vehicle    Vehicle   `gorm:"foreignKey:VehicleID" json:"-"`
	Latitude   float64   `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude  float64   `gorm:"type:decimal(11,8)" json:"longitude"`
	RouteID    *uint     `json:"route_id"`
	Route      *Route    `gorm:"foreignKey:RouteID" json:"route,omitempty"`
	RecordedAt time.Time `json:"recorded_at"`
}

type VehiclePhoto struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	VehicleID   uint             `json:"vehicle_id"`
	Vehicle     Vehicle          `gorm:"foreignKey:VehicleID" json:"-"`
	URL         string           `json:"url"`
	Perspective PhotoPerspective `json:"perspective"`
	Purpose     PhotoPurpose     `json:"purpose"`
	SortOrder   int              `gorm:"default:0" json:"sort_order"`
	IsPrimary   bool             `gorm:"default:false" json:"is_primary"`
	UploadedBy  uint             `json:"uploaded_by"`
	CreatedAt   time.Time        `json:"created_at"`
}

type VehicleZoneMark struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	VehicleID   uint             `json:"vehicle_id"`
	Vehicle     Vehicle          `gorm:"foreignKey:VehicleID" json:"-"`
	ZoneID      uint             `json:"zone_id"`
	Zone        VehicleModelZone `gorm:"foreignKey:ZoneID" json:"zone"`
	Type        string           `json:"type"`
	Severity    string           `json:"severity"`
	Description string           `json:"description"`
	PhotoID     *uint            `json:"photo_id"`
	Photo       *VehiclePhoto    `gorm:"foreignKey:PhotoID" json:"photo,omitempty"`
	ReportedBy  uint             `json:"reported_by"`
	Resolved    bool             `gorm:"default:false" json:"resolved"`
	ResolvedBy  *uint            `json:"resolved_by"`
	ResolvedAt  *time.Time       `json:"resolved_at"`
	CreatedAt   time.Time        `json:"created_at"`
}