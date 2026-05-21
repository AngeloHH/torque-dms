package domain

import (
	"testing"
	"time"
)

// --- Resource ---

func TestNewResource_Valid(t *testing.T) {
	r, err := NewResource("vehicles.list", "List Vehicles", "/api/vehicles", "GET", "inventory")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Code != "vehicles.list" {
		t.Errorf("Code = %v, want vehicles.list", r.Code)
	}
	if r.RequiresOwnership() {
		t.Error("new resource should not require ownership")
	}
}

func TestNewResource_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		url     string
		method  string
	}{
		{"empty code", "", "/api/test", "GET"},
		{"empty url", "test", "", "GET"},
		{"empty method", "test", "/api/test", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewResource(tt.code, "Test", tt.url, tt.method, "mod")
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestResource_Ownership(t *testing.T) {
	r, _ := NewResource("leads.view", "View Lead", "/api/leads/:id", "GET", "sales")
	r.SetOwnershipField("created_by")

	if !r.RequiresOwnership() {
		t.Error("should require ownership after SetOwnershipField")
	}
	if r.OwnershipField != "created_by" {
		t.Errorf("OwnershipField = %v, want created_by", r.OwnershipField)
	}
}

// --- Role ---

func TestNewRole_Valid(t *testing.T) {
	role, err := NewRole("admin", "Administrator role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role.Name != "admin" {
		t.Errorf("Name = %v, want admin", role.Name)
	}
	if role.IsSystemRole {
		t.Error("new role should not be system role")
	}
}

func TestNewRole_EmptyName(t *testing.T) {
	_, err := NewRole("", "some description")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestRole_SetAsSystemRole(t *testing.T) {
	role, _ := NewRole("superadmin", "")
	role.SetAsSystemRole()
	if !role.IsSystemRole {
		t.Error("should be system role after SetAsSystemRole")
	}
}

// --- RoleResource ---

func TestNewRoleResource_Valid(t *testing.T) {
	rr, err := NewRoleResource(1, 2, AccessScopeAll)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rr.RoleID != 1 || rr.ResourceID != 2 {
		t.Errorf("RoleID=%v ResourceID=%v, want 1,2", rr.RoleID, rr.ResourceID)
	}
}

func TestNewRoleResource_Invalid(t *testing.T) {
	tests := []struct {
		name       string
		roleID     uint
		resourceID uint
		scope      AccessScope
	}{
		{"zero role", 0, 1, AccessScopeAll},
		{"zero resource", 1, 0, AccessScopeAll},
		{"invalid scope", 1, 1, "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRoleResource(tt.roleID, tt.resourceID, tt.scope)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// --- EntityResource ---

func TestNewEntityResource_Valid(t *testing.T) {
	er, err := NewEntityResource(1, 2, AccessScopeOwn, 3, "testing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if er.EntityID != 1 || er.ResourceID != 2 || er.AssignedBy != 3 {
		t.Error("fields not set correctly")
	}
	if er.IsExpired() {
		t.Error("should not be expired without expiration set")
	}
}

func TestNewEntityResource_Invalid(t *testing.T) {
	tests := []struct {
		name       string
		entityID   uint
		resourceID uint
		scope      AccessScope
		assignedBy uint
	}{
		{"zero entity", 0, 1, AccessScopeAll, 1},
		{"zero resource", 1, 0, AccessScopeAll, 1},
		{"invalid scope", 1, 1, "bad", 1},
		{"zero assigned_by", 1, 1, AccessScopeAll, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEntityResource(tt.entityID, tt.resourceID, tt.scope, tt.assignedBy, "reason")
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestEntityResource_Expiration(t *testing.T) {
	er, _ := NewEntityResource(1, 1, AccessScopeAll, 1, "test")

	// Not expired without date
	if er.IsExpired() {
		t.Error("should not be expired without expiration")
	}

	// Set future expiration
	future := time.Now().Add(1 * time.Hour)
	er.SetExpiration(future)
	if er.IsExpired() {
		t.Error("should not be expired with future date")
	}

	// Set past expiration
	past := time.Now().Add(-1 * time.Hour)
	er.SetExpiration(past)
	if !er.IsExpired() {
		t.Error("should be expired with past date")
	}
}

// --- EntityRole ---

func TestNewEntityRole_Valid(t *testing.T) {
	er, err := NewEntityRole(1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if er.EntityID != 1 || er.RoleID != 2 {
		t.Error("fields not set correctly")
	}
}

func TestNewEntityRole_Invalid(t *testing.T) {
	if _, err := NewEntityRole(0, 1); err == nil {
		t.Error("expected error for zero entity")
	}
	if _, err := NewEntityRole(1, 0); err == nil {
		t.Error("expected error for zero role")
	}
}

// --- PermissionChecker ---

func TestPermissionChecker_DirectPermission(t *testing.T) {
	entityResources := []EntityResource{
		{EntityID: 1, ResourceID: 10, Scope: AccessScopeAll},
	}

	checker := NewPermissionChecker(nil, nil, entityResources)

	if scope := checker.GetScope(1, 10); scope != AccessScopeAll {
		t.Errorf("scope = %v, want all", scope)
	}
	if !checker.CanAccess(1, 10) {
		t.Error("should have access via direct permission")
	}
}

func TestPermissionChecker_RolePermission(t *testing.T) {
	entityRoles := []EntityRole{
		{EntityID: 1, RoleID: 100},
	}
	roleResources := []RoleResource{
		{RoleID: 100, ResourceID: 10, Scope: AccessScopeOwn},
	}

	checker := NewPermissionChecker(entityRoles, roleResources, nil)

	if scope := checker.GetScope(1, 10); scope != AccessScopeOwn {
		t.Errorf("scope = %v, want own", scope)
	}
	if !checker.CanAccess(1, 10) {
		t.Error("should have access via role")
	}
}

func TestPermissionChecker_NoPermission(t *testing.T) {
	checker := NewPermissionChecker(nil, nil, nil)

	if scope := checker.GetScope(1, 10); scope != AccessScopeNone {
		t.Errorf("scope = %v, want none", scope)
	}
	if checker.CanAccess(1, 10) {
		t.Error("should not have access without any permissions")
	}
}

func TestPermissionChecker_DirectOverridesRole(t *testing.T) {
	entityRoles := []EntityRole{
		{EntityID: 1, RoleID: 100},
	}
	roleResources := []RoleResource{
		{RoleID: 100, ResourceID: 10, Scope: AccessScopeAll},
	}
	entityResources := []EntityResource{
		{EntityID: 1, ResourceID: 10, Scope: AccessScopeOwn},
	}

	checker := NewPermissionChecker(entityRoles, roleResources, entityResources)

	// Direct permission takes precedence
	if scope := checker.GetScope(1, 10); scope != AccessScopeOwn {
		t.Errorf("scope = %v, want own (direct overrides role)", scope)
	}
}

func TestPermissionChecker_BestRoleScope(t *testing.T) {
	entityRoles := []EntityRole{
		{EntityID: 1, RoleID: 100},
		{EntityID: 1, RoleID: 200},
	}
	roleResources := []RoleResource{
		{RoleID: 100, ResourceID: 10, Scope: AccessScopeOwn},
		{RoleID: 200, ResourceID: 10, Scope: AccessScopeAll},
	}

	checker := NewPermissionChecker(entityRoles, roleResources, nil)

	// Should pick the most permissive scope
	if scope := checker.GetScope(1, 10); scope != AccessScopeAll {
		t.Errorf("scope = %v, want all (best of multiple roles)", scope)
	}
}

func TestPermissionChecker_ExpiredDirectPermission(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	entityResources := []EntityResource{
		{EntityID: 1, ResourceID: 10, Scope: AccessScopeAll, ExpiresAt: &past},
	}

	checker := NewPermissionChecker(nil, nil, entityResources)

	if scope := checker.GetScope(1, 10); scope != AccessScopeNone {
		t.Errorf("scope = %v, want none (expired)", scope)
	}
}

func TestPermissionChecker_CanAccessOwn(t *testing.T) {
	entityResources := []EntityResource{
		{EntityID: 1, ResourceID: 10, Scope: AccessScopeOwn},
	}
	checker := NewPermissionChecker(nil, nil, entityResources)

	// Owner matches
	if !checker.CanAccessOwn(1, 10, 1) {
		t.Error("should access own resource when entity is owner")
	}

	// Owner doesn't match
	if checker.CanAccessOwn(1, 10, 2) {
		t.Error("should not access when entity is not owner with scope=own")
	}
}

func TestPermissionChecker_CanAccessOwn_AllScope(t *testing.T) {
	entityResources := []EntityResource{
		{EntityID: 1, ResourceID: 10, Scope: AccessScopeAll},
	}
	checker := NewPermissionChecker(nil, nil, entityResources)

	// scope=all can access anyone's resources
	if !checker.CanAccessOwn(1, 10, 999) {
		t.Error("scope=all should access any owner's resource")
	}
}