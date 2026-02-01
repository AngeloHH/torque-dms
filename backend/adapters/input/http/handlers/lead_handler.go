package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"torque-dms/adapters/input/http/dto/request"
	"torque-dms/adapters/input/http/dto/response"
	"torque-dms/core/sales/domain"
	"torque-dms/core/sales/ports/input"
)

type LeadHandler struct {
	leadService input.LeadService
	stepService input.StepService
}

func NewLeadHandler(leadService input.LeadService, stepService input.StepService) *LeadHandler {
	return &LeadHandler{
		leadService: leadService,
		stepService: stepService,
	}
}

// Lead CRUD

func (h *LeadHandler) Create(c *gin.Context) {
	var req request.CreateLeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lead, err := h.leadService.Create(input.CreateLeadInput{
		EntityID:      req.EntityID,
		VehicleID:     req.VehicleID,
		InterestType:  req.InterestType,
		InterestMake:  req.InterestMake,
		InterestModel: req.InterestModel,
		BudgetMin:     req.BudgetMin,
		BudgetMax:     req.BudgetMax,
		SourceID:      req.SourceID,
		SourceDetail:  req.SourceDetail,
		PresetID:      req.PresetID,
		AssignedTo:    req.AssignedTo,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Inicializar progreso si hay preset
	if req.PresetID != nil {
		h.stepService.InitializeProgress(lead.ID, *req.PresetID)
	}

	c.JSON(http.StatusCreated, toLeadResponse(lead))
}

func (h *LeadHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	lead, err := h.leadService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lead not found"})
		return
	}

	c.JSON(http.StatusOK, toLeadResponse(lead))
}

func (h *LeadHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	leads, err := h.leadService.List(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.LeadResponse, len(leads))
	for i, lead := range leads {
		responseList[i] = *toLeadResponse(lead)
	}

	c.JSON(http.StatusOK, response.LeadListResponse{
		Leads: responseList,
		Total: len(responseList),
	})
}

func (h *LeadHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req request.UpdateLeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lead, err := h.leadService.Update(uint(id), input.UpdateLeadInput{
		VehicleID:     req.VehicleID,
		InterestType:  req.InterestType,
		InterestMake:  req.InterestMake,
		InterestModel: req.InterestModel,
		BudgetMin:     req.BudgetMin,
		BudgetMax:     req.BudgetMax,
		SourceDetail:  req.SourceDetail,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toLeadResponse(lead))
}

func (h *LeadHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.leadService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "lead deleted successfully"})
}

// Sources

func (h *LeadHandler) CreateSource(c *gin.Context) {
	var req request.CreateLeadSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	source, err := h.leadService.CreateSource(input.CreateLeadSourceInput{
		Code:       req.Code,
		Name:       req.Name,
		IsExternal: req.IsExternal,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toLeadSourceResponse(source))
}

func (h *LeadHandler) GetSources(c *gin.Context) {
	sources, err := h.leadService.GetSources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.LeadSourceResponse, len(sources))
	for i, source := range sources {
		responseList[i] = *toLeadSourceResponse(source)
	}

	c.JSON(http.StatusOK, response.LeadSourceListResponse{
		Sources: responseList,
		Total:   len(responseList),
	})
}

func (h *LeadHandler) GetActiveSources(c *gin.Context) {
	sources, err := h.leadService.GetActiveSources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.LeadSourceResponse, len(sources))
	for i, source := range sources {
		responseList[i] = *toLeadSourceResponse(source)
	}

	c.JSON(http.StatusOK, response.LeadSourceListResponse{
		Sources: responseList,
		Total:   len(responseList),
	})
}

func (h *LeadHandler) DeactivateSource(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.leadService.DeactivateSource(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "source deactivated"})
}

func (h *LeadHandler) ActivateSource(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.leadService.ActivateSource(uint(id)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "source activated"})
}

// Assignments

