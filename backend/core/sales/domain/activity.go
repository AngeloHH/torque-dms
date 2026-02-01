package domain

import (
	"errors"
	"time"
)

type ActivityType string

const (
	ActivityTypeCallOutbound         ActivityType = "call_outbound"
	ActivityTypeCallInbound          ActivityType = "call_inbound"
	ActivityTypeEmailSent            ActivityType = "email_sent"
	ActivityTypeEmailReceived        ActivityType = "email_received"
	ActivityTypeSMSSent              ActivityType = "sms_sent"
	ActivityTypeSMSReceived          ActivityType = "sms_received"
	ActivityTypeAppointmentScheduled ActivityType = "appointment_scheduled"
	ActivityTypeAppointmentCompleted ActivityType = "appointment_completed"
	ActivityTypeAppointmentCancelled ActivityType = "appointment_cancelled"
	ActivityTypeDemo                 ActivityType = "demo"
	ActivityTypeQuoteSent            ActivityType = "quote_sent"
	ActivityTypeOther                ActivityType = "other"
)

type LeadActivity struct {
	ID          uint
	LeadID      uint
	Type        ActivityType
	Description string
	Outcome     string
	PhoneID     *uint
	Email       string
	PerformedBy uint
	ScheduledAt *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
}

func NewLeadActivity(leadID uint, activityType ActivityType, performedBy uint) (*LeadActivity, error) {
	if leadID == 0 {
		return nil, errors.New("lead is required")
	}
	if performedBy == 0 {
		return nil, errors.New("performed_by is required")
	}

	return &LeadActivity{
		LeadID:      leadID,
		Type:        activityType,
		PerformedBy: performedBy,
		CreatedAt:   time.Now(),
	}, nil
}

func (a *LeadActivity) SetDescription(description string) {
	a.Description = description
}

func (a *LeadActivity) SetOutcome(outcome string) {
	a.Outcome = outcome
}

func (a *LeadActivity) SetPhone(phoneID uint) {
	a.PhoneID = &phoneID
}

func (a *LeadActivity) SetEmail(email string) {
	a.Email = email
}

func (a *LeadActivity) Schedule(scheduledAt time.Time) {
	a.ScheduledAt = &scheduledAt
}

func (a *LeadActivity) Complete() {
	now := time.Now()
	a.CompletedAt = &now
}

func (a *LeadActivity) IsCompleted() bool {
	return a.CompletedAt != nil
}

func (a *LeadActivity) IsScheduled() bool {
	return a.ScheduledAt != nil && a.CompletedAt == nil
}

func (a *LeadActivity) IsOverdue() bool {
	if a.ScheduledAt == nil || a.CompletedAt != nil {
		return false
	}
	return time.Now().After(*a.ScheduledAt)
}