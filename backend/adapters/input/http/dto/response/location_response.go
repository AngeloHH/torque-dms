package response

import "time"

type LocationResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	Zip       string    `json:"zip"`
	CountryID uint      `json:"country_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Capacity  int       `json:"capacity"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type LocationListResponse struct {
	Locations []LocationResponse `json:"locations"`
	Total     int                `json:"total"`
}