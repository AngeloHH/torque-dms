package request

type CreateLeadRequest struct {
	EntityID      uint    `json:"entity_id" binding:"required"`
	VehicleID     *uint   `json:"vehicle_id"`
	InterestType  string  `json:"interest_type"`
	InterestMake  string  `json:"interest_make"`
	InterestModel string  `json:"interest_model"`
	BudgetMin     float64 `json:"budget_min"`
	BudgetMax     float64 `json:"budget_max"`
	SourceID      uint    `json:"source_id" binding:"required"`
	SourceDetail  string  `json:"source_detail"`
	PresetID      *uint   `json:"preset_id"`
	AssignedTo    uint    `json:"assigned_to"`
}

type UpdateLeadRequest struct {
	VehicleID     *uint    `json:"vehicle_id"`
	InterestType  *string  `json:"interest_type"`
	InterestMake  *string  `json:"interest_make"`
	InterestModel *string  `json:"interest_model"`
	BudgetMin     *float64 `json:"budget_min"`
	BudgetMax     *float64 `json:"budget_max"`
	SourceDetail  *string  `json:"source_detail"`
}

type CreateLeadSourceRequest struct {
	Code       string `json:"code" binding:"required"`
	Name       string `json:"name" binding:"required"`
	IsExternal bool   `json:"is_external"`
}

type AssignLeadRequest struct {
	EntityID  uint   `json:"entity_id" binding:"required"`
	Role      string `json:"role" binding:"required"`
	IsPrimary bool   `json:"is_primary"`
}

type SetPrimaryAssignmentRequest struct {
	AssignmentID uint `json:"assignment_id" binding:"required"`
}

type AddNoteRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateNoteRequest struct {
	Content string `json:"content" binding:"required"`
}

type AddActivityRequest struct {
	Type        string  `json:"type" binding:"required"`
	Description string  `json:"description"`
	Outcome     string  `json:"outcome"`
	PhoneID     *uint   `json:"phone_id"`
	Email       string  `json:"email"`
	ScheduledAt *string `json:"scheduled_at"`
}

type CreatePresetRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
	IsShared    bool   `json:"is_shared"`
}

type CreateStepRequest struct {
	Code      string `json:"code" binding:"required"`
	Name      string `json:"name" binding:"required"`
	SortOrder int    `json:"sort_order"`
	IsFinal   bool   `json:"is_final"`
}

type UpdateProgressRequest struct {
	Status string `json:"status" binding:"required"`
	Notes  string `json:"notes"`
}