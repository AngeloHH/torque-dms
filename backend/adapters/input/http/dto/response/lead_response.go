package response

import "time"

type LeadResponse struct {
	ID            uint      `json:"id"`
	EntityID      uint      `json:"entity_id"`
	VehicleID     *uint     `json:"vehicle_id"`
	InterestType  string    `json:"interest_type"`
	InterestMake  string    `json:"interest_make"`
	InterestModel string    `json:"interest_model"`
	BudgetMin     float64   `json:"budget_min"`
	BudgetMax     float64   `json:"budget_max"`
	SourceID      uint      `json:"source_id"`
	SourceDetail  string    `json:"source_detail"`
	PresetID      *uint     `json:"preset_id"`
	CreatedAt     time.Time `json:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`
}

type LeadListResponse struct {
	Leads []LeadResponse `json:"leads"`
	Total int            `json:"total"`
}

type LeadSourceResponse struct {
	ID         uint      `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	IsExternal bool      `json:"is_external"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
}

type LeadSourceListResponse struct {
	Sources []LeadSourceResponse `json:"sources"`
	Total   int                  `json:"total"`
}

type LeadAssignmentResponse struct {
	ID         uint      `json:"id"`
	LeadID     uint      `json:"lead_id"`
	EntityID   uint      `json:"entity_id"`
	Role       string    `json:"role"`
	IsPrimary  bool      `json:"is_primary"`
	AssignedBy uint      `json:"assigned_by"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
}

type LeadAssignmentListResponse struct {
	Assignments []LeadAssignmentResponse `json:"assignments"`
	Total       int                      `json:"total"`
}

type LeadNoteResponse struct {
	ID         uint      `json:"id"`
	LeadID     uint      `json:"lead_id"`
	Content    string    `json:"content"`
	CreatedBy  uint      `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type LeadNoteListResponse struct {
	Notes []LeadNoteResponse `json:"notes"`
	Total int                `json:"total"`
}

type LeadActivityResponse struct {
	ID          uint       `json:"id"`
	LeadID      uint       `json:"lead_id"`
	Type        string     `json:"type"`
	Description string     `json:"description"`
	Outcome     string     `json:"outcome"`
	PhoneID     *uint      `json:"phone_id"`
	Email       string     `json:"email"`
	PerformedBy uint       `json:"performed_by"`
	ScheduledAt *time.Time `json:"scheduled_at"`
	CompletedAt *time.Time `json:"completed_at"`
	IsOverdue   bool       `json:"is_overdue"`
	CreatedAt   time.Time  `json:"created_at"`
}

type LeadActivityListResponse struct {
	Activities []LeadActivityResponse `json:"activities"`
	Total      int                    `json:"total"`
}

type LeadStepPresetResponse struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SortOrder   int       `json:"sort_order"`
	IsPublic    bool      `json:"is_public"`
	IsShared    bool      `json:"is_shared"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type LeadStepPresetListResponse struct {
	Presets []LeadStepPresetResponse `json:"presets"`
	Total   int                      `json:"total"`
}

type LeadStepResponse struct {
	ID        uint      `json:"id"`
	PresetID  uint      `json:"preset_id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
	IsFinal   bool      `json:"is_final"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type LeadStepListResponse struct {
	Steps []LeadStepResponse `json:"steps"`
	Total int                `json:"total"`
}

type LeadStepProgressResponse struct {
	ID          uint       `json:"id"`
	LeadID      uint       `json:"lead_id"`
	StepID      uint       `json:"step_id"`
	Status      string     `json:"status"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CompletedBy *uint      `json:"completed_by"`
	Notes       string     `json:"notes"`
	CreatedAt   time.Time  `json:"created_at"`
}

type LeadStepProgressListResponse struct {
	Progress []LeadStepProgressResponse `json:"progress"`
	Total    int                        `json:"total"`
}