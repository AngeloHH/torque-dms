package services

import (
	"errors"
	"testing"

	"torque-dms/core/sales/domain"
	"torque-dms/core/sales/ports/input"
)

// --- Mocks ---

type mockLeadRepo struct {
	saved    *domain.Lead
	findByID *domain.Lead
	exists   bool
	leads    []*domain.Lead
}

func (m *mockLeadRepo) Save(l *domain.Lead) error {
	l.ID = 1
	m.saved = l
	return nil
}
func (m *mockLeadRepo) Update(l *domain.Lead) error { return nil }
func (m *mockLeadRepo) FindByID(id uint) (*domain.Lead, error) {
	if m.findByID != nil {
		return m.findByID, nil
	}
	return nil, errors.New("not found")
}
func (m *mockLeadRepo) FindByEntityID(id uint) ([]*domain.Lead, error) { return m.leads, nil }
func (m *mockLeadRepo) FindAll(l int, o int) ([]*domain.Lead, error)   { return m.leads, nil }
func (m *mockLeadRepo) Delete(id uint) error                           { return nil }
func (m *mockLeadRepo) Exists(id uint) (bool, error)                   { return m.exists, nil }

type mockSourceRepo struct {
	saved    *domain.LeadSource
	findByID *domain.LeadSource
	exists   bool
	sources  []*domain.LeadSource
}

func (m *mockSourceRepo) Save(s *domain.LeadSource) error {
	s.ID = 1
	m.saved = s
	return nil
}
func (m *mockSourceRepo) Update(s *domain.LeadSource) error { return nil }
func (m *mockSourceRepo) FindByID(id uint) (*domain.LeadSource, error) {
	if m.findByID != nil {
		return m.findByID, nil
	}
	return nil, errors.New("not found")
}
func (m *mockSourceRepo) FindByCode(code string) (*domain.LeadSource, error) { return nil, nil }
func (m *mockSourceRepo) FindAll() ([]*domain.LeadSource, error)             { return m.sources, nil }
func (m *mockSourceRepo) FindActive() ([]*domain.LeadSource, error)          { return nil, nil }
func (m *mockSourceRepo) Delete(id uint) error                               { return nil }
func (m *mockSourceRepo) Exists(id uint) (bool, error)                       { return m.exists, nil }

type mockAssignmentRepo struct {
	saved       *domain.LeadAssignment
	findByID    *domain.LeadAssignment
	assignments []*domain.LeadAssignment
}

func (m *mockAssignmentRepo) Save(a *domain.LeadAssignment) error {
	a.ID = 1
	m.saved = a
	return nil
}
func (m *mockAssignmentRepo) Update(a *domain.LeadAssignment) error { return nil }
func (m *mockAssignmentRepo) FindByID(id uint) (*domain.LeadAssignment, error) {
	if m.findByID != nil {
		return m.findByID, nil
	}
	return nil, errors.New("not found")
}
func (m *mockAssignmentRepo) FindByLeadID(id uint) ([]*domain.LeadAssignment, error) {
	return m.assignments, nil
}
func (m *mockAssignmentRepo) FindPrimaryByLeadID(id uint) (*domain.LeadAssignment, error) {
	return nil, nil
}
func (m *mockAssignmentRepo) FindByEntityID(id uint) ([]*domain.LeadAssignment, error) {
	return nil, nil
}
func (m *mockAssignmentRepo) Delete(id uint) error { return nil }

type mockNoteRepo struct {
	saved    *domain.LeadNote
	findByID *domain.LeadNote
	notes    []*domain.LeadNote
}

func (m *mockNoteRepo) Save(n *domain.LeadNote) error {
	n.ID = 1
	m.saved = n
	return nil
}
func (m *mockNoteRepo) Update(n *domain.LeadNote) error { return nil }
func (m *mockNoteRepo) FindByID(id uint) (*domain.LeadNote, error) {
	if m.findByID != nil {
		return m.findByID, nil
	}
	return nil, errors.New("not found")
}
func (m *mockNoteRepo) FindByLeadID(id uint) ([]*domain.LeadNote, error) { return m.notes, nil }
func (m *mockNoteRepo) Delete(id uint) error                             { return nil }

type mockActivityRepo struct {
	saved    *domain.LeadActivity
	findByID *domain.LeadActivity
}

func (m *mockActivityRepo) Save(a *domain.LeadActivity) error {
	a.ID = 1
	m.saved = a
	return nil
}
func (m *mockActivityRepo) Update(a *domain.LeadActivity) error { return nil }
func (m *mockActivityRepo) FindByID(id uint) (*domain.LeadActivity, error) {
	if m.findByID != nil {
		return m.findByID, nil
	}
	return nil, errors.New("not found")
}
func (m *mockActivityRepo) FindByLeadID(id uint) ([]*domain.LeadActivity, error) { return nil, nil }
func (m *mockActivityRepo) FindScheduledByEntityID(id uint) ([]*domain.LeadActivity, error) {
	return nil, nil
}
func (m *mockActivityRepo) FindOverdue() ([]*domain.LeadActivity, error) { return nil, nil }
func (m *mockActivityRepo) Delete(id uint) error                        { return nil }

