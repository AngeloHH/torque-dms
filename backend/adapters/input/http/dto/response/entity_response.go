package response

import "time"

type EntityResponse struct {
	ID             uint            `json:"id"`
	Type           string          `json:"type"`
	FirstName      string          `json:"first_name"`
	LastName       string          `json:"last_name"`
	BusinessName   string          `json:"business_name"`
	TaxID          string          `json:"tax_id"`
	Email          string          `json:"email"`
	Address        string          `json:"address"`
	City           string          `json:"city"`
	State          string          `json:"state"`
	Zip            string          `json:"zip"`
	CountryID      *uint           `json:"country_id"`
	IsSystemUser   bool            `json:"is_system_user"`
	IsInternal     bool            `json:"is_internal"`
	ParentEntityID *uint           `json:"parent_entity_id"`
	Status         string          `json:"status"`
	CreatedAt      time.Time       `json:"created_at"`
	ModifiedAt     time.Time       `json:"modified_at"`
}

type EntityListResponse struct {
	Entities []EntityResponse `json:"entities"`
	Total    int              `json:"total"`
}