func (h *LeadHandler) Assign(c *gin.Context) {
	leadID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req request.AssignLeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignedBy, _ := c.Get("entity_id")

	assignment, err := h.leadService.Assign(input.AssignLeadInput{
		LeadID:     uint(leadID),
		EntityID:   req.EntityID,
		Role:       req.Role,
		IsPrimary:  req.IsPrimary,
		AssignedBy: assignedBy.(uint),
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toLeadAssignmentResponse(assignment))
}

func (h *LeadHandler) GetAssignments(c *gin.Context) {
	leadID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	assignments, err := h.leadService.GetAssignments(uint(leadID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.LeadAssignmentResponse, len(assignments))
	for i, a := range assignments {
		responseList[i] = *toLeadAssignmentResponse(a)
	}

	c.JSON(http.StatusOK, response.LeadAssignmentListResponse{
		Assignments: responseList,
		Total:       len(responseList),
	})
}

func (h *LeadHandler) RemoveAssignment(c *gin.Context) {
	assignmentID, err := strconv.ParseUint(c.Param("assignmentId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assignment id"})
		return
	}

	if err := h.leadService.RemoveAssignment(uint(assignmentID)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assignment removed"})
}

func (h *LeadHandler) SetPrimaryAssignment(c *gin.Context) {
	leadID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req request.SetPrimaryAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.leadService.SetPrimaryAssignment(uint(leadID), req.AssignmentID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "primary assignment set"})
}

// Notes

func (h *LeadHandler) AddNote(c *gin.Context) {
	leadID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req request.AddNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBy, _ := c.Get("entity_id")

	note, err := h.leadService.AddNote(input.AddNoteInput{
		LeadID:    uint(leadID),
		Content:   req.Content,
		CreatedBy: createdBy.(uint),
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toLeadNoteResponse(note))
}

func (h *LeadHandler) GetNotes(c *gin.Context) {
	leadID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	notes, err := h.leadService.GetNotes(uint(leadID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.LeadNoteResponse, len(notes))
	for i, n := range notes {
		responseList[i] = *toLeadNoteResponse(n)
	}

	c.JSON(http.StatusOK, response.LeadNoteListResponse{
		Notes: responseList,
		Total: len(responseList),
	})
}

func (h *LeadHandler) UpdateNote(c *gin.Context) {
	noteID, err := strconv.ParseUint(c.Param("noteId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	var req request.UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.leadService.UpdateNote(uint(noteID), req.Content)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toLeadNoteResponse(note))
}

func (h *LeadHandler) DeleteNote(c *gin.Context) {
	noteID, err := strconv.ParseUint(c.Param("noteId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	if err := h.leadService.DeleteNote(uint(noteID)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "note deleted"})
}

// Activities

func (h *LeadHandler) AddActivity(c *gin.Context) {
	leadID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req request.AddActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	performedBy, _ := c.Get("entity_id")

	activity, err := h.leadService.AddActivity(input.AddActivityInput{
		LeadID:      uint(leadID),
		Type:        req.Type,
		Description: req.Description,
		Outcome:     req.Outcome,
		PhoneID:     req.PhoneID,
		Email:       req.Email,
		PerformedBy: performedBy.(uint),
		ScheduledAt: req.ScheduledAt,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toLeadActivityResponse(activity))
}

func (h *LeadHandler) GetActivities(c *gin.Context) {
	leadID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	activities, err := h.leadService.GetActivities(uint(leadID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.LeadActivityResponse, len(activities))
	for i, a := range activities {
		responseList[i] = *toLeadActivityResponse(a)
	}

	c.JSON(http.StatusOK, response.LeadActivityListResponse{
		Activities: responseList,
		Total:      len(responseList),
	})
}

func (h *LeadHandler) CompleteActivity(c *gin.Context) {
	activityID, err := strconv.ParseUint(c.Param("activityId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
		return
	}

	if err := h.leadService.CompleteActivity(uint(activityID)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "activity completed"})
}

func (h *LeadHandler) GetMyScheduledActivities(c *gin.Context) {
	entityID, _ := c.Get("entity_id")

	activities, err := h.leadService.GetScheduledActivities(entityID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.LeadActivityResponse, len(activities))
	for i, a := range activities {
		responseList[i] = *toLeadActivityResponse(a)
	}

	c.JSON(http.StatusOK, response.LeadActivityListResponse{
		Activities: responseList,
		Total:      len(responseList),
	})
}

func (h *LeadHandler) GetOverdueActivities(c *gin.Context) {
	activities, err := h.leadService.GetOverdueActivities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.LeadActivityResponse, len(activities))
	for i, a := range activities {
		responseList[i] = *toLeadActivityResponse(a)
	}

	c.JSON(http.StatusOK, response.LeadActivityListResponse{
		Activities: responseList,
		Total:      len(responseList),
	})
}

// Progress

func (h *LeadHandler) GetProgress(c *gin.Context) {
	leadID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	progress, err := h.stepService.GetProgress(uint(leadID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseList := make([]response.LeadStepProgressResponse, len(progress))
	for i, p := range progress {
		responseList[i] = *toLeadStepProgressResponse(p)
	}

	c.JSON(http.StatusOK, response.LeadStepProgressListResponse{
		Progress: responseList,
		Total:    len(responseList),
	})
}

func (h *LeadHandler) UpdateProgress(c *gin.Context) {
	leadID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	stepID, err := strconv.ParseUint(c.Param("stepId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step id"})
		return
	}

	var req request.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	completedBy, _ := c.Get("entity_id")

	progress, err := h.stepService.UpdateProgress(input.UpdateProgressInput{
		LeadID:      uint(leadID),
		StepID:      uint(stepID),
		Status:      req.Status,
		CompletedBy: completedBy.(uint),
		Notes:       req.Notes,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toLeadStepProgressResponse(progress))
}

// Helpers

func toLeadResponse(l *domain.Lead) *response.LeadResponse {
	return &response.LeadResponse{
		ID:            l.ID,
		EntityID:      l.EntityID,
		VehicleID:     l.VehicleID,
		InterestType:  l.InterestType,
		InterestMake:  l.InterestMake,
		InterestModel: l.InterestModel,
		BudgetMin:     l.BudgetMin,
		BudgetMax:     l.BudgetMax,
		SourceID:      l.SourceID,
		SourceDetail:  l.SourceDetail,
		PresetID:      l.PresetID,
		CreatedAt:     l.CreatedAt,
		ModifiedAt:    l.ModifiedAt,
	}
}

func toLeadSourceResponse(s *domain.LeadSource) *response.LeadSourceResponse {
	return &response.LeadSourceResponse{
		ID:         s.ID,
		Code:       s.Code,
		Name:       s.Name,
		IsExternal: s.IsExternal,
		Active:     s.Active,
		CreatedAt:  s.CreatedAt,
	}
}

func toLeadAssignmentResponse(a *domain.LeadAssignment) *response.LeadAssignmentResponse {
	return &response.LeadAssignmentResponse{
		ID:         a.ID,
		LeadID:     a.LeadID,
		EntityID:   a.EntityID,
		Role:       string(a.Role),
		IsPrimary:  a.IsPrimary,
		AssignedBy: a.AssignedBy,
		Active:     a.Active,
		CreatedAt:  a.CreatedAt,
	}
}

func toLeadNoteResponse(n *domain.LeadNote) *response.LeadNoteResponse {
	return &response.LeadNoteResponse{
		ID:         n.ID,
		LeadID:     n.LeadID,
		Content:    n.Content,
		CreatedBy:  n.CreatedBy,
		CreatedAt:  n.CreatedAt,
		ModifiedAt: n.ModifiedAt,
	}
}

func toLeadActivityResponse(a *domain.LeadActivity) *response.LeadActivityResponse {
	return &response.LeadActivityResponse{
		ID:          a.ID,
		LeadID:      a.LeadID,
		Type:        string(a.Type),
		Description: a.Description,
		Outcome:     a.Outcome,
		PhoneID:     a.PhoneID,
		Email:       a.Email,
		PerformedBy: a.PerformedBy,
		ScheduledAt: a.ScheduledAt,
		CompletedAt: a.CompletedAt,
		IsOverdue:   a.IsOverdue(),
		CreatedAt:   a.CreatedAt,
	}
}

func toLeadStepProgressResponse(p *domain.LeadStepProgress) *response.LeadStepProgressResponse {
	return &response.LeadStepProgressResponse{
		ID:          p.ID,
		LeadID:      p.LeadID,
		StepID:      p.StepID,
		Status:      string(p.Status),
		StartedAt:   p.StartedAt,
		CompletedAt: p.CompletedAt,
		CompletedBy: p.CompletedBy,
		Notes:       p.Notes,
		CreatedAt:   p.CreatedAt,
	}
}