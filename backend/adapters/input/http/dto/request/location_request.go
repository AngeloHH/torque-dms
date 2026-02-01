package request

type CreateLocationRequest struct {
	Name      string  `json:"name" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	Address   string  `json:"address"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Zip       string  `json:"zip"`
	CountryID uint    `json:"country_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Capacity  int     `json:"capacity"`
}

type UpdateLocationRequest struct {
	Name      *string  `json:"name"`
	Address   *string  `json:"address"`
	City      *string  `json:"city"`
	State     *string  `json:"state"`
	Zip       *string  `json:"zip"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Capacity  *int     `json:"capacity"`
}