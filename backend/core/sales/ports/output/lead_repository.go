package output

import "torque-dms/core/sales/domain"

type LeadRepository interface {
	Save(lead *domain.Lead) error
	Update(lead *domain.Lead) error
	FindByID(id uint) (*domain.Lead, error)
	FindByEntityID(entityID uint) ([]*domain.Lead, error)
	FindAll(limit int, offset int) ([]*domain.Lead, error)
	Delete(id uint) error
	Exists(id uint) (bool, error)
}

type LeadSourceRepository interface {
	Save(source *domain.LeadSource) error
	Update(source *domain.LeadSource) error
	FindByID(id uint) (*domain.LeadSource, error)
	FindByCode(code string) (*domain.LeadSource, error)
	FindAll() ([]*domain.LeadSource, error)
	FindActive() ([]*domain.LeadSource, error)
	Delete(id uint) error
	Exists(id uint) (bool, error)
}

type LeadAssignmentRepository interface {
	Save(assignment *domain.LeadAssignment) error
	Update(assignment *domain.LeadAssignment) error
	FindByID(id uint) (*domain.LeadAssignment, error)
	FindByLeadID(leadID uint) ([]*domain.LeadAssignment, error)
	FindPrimaryByLeadID(leadID uint) (*domain.LeadAssignment, error)
	FindByEntityID(entityID uint) ([]*domain.LeadAssignment, error)
	Delete(id uint) error
}

type LeadNoteRepository interface {
	Save(note *domain.LeadNote) error
	Update(note *domain.LeadNote) error
	FindByID(id uint) (*domain.LeadNote, error)
	FindByLeadID(leadID uint) ([]*domain.LeadNote, error)
	Delete(id uint) error
}

type LeadActivityRepository interface {
	Save(activity *domain.LeadActivity) error
	Update(activity *domain.LeadActivity) error
	FindByID(id uint) (*domain.LeadActivity, error)
	FindByLeadID(leadID uint) ([]*domain.LeadActivity, error)
	FindScheduledByEntityID(entityID uint) ([]*domain.LeadActivity, error)
	FindOverdue() ([]*domain.LeadActivity, error)
	Delete(id uint) error
}