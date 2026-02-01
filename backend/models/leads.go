package models

import "time"

type LeadSource struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Code       string    `gorm:"unique" json:"code"`
	Name       string    `json:"name"`
	IsExternal bool      `gorm:"default:false" json:"is_external"`
	Active     bool      `gorm:"default:true" json:"active"`
	CreatedAt  time.Time `json:"created_at"`
}

type Lead struct {
	ID            uint             `gorm:"primaryKey" json:"id"`
	EntityID      uint             `json:"entity_id"`
	Entity        Entity           `gorm:"foreignKey:EntityID" json:"entity"`
	VehicleID     *uint            `json:"vehicle_id"`
	Vehicle       *Vehicle         `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	InterestType  VehicleCondition `json:"interest_type"`
	InterestMake  string           `json:"interest_make"`
	InterestModel string           `json:"interest_model"`
	BudgetMin     float64          `json:"budget_min"`
	BudgetMax     float64          `json:"budget_max"`
	SourceID      uint             `json:"source_id"`
	Source        LeadSource       `gorm:"foreignKey:SourceID" json:"source"`
	SourceDetail  string           `json:"source_detail"`
	PresetID      *uint            `json:"preset_id"`
	Preset        *LeadStepPreset  `gorm:"foreignKey:PresetID" json:"preset,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	ModifiedAt    time.Time        `json:"modified_at"`
}

type LeadStepPreset struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"unique" json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SortOrder   int       `json:"sort_order"`
	IsPublic    bool      `gorm:"default:false" json:"is_public"`
	IsShared    bool      `gorm:"default:false" json:"is_shared"`
	CreatedBy   uint      `json:"created_by"`
	Creator     Entity    `gorm:"foreignKey:CreatedBy" json:"creator"`
	CreatedAt   time.Time `json:"created_at"`
}

type LeadStep struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PresetID  uint           `json:"preset_id"`
	Preset    LeadStepPreset `gorm:"foreignKey:PresetID" json:"-"`
	Code      string         `json:"code"`
	Name      string         `json:"name"`
	SortOrder int            `json:"sort_order"`
	IsFinal   bool           `gorm:"default:false" json:"is_final"`
	Active    bool           `gorm:"default:true" json:"active"`
	CreatedAt time.Time      `json:"created_at"`
}

type LeadStepProgress struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	LeadID      uint       `json:"lead_id"`
	Lead        Lead       `gorm:"foreignKey:LeadID" json:"-"`
	StepID      uint       `json:"step_id"`
	Step        LeadStep   `gorm:"foreignKey:StepID" json:"step"`
	Status      StepStatus `gorm:"default:'pending'" json:"status"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CompletedBy *uint      `json:"completed_by"`
	Notes       string     `json:"notes"`
	CreatedAt   time.Time  `json:"created_at"`
}

type LeadAssignment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	LeadID     uint      `json:"lead_id"`
	Lead       Lead      `gorm:"foreignKey:LeadID" json:"-"`
	EntityID   uint      `json:"entity_id"`
	Entity     Entity    `gorm:"foreignKey:EntityID" json:"entity"`
	Role       string    `json:"role"`
	IsPrimary  bool      `gorm:"default:false" json:"is_primary"`
	AssignedBy uint      `json:"assigned_by"`
	Active     bool      `gorm:"default:true" json:"active"`
	CreatedAt  time.Time `json:"created_at"`
}

type LeadNote struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	LeadID     uint      `json:"lead_id"`
	Lead       Lead      `gorm:"foreignKey:LeadID" json:"-"`
	Content    string    `gorm:"type:text" json:"content"`
	CreatedBy  uint      `json:"created_by"`
	Creator    Entity    `gorm:"foreignKey:CreatedBy" json:"creator"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type LeadActivity struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	LeadID      uint         `json:"lead_id"`
	Lead        Lead         `gorm:"foreignKey:LeadID" json:"-"`
	Type        ActivityType `json:"type"`
	Description string       `json:"description"`
	Outcome     string       `json:"outcome"`
	PhoneID     *uint        `json:"phone_id"`
	Phone       *EntityPhone `gorm:"foreignKey:PhoneID" json:"phone,omitempty"`
	Email       string       `json:"email"`
	PerformedBy uint         `json:"performed_by"`
	Performer   Entity       `gorm:"foreignKey:PerformedBy" json:"performer"`
	ScheduledAt *time.Time   `json:"scheduled_at"`
	CompletedAt *time.Time   `json:"completed_at"`
	CreatedAt   time.Time    `json:"created_at"`
}