package services

import (
	"errors"
	"time"

	"torque-dms/core/sales/domain"
	"torque-dms/core/sales/ports/input"
	"torque-dms/core/sales/ports/output"
)

type leadService struct {
	leadRepo       output.LeadRepository
	sourceRepo     output.LeadSourceRepository
	assignmentRepo output.LeadAssignmentRepository
	noteRepo       output.LeadNoteRepository
	activityRepo   output.LeadActivityRepository
}

func NewLeadService(
	leadRepo output.LeadRepository,
	sourceRepo output.LeadSourceRepository,
	assignmentRepo output.LeadAssignmentRepository,
	noteRepo output.LeadNoteRepository,
	activityRepo output.LeadActivityRepository,
) input.LeadService {
	return &leadService{
		leadRepo:       leadRepo,
		sourceRepo:     sourceRepo,
		assignmentRepo: assignmentRepo,
		noteRepo:       noteRepo,
		activityRepo:   activityRepo,
	}
}

// Lead CRUD

func (s *leadService) Create(inp input.CreateLeadInput) (*domain.Lead, error) {
	// Verificar que source exista
	sourceExists, err := s.sourceRepo.Exists(inp.SourceID)
	if err != nil {
		return nil, err
	}
	if !sourceExists {
		return nil, errors.New("source not found")
	}

	lead, err := domain.NewLead(inp.EntityID, inp.SourceID)
	if err != nil {
		return nil, err
	}

	if inp.VehicleID != nil {
		lead.SetVehicle(*inp.VehicleID)
	}

	lead.SetInterest(inp.InterestType, inp.InterestMake, inp.InterestModel)

	if inp.BudgetMin > 0 || inp.BudgetMax > 0 {
		if err := lead.SetBudget(inp.BudgetMin, inp.BudgetMax); err != nil {
			return nil, err
		}
	}

	lead.SourceDetail = inp.SourceDetail

	if inp.PresetID != nil {
		lead.SetPreset(*inp.PresetID)
	}

	if err := s.leadRepo.Save(lead); err != nil {
		return nil, err
	}

	// Si hay asignación inicial
	if inp.AssignedTo > 0 {
		assignment, err := domain.NewLeadAssignment(
			lead.ID,
			inp.AssignedTo,
			domain.AssignmentRoleSalesperson,
			inp.AssignedTo,
		)
		if err != nil {
			return nil, err
		}
		assignment.SetAsPrimary()
		if err := s.assignmentRepo.Save(assignment); err != nil {
			return nil, err
		}
	}

	return lead, nil
}

func (s *leadService) GetByID(id uint) (*domain.Lead, error) {
	return s.leadRepo.FindByID(id)
}

func (s *leadService) Update(id uint, inp input.UpdateLeadInput) (*domain.Lead, error) {
	lead, err := s.leadRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("lead not found")
	}

	if inp.VehicleID != nil {
		if *inp.VehicleID == 0 {
			lead.RemoveVehicle()
		} else {
			lead.SetVehicle(*inp.VehicleID)
		}
	}

	interestType := lead.InterestType
	interestMake := lead.InterestMake
	interestModel := lead.InterestModel

	if inp.InterestType != nil {
		interestType = *inp.InterestType
	}
	if inp.InterestMake != nil {
		interestMake = *inp.InterestMake
	}
	if inp.InterestModel != nil {
		interestModel = *inp.InterestModel
	}
	lead.SetInterest(interestType, interestMake, interestModel)

	if inp.BudgetMin != nil || inp.BudgetMax != nil {
		budgetMin := lead.BudgetMin
		budgetMax := lead.BudgetMax
		if inp.BudgetMin != nil {
			budgetMin = *inp.BudgetMin
		}
		if inp.BudgetMax != nil {
			budgetMax = *inp.BudgetMax
		}
		if err := lead.SetBudget(budgetMin, budgetMax); err != nil {
			return nil, err
		}
	}

	if inp.SourceDetail != nil {
		lead.SourceDetail = *inp.SourceDetail
	}

	lead.ModifiedAt = time.Now()

	if err := s.leadRepo.Update(lead); err != nil {
		return nil, err
	}

	return lead, nil
}

func (s *leadService) Delete(id uint) error {
	exists, err := s.leadRepo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("lead not found")
	}

	return s.leadRepo.Delete(id)
}

func (s *leadService) List(limit int, offset int) ([]*domain.Lead, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.leadRepo.FindAll(limit, offset)
}

func (s *leadService) ListByEntity(entityID uint) ([]*domain.Lead, error) {
	return s.leadRepo.FindByEntityID(entityID)
}

// Lead Sources

func (s *leadService) CreateSource(inp input.CreateLeadSourceInput) (*domain.LeadSource, error) {
	source, err := domain.NewLeadSource(inp.Code, inp.Name, inp.IsExternal)
	if err != nil {
		return nil, err
	}

	if err := s.sourceRepo.Save(source); err != nil {
		return nil, err
	}

	return source, nil
}

func (s *leadService) GetSources() ([]*domain.LeadSource, error) {
	return s.sourceRepo.FindAll()
}

func (s *leadService) GetActiveSources() ([]*domain.LeadSource, error) {
	return s.sourceRepo.FindActive()
}

func (s *leadService) DeactivateSource(id uint) error {
	source, err := s.sourceRepo.FindByID(id)
	if err != nil {
		return errors.New("source not found")
	}

	source.Deactivate()
	return s.sourceRepo.Update(source)
}

