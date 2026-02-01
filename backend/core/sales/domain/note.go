package domain

import (
	"errors"
	"time"
)

type LeadNote struct {
	ID         uint
	LeadID     uint
	Content    string
	CreatedBy  uint
	CreatedAt  time.Time
	ModifiedAt time.Time
}

func NewLeadNote(leadID uint, content string, createdBy uint) (*LeadNote, error) {
	if leadID == 0 {
		return nil, errors.New("lead is required")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}
	if createdBy == 0 {
		return nil, errors.New("created_by is required")
	}

	now := time.Now()
	return &LeadNote{
		LeadID:     leadID,
		Content:    content,
		CreatedBy:  createdBy,
		CreatedAt:  now,
		ModifiedAt: now,
	}, nil
}

func (n *LeadNote) Update(content string) error {
	if content == "" {
		return errors.New("content is required")
	}
	n.Content = content
	n.ModifiedAt = time.Now()
	return nil
}