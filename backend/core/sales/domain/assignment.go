package domain

import (
	"errors"
	"time"
)

type AssignmentRole string

const (
	AssignmentRoleSalesperson AssignmentRole = "salesperson"
	AssignmentRoleManager     AssignmentRole = "manager"
	AssignmentRoleFinance     AssignmentRole = "finance"
	AssignmentRoleCloser      AssignmentRole = "closer"
)

type LeadAssignment struct {
	ID         uint
	LeadID     uint
	EntityID   uint
	Role       AssignmentRole
	IsPrimary  bool
	AssignedBy uint
	Active     bool
	CreatedAt  time.Time
}

func NewLeadAssignment(leadID uint, entityID uint, role AssignmentRole, assignedBy uint) (*LeadAssignment, error) {
	if leadID == 0 {
		return nil, errors.New("lead is required")
	}
	if entityID == 0 {
		return nil, errors.New("entity is required")
	}
	if assignedBy == 0 {
		return nil, errors.New("assigned_by is required")
	}

	return &LeadAssignment{
		LeadID:     leadID,
		EntityID:   entityID,
		Role:       role,
		AssignedBy: assignedBy,
		Active:     true,
		CreatedAt:  time.Now(),
	}, nil
}

func (a *LeadAssignment) SetAsPrimary() {
	a.IsPrimary = true
}

func (a *LeadAssignment) RemovePrimary() {
	a.IsPrimary = false
}

func (a *LeadAssignment) Deactivate() {
	a.Active = false
}

func (a *LeadAssignment) Activate() {
	a.Active = true
}