func (s *leadService) ActivateSource(id uint) error {
	source, err := s.sourceRepo.FindByID(id)
	if err != nil {
		return errors.New("source not found")
	}

	source.Activate()
	return s.sourceRepo.Update(source)
}

// Assignments

func (s *leadService) Assign(inp input.AssignLeadInput) (*domain.LeadAssignment, error) {
	// Verificar que lead exista
	exists, err := s.leadRepo.Exists(inp.LeadID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("lead not found")
	}

	assignment, err := domain.NewLeadAssignment(
		inp.LeadID,
		inp.EntityID,
		domain.AssignmentRole(inp.Role),
		inp.AssignedBy,
	)
	if err != nil {
		return nil, err
	}

	if inp.IsPrimary {
		// Quitar primary de otros
		existingAssignments, err := s.assignmentRepo.FindByLeadID(inp.LeadID)
		if err != nil {
			return nil, err
		}
		for _, a := range existingAssignments {
			if a.IsPrimary {
				a.RemovePrimary()
				s.assignmentRepo.Update(a)
			}
		}
		assignment.SetAsPrimary()
	}

	if err := s.assignmentRepo.Save(assignment); err != nil {
		return nil, err
	}

	return assignment, nil
}

func (s *leadService) GetAssignments(leadID uint) ([]*domain.LeadAssignment, error) {
	return s.assignmentRepo.FindByLeadID(leadID)
}

func (s *leadService) RemoveAssignment(assignmentID uint) error {
	assignment, err := s.assignmentRepo.FindByID(assignmentID)
	if err != nil {
		return errors.New("assignment not found")
	}

	assignment.Deactivate()
	return s.assignmentRepo.Update(assignment)
}

func (s *leadService) SetPrimaryAssignment(leadID uint, assignmentID uint) error {
	assignment, err := s.assignmentRepo.FindByID(assignmentID)
	if err != nil {
		return errors.New("assignment not found")
	}

	if assignment.LeadID != leadID {
		return errors.New("assignment does not belong to this lead")
	}

	// Quitar primary de otros
	existingAssignments, err := s.assignmentRepo.FindByLeadID(leadID)
	if err != nil {
		return err
	}
	for _, a := range existingAssignments {
		if a.IsPrimary && a.ID != assignmentID {
			a.RemovePrimary()
			s.assignmentRepo.Update(a)
		}
	}

	assignment.SetAsPrimary()
	return s.assignmentRepo.Update(assignment)
}

// Notes

func (s *leadService) AddNote(inp input.AddNoteInput) (*domain.LeadNote, error) {
	exists, err := s.leadRepo.Exists(inp.LeadID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("lead not found")
	}

	note, err := domain.NewLeadNote(inp.LeadID, inp.Content, inp.CreatedBy)
	if err != nil {
		return nil, err
	}

	if err := s.noteRepo.Save(note); err != nil {
		return nil, err
	}

	return note, nil
}

func (s *leadService) GetNotes(leadID uint) ([]*domain.LeadNote, error) {
	return s.noteRepo.FindByLeadID(leadID)
}

func (s *leadService) UpdateNote(noteID uint, content string) (*domain.LeadNote, error) {
	note, err := s.noteRepo.FindByID(noteID)
	if err != nil {
		return nil, errors.New("note not found")
	}

	if err := note.Update(content); err != nil {
		return nil, err
	}

	if err := s.noteRepo.Update(note); err != nil {
		return nil, err
	}

	return note, nil
}

func (s *leadService) DeleteNote(noteID uint) error {
	_, err := s.noteRepo.FindByID(noteID)
	if err != nil {
		return errors.New("note not found")
	}

	return s.noteRepo.Delete(noteID)
}

// Activities

func (s *leadService) AddActivity(inp input.AddActivityInput) (*domain.LeadActivity, error) {
	exists, err := s.leadRepo.Exists(inp.LeadID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("lead not found")
	}

	activity, err := domain.NewLeadActivity(inp.LeadID, domain.ActivityType(inp.Type), inp.PerformedBy)
	if err != nil {
		return nil, err
	}

	activity.SetDescription(inp.Description)
	activity.SetOutcome(inp.Outcome)

	if inp.PhoneID != nil {
		activity.SetPhone(*inp.PhoneID)
	}

	if inp.Email != "" {
		activity.SetEmail(inp.Email)
	}

	if inp.ScheduledAt != nil {
		scheduledAt, err := time.Parse(time.RFC3339, *inp.ScheduledAt)
		if err != nil {
			return nil, errors.New("invalid scheduled_at format")
		}
		activity.Schedule(scheduledAt)
	}

	if err := s.activityRepo.Save(activity); err != nil {
		return nil, err
	}

	return activity, nil
}

func (s *leadService) GetActivities(leadID uint) ([]*domain.LeadActivity, error) {
	return s.activityRepo.FindByLeadID(leadID)
}

func (s *leadService) CompleteActivity(activityID uint) error {
	activity, err := s.activityRepo.FindByID(activityID)
	if err != nil {
		return errors.New("activity not found")
	}

	activity.Complete()
	return s.activityRepo.Update(activity)
}

func (s *leadService) GetScheduledActivities(entityID uint) ([]*domain.LeadActivity, error) {
	return s.activityRepo.FindScheduledByEntityID(entityID)
}

func (s *leadService) GetOverdueActivities() ([]*domain.LeadActivity, error) {
	return s.activityRepo.FindOverdue()
}