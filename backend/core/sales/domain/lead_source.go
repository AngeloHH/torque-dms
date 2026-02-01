package domain

import (
	"errors"
	"time"
)

type LeadSource struct {
	ID         uint
	Code       string
	Name       string
	IsExternal bool
	Active     bool
	CreatedAt  time.Time
}

func NewLeadSource(code string, name string, isExternal bool) (*LeadSource, error) {
	if code == "" {
		return nil, errors.New("code is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	return &LeadSource{
		Code:       code,
		Name:       name,
		IsExternal: isExternal,
		Active:     true,
		CreatedAt:  time.Now(),
	}, nil
}

func (s *LeadSource) Deactivate() {
	s.Active = false
}

func (s *LeadSource) Activate() {
	s.Active = true
}