func newLeadService() (input.LeadService, *mockLeadRepo, *mockSourceRepo, *mockAssignmentRepo, *mockNoteRepo, *mockActivityRepo) {
	lr := &mockLeadRepo{exists: true}
	sr := &mockSourceRepo{exists: true}
	ar := &mockAssignmentRepo{}
	nr := &mockNoteRepo{}
	acr := &mockActivityRepo{}
	return NewLeadService(lr, sr, ar, nr, acr), lr, sr, ar, nr, acr
}

// --- Lead CRUD ---

func TestLeadService_Create_Success(t *testing.T) {
	svc, lr, _, _, _, _ := newLeadService()

	lead, err := svc.Create(input.CreateLeadInput{
		EntityID:     1,
		SourceID:     1,
		InterestType: "purchase",
		InterestMake: "Toyota",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lead.EntityID != 1 {
		t.Errorf("EntityID = %v, want 1", lead.EntityID)
	}
	if lr.saved == nil {
		t.Error("lead should be saved")
	}
}

func TestLeadService_Create_SourceNotFound(t *testing.T) {
	svc, _, sr, _, _, _ := newLeadService()
	sr.exists = false

	_, err := svc.Create(input.CreateLeadInput{EntityID: 1, SourceID: 999})
	if err == nil {
		t.Error("expected error for non-existent source")
	}
}

func TestLeadService_Create_WithAssignment(t *testing.T) {
	svc, _, _, ar, _, _ := newLeadService()

	_, err := svc.Create(input.CreateLeadInput{
		EntityID:   1,
		SourceID:   1,
		AssignedTo: 5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ar.saved == nil {
		t.Error("assignment should be saved")
	}
	if !ar.saved.IsPrimary {
		t.Error("initial assignment should be primary")
	}
}

func TestLeadService_Create_WithBudget(t *testing.T) {
	svc, _, _, _, _, _ := newLeadService()

	lead, err := svc.Create(input.CreateLeadInput{
		EntityID:  1,
		SourceID:  1,
		BudgetMin: 20000,
		BudgetMax: 40000,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lead.BudgetMin != 20000 || lead.BudgetMax != 40000 {
		t.Errorf("budget = %v-%v, want 20000-40000", lead.BudgetMin, lead.BudgetMax)
	}
}

func TestLeadService_Update_Success(t *testing.T) {
	lead, _ := domain.NewLead(1, 1)
	lead.ID = 1
	svc, lr, _, _, _, _ := newLeadService()
	lr.findByID = lead

	interestType := "lease"
	updated, err := svc.Update(1, input.UpdateLeadInput{InterestType: &interestType})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.InterestType != "lease" {
		t.Errorf("InterestType = %v, want lease", updated.InterestType)
	}
}

func TestLeadService_Update_NotFound(t *testing.T) {
	svc, _, _, _, _, _ := newLeadService()
	interestType := "lease"
	_, err := svc.Update(999, input.UpdateLeadInput{InterestType: &interestType})
	if err == nil {
		t.Error("expected error for non-existent lead")
	}
}

func TestLeadService_Delete_Success(t *testing.T) {
	svc, _, _, _, _, _ := newLeadService()
	if err := svc.Delete(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLeadService_Delete_NotFound(t *testing.T) {
	svc, lr, _, _, _, _ := newLeadService()
	lr.exists = false
	if err := svc.Delete(999); err == nil {
		t.Error("expected error for non-existent lead")
	}
}

// --- Sources ---

func TestLeadService_CreateSource(t *testing.T) {
	svc, _, _, _, _, _ := newLeadService()
	src, err := svc.CreateSource(input.CreateLeadSourceInput{
		Code: "web", Name: "Website", IsExternal: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Code != "web" {
		t.Errorf("Code = %v, want web", src.Code)
	}
}

func TestLeadService_DeactivateSource(t *testing.T) {
	src, _ := domain.NewLeadSource("web", "Website", true)
	svc, _, sr, _, _, _ := newLeadService()
	sr.findByID = src

	if err := svc.DeactivateSource(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLeadService_DeactivateSource_NotFound(t *testing.T) {
	svc, _, _, _, _, _ := newLeadService()
	if err := svc.DeactivateSource(999); err == nil {
		t.Error("expected error for non-existent source")
	}
}

// --- Assignments ---

func TestLeadService_Assign_Success(t *testing.T) {
	svc, _, _, _, _, _ := newLeadService()

	a, err := svc.Assign(input.AssignLeadInput{
		LeadID: 1, EntityID: 5, Role: "salesperson", AssignedBy: 3,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.EntityID != 5 {
		t.Errorf("EntityID = %v, want 5", a.EntityID)
	}
}

func TestLeadService_Assign_LeadNotFound(t *testing.T) {
	svc, lr, _, _, _, _ := newLeadService()
	lr.exists = false

	_, err := svc.Assign(input.AssignLeadInput{
		LeadID: 999, EntityID: 1, Role: "salesperson", AssignedBy: 1,
	})
	if err == nil {
		t.Error("expected error for non-existent lead")
	}
}

func TestLeadService_RemoveAssignment(t *testing.T) {
	a, _ := domain.NewLeadAssignment(1, 1, domain.AssignmentRoleSalesperson, 1)
	a.ID = 10
	svc, _, _, ar, _, _ := newLeadService()
	ar.findByID = a

	if err := svc.RemoveAssignment(10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLeadService_SetPrimaryAssignment_WrongLead(t *testing.T) {
	a, _ := domain.NewLeadAssignment(2, 1, domain.AssignmentRoleSalesperson, 1)
	a.ID = 10
	svc, _, _, ar, _, _ := newLeadService()
	ar.findByID = a

	err := svc.SetPrimaryAssignment(1, 10)
	if err == nil {
		t.Error("expected error when assignment belongs to different lead")
	}
}

// --- Notes ---

func TestLeadService_AddNote(t *testing.T) {
	svc, _, _, _, _, _ := newLeadService()

	note, err := svc.AddNote(input.AddNoteInput{
		LeadID: 1, Content: "Interested in financing", CreatedBy: 5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.Content != "Interested in financing" {
		t.Errorf("Content = %v", note.Content)
	}
}

func TestLeadService_AddNote_LeadNotFound(t *testing.T) {
	svc, lr, _, _, _, _ := newLeadService()
	lr.exists = false

	_, err := svc.AddNote(input.AddNoteInput{LeadID: 999, Content: "test", CreatedBy: 1})
	if err == nil {
		t.Error("expected error for non-existent lead")
	}
}

func TestLeadService_UpdateNote(t *testing.T) {
	note, _ := domain.NewLeadNote(1, "original", 1)
	note.ID = 10
	svc, _, _, _, nr, _ := newLeadService()
	nr.findByID = note

	updated, err := svc.UpdateNote(10, "new content")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Content != "new content" {
		t.Errorf("Content = %v, want new content", updated.Content)
	}
}

func TestLeadService_DeleteNote_NotFound(t *testing.T) {
	svc, _, _, _, _, _ := newLeadService()
	if err := svc.DeleteNote(999); err == nil {
		t.Error("expected error for non-existent note")
	}
}

// --- Activities ---

func TestLeadService_AddActivity(t *testing.T) {
	svc, _, _, _, _, _ := newLeadService()

	act, err := svc.AddActivity(input.AddActivityInput{
		LeadID:      1,
		Type:        "call_outbound",
		Description: "Follow up",
		PerformedBy: 5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if act.Type != domain.ActivityTypeCallOutbound {
		t.Errorf("Type = %v, want call_outbound", act.Type)
	}
}

func TestLeadService_AddActivity_WithSchedule(t *testing.T) {
	svc, _, _, _, _, _ := newLeadService()

	scheduled := "2026-06-01T10:00:00Z"
	act, err := svc.AddActivity(input.AddActivityInput{
		LeadID:      1,
		Type:        "appointment_scheduled",
		PerformedBy: 5,
		ScheduledAt: &scheduled,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !act.IsScheduled() {
		t.Error("activity should be scheduled")
	}
}

func TestLeadService_AddActivity_InvalidSchedule(t *testing.T) {
	svc, _, _, _, _, _ := newLeadService()

	bad := "not-a-date"
	_, err := svc.AddActivity(input.AddActivityInput{
		LeadID: 1, Type: "call_outbound", PerformedBy: 1, ScheduledAt: &bad,
	})
	if err == nil {
		t.Error("expected error for invalid date format")
	}
}

func TestLeadService_CompleteActivity(t *testing.T) {
	act, _ := domain.NewLeadActivity(1, domain.ActivityTypeCallOutbound, 1)
	act.ID = 10
	svc, _, _, _, _, acr := newLeadService()
	acr.findByID = act

	if err := svc.CompleteActivity(10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLeadService_CompleteActivity_NotFound(t *testing.T) {
	svc, _, _, _, _, _ := newLeadService()
	if err := svc.CompleteActivity(999); err == nil {
		t.Error("expected error for non-existent activity")
	}
}