package input

import "torque-dms/core/sales/domain"

type CreateLeadInput struct {
	EntityID      uint
	VehicleID     *uint
	InterestType  string
	InterestMake  string
	InterestModel string
	BudgetMin     float64
	BudgetMax     float64
	SourceID      uint
	SourceDetail  string
	PresetID      *uint
	AssignedTo    uint
}

type UpdateLeadInput struct {
	VehicleID     *uint
	InterestType  *string
	InterestMake  *string
	InterestModel *string
	BudgetMin     *float64
	BudgetMax     *float64
	SourceDetail  *string
}

type CreateLeadSourceInput struct {
	Code       string
	Name       string
	IsExternal bool
}

type AssignLeadInput struct {
	LeadID     uint
	EntityID   uint
	Role       string
	IsPrimary  bool
	AssignedBy uint
}

type AddNoteInput struct {
	LeadID    uint
	Content   string
	CreatedBy uint
}

type AddActivityInput struct {
	LeadID      uint
	Type        string
	Description string
	Outcome     string
	PhoneID     *uint
	Email       string
	PerformedBy uint
	ScheduledAt *string
}

type LeadService interface {
	// Lead CRUD
	Create(input CreateLeadInput) (*domain.Lead, error)
	GetByID(id uint) (*domain.Lead, error)
	Update(id uint, input UpdateLeadInput) (*domain.Lead, error)
	Delete(id uint) error
	List(limit int, offset int) ([]*domain.Lead, error)
	ListByEntity(entityID uint) ([]*domain.Lead, error)

	// Lead Sources
	CreateSource(input CreateLeadSourceInput) (*domain.LeadSource, error)
	GetSources() ([]*domain.LeadSource, error)
	GetActiveSources() ([]*domain.LeadSource, error)
	DeactivateSource(id uint) error
	ActivateSource(id uint) error

	// Assignments
	Assign(input AssignLeadInput) (*domain.LeadAssignment, error)
	GetAssignments(leadID uint) ([]*domain.LeadAssignment, error)
	RemoveAssignment(assignmentID uint) error
	SetPrimaryAssignment(leadID uint, assignmentID uint) error

	// Notes
	AddNote(input AddNoteInput) (*domain.LeadNote, error)
	GetNotes(leadID uint) ([]*domain.LeadNote, error)
	UpdateNote(noteID uint, content string) (*domain.LeadNote, error)
	DeleteNote(noteID uint) error

	// Activities
	AddActivity(input AddActivityInput) (*domain.LeadActivity, error)
	GetActivities(leadID uint) ([]*domain.LeadActivity, error)
	CompleteActivity(activityID uint) error
	GetScheduledActivities(entityID uint) ([]*domain.LeadActivity, error)
	GetOverdueActivities() ([]*domain.LeadActivity, error)
}