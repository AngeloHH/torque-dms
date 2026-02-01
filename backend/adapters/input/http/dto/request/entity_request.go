package request

type CreateEntityRequest struct {
	Type         string `json:"type" binding:"required"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	BusinessName string `json:"business_name"`
	TaxID        string `json:"tax_id"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
	CountryID    *uint  `json:"country_id"`
}

type UpdateEntityRequest struct {
	Field string `json:"field" binding:"required"`
	Value string `json:"value" binding:"required"`
}