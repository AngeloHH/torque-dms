package request

type CreateVehicleRequest struct {
	StockNumber       string  `json:"stock_number" binding:"required"`
	VIN               string  `json:"vin" binding:"required"`
	Plate             string  `json:"plate"`
	Make              string  `json:"make" binding:"required"`
	Model             string  `json:"model" binding:"required"`
	Trim              string  `json:"trim"`
	Year              int     `json:"year" binding:"required"`
	Mileage           int     `json:"mileage"`
	ExteriorColor     string  `json:"exterior_color"`
	InteriorColor     string  `json:"interior_color"`
	MSRP              float64 `json:"msrp"`
	InvoicePrice      float64 `json:"invoice_price"`
	AskingPrice       float64 `json:"asking_price"`
	Condition         string  `json:"condition"`
	LocationID        uint    `json:"location_id"`
	AcquisitionSource string  `json:"acquisition_source"`
	AcquisitionCost   float64 `json:"acquisition_cost"`
}

type UpdateVehicleRequest struct {
	Plate         *string  `json:"plate"`
	Trim          *string  `json:"trim"`
	Mileage       *int     `json:"mileage"`
	ExteriorColor *string  `json:"exterior_color"`
	InteriorColor *string  `json:"interior_color"`
	MSRP          *float64 `json:"msrp"`
	InvoicePrice  *float64 `json:"invoice_price"`
	AskingPrice   *float64 `json:"asking_price"`
}

type ChangeLocationRequest struct {
	LocationID uint `json:"location_id" binding:"required"`
}

type AddPhotoRequest struct {
	URL         string `json:"url" binding:"required"`
	Perspective string `json:"perspective" binding:"required"`
	Purpose     string `json:"purpose" binding:"required"`
	IsPrimary   bool   `json:"is_primary"`
}

type SetPrimaryPhotoRequest struct {
	PhotoID uint `json:"photo_id" binding:"required"`
}