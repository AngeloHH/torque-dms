package response

import "time"

type VehicleResponse struct {
	ID                uint      `json:"id"`
	StockNumber       string    `json:"stock_number"`
	VIN               string    `json:"vin"`
	Plate             string    `json:"plate"`
	Make              string    `json:"make"`
	Model             string    `json:"model"`
	Trim              string    `json:"trim"`
	Year              int       `json:"year"`
	Mileage           int       `json:"mileage"`
	ExteriorColor     string    `json:"exterior_color"`
	InteriorColor     string    `json:"interior_color"`
	MSRP              float64   `json:"msrp"`
	InvoicePrice      float64   `json:"invoice_price"`
	AskingPrice       float64   `json:"asking_price"`
	Condition         string    `json:"condition"`
	Status            string    `json:"status"`
	LotType           string    `json:"lot_type"`
	LocationID        uint      `json:"location_id"`
	AcquisitionSource string    `json:"acquisition_source"`
	AcquisitionDate   time.Time `json:"acquisition_date"`
	AcquisitionCost   float64   `json:"acquisition_cost"`
	Profit            float64   `json:"profit"`
	Margin            float64   `json:"margin"`
	CreatedAt         time.Time `json:"created_at"`
	ModifiedAt        time.Time `json:"modified_at"`
}

type VehicleListResponse struct {
	Vehicles []VehicleResponse `json:"vehicles"`
	Total    int               `json:"total"`
}

type VehiclePhotoResponse struct {
	ID          uint      `json:"id"`
	VehicleID   uint      `json:"vehicle_id"`
	URL         string    `json:"url"`
	Perspective string    `json:"perspective"`
	Purpose     string    `json:"purpose"`
	SortOrder   int       `json:"sort_order"`
	IsPrimary   bool      `json:"is_primary"`
	CreatedAt   time.Time `json:"created_at"`
}

type VehiclePhotosResponse struct {
	Photos []VehiclePhotoResponse `json:"photos"`
	Total  int                    `json:"total"`
}