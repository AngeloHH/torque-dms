package services

import (
	"errors"
	"testing"

	"torque-dms/core/identity/domain"
	"torque-dms/core/identity/ports/input"
)

// --- Mocks ---

type mockRoleRepo struct {
	roles       []*domain.Role
	savedRole   *domain.Role
	saveErr     error
	entityRoles []*domain.Role
	entityIDs   []uint
}

func (m *mockRoleRepo) Save(r *domain.Role) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	r.ID = 1
	m.savedRole = r
	return nil
}
func (m *mockRoleRepo) Update(r *domain.Role) error               { return nil }
func (m *mockRoleRepo) FindByID(id uint) (*domain.Role, error)    { return nil, nil }
func (m *mockRoleRepo) FindByName(name string) (*domain.Role, error) { return nil, nil }
func (m *mockRoleRepo) FindAll() ([]*domain.Role, error)          { return m.roles, nil }
func (m *mockRoleRepo) Delete(id uint) error                      { return nil }
func (m *mockRoleRepo) AssignRoleToEntity(er *domain.EntityRole) error { return nil }
func (m *mockRoleRepo) RemoveRoleFromEntity(entityID uint, roleID uint) error { return nil }
func (m *mockRoleRepo) FindRolesByEntityID(entityID uint) ([]*domain.Role, error) {
	return m.entityRoles, nil
}
func (m *mockRoleRepo) FindEntitiesByRoleID(roleID uint) ([]uint, error) {
	return m.entityIDs, nil
}

type mockResourceRepo struct {
	resources        []*domain.Resource
	savedResource    *domain.Resource
	saveErr          error
	findByMethod     *domain.Resource
	findByMethodErr  error
	roleResources    []*domain.RoleResource
	entityResources  []*domain.EntityResource
}

func (m *mockResourceRepo) Save(r *domain.Resource) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	r.ID = 1
	m.savedResource = r
	return nil
}
func (m *mockResourceRepo) Update(r *domain.Resource) error                 { return nil }
func (m *mockResourceRepo) FindByID(id uint) (*domain.Resource, error)      { return nil, nil }
func (m *mockResourceRepo) FindByCode(code string) (*domain.Resource, error) { return nil, nil }
func (m *mockResourceRepo) FindByMethodAndPattern(method string, pattern string) (*domain.Resource, error) {
	if m.findByMethod != nil {
		return m.findByMethod, nil
	}
	return nil, m.findByMethodErr
}
func (m *mockResourceRepo) FindAll() ([]*domain.Resource, error) { return m.resources, nil }
func (m *mockResourceRepo) Delete(id uint) error                 { return nil }
func (m *mockResourceRepo) AssignResourceToRole(rr *domain.RoleResource) error { return nil }
func (m *mockResourceRepo) RemoveResourceFromRole(roleID uint, resourceID uint) error { return nil }
func (m *mockResourceRepo) FindResourcesByRoleID(roleID uint) ([]*domain.RoleResource, error) {
	return m.roleResources, nil
}
func (m *mockResourceRepo) AssignResourceToEntity(er *domain.EntityResource) error { return nil }
func (m *mockResourceRepo) RemoveResourceFromEntity(entityID uint, resourceID uint) error { return nil }
func (m *mockResourceRepo) FindResourcesByEntityID(entityID uint) ([]*domain.EntityResource, error) {
	return m.entityResources, nil
}

// --- Tests ---

func TestPermissionService_CreateRole(t *testing.T) {
	roleRepo := &mockRoleRepo{}
	resourceRepo := &mockResourceRepo{}
	svc := NewPermissionService(roleRepo, resourceRepo)

	role, err := svc.CreateRole("admin", "Administrator")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role.Name != "admin" {
		t.Errorf("Name = %v, want admin", role.Name)
	}
}

func TestPermissionService_CreateRole_EmptyName(t *testing.T) {
	roleRepo := &mockRoleRepo{}
	resourceRepo := &mockResourceRepo{}
	svc := NewPermissionService(roleRepo, resourceRepo)

	_, err := svc.CreateRole("", "desc")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestPermissionService_CreateResource(t *testing.T) {
	roleRepo := &mockRoleRepo{}
	resourceRepo := &mockResourceRepo{}
	svc := NewPermissionService(roleRepo, resourceRepo)

	resource, err := svc.CreateResource("vehicles.list", "List Vehicles", "/api/vehicles", "GET", "inventory")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resource.Code != "vehicles.list" {
		t.Errorf("Code = %v, want vehicles.list", resource.Code)
	}
}

func TestPermissionService_CheckRouteAccess_NoResource(t *testing.T) {
	roleRepo := &mockRoleRepo{}
	resourceRepo := &mockResourceRepo{
		findByMethodErr: errors.New("not found"),
	}
	svc := NewPermissionService(roleRepo, resourceRepo)

	result, err := svc.CheckRouteAccess(1, "GET", "/api/unknown")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Allowed {
		t.Error("should allow access when resource is not registered")
	}
	if result.Scope != domain.AccessScopeAll {
		t.Errorf("scope = %v, want all", result.Scope)
	}
}

func TestPermissionService_CheckRouteAccess_WithPermission(t *testing.T) {
	resource := &domain.Resource{ID: 10, Code: "vehicles.list", Method: "GET", URLPattern: "/api/vehicles"}

	roleRepo := &mockRoleRepo{
		entityRoles: []*domain.Role{{ID: 100, Name: "admin"}},
	}
	resourceRepo := &mockResourceRepo{
		findByMethod: resource,
		roleResources: []*domain.RoleResource{
			{RoleID: 100, ResourceID: 10, Scope: domain.AccessScopeAll},
		},
		entityResources: []*domain.EntityResource{},
	}
	svc := NewPermissionService(roleRepo, resourceRepo)

	result, err := svc.CheckRouteAccess(1, "GET", "/api/vehicles")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Allowed {
		t.Error("should allow access with role permission")
	}
	if result.Scope != domain.AccessScopeAll {
		t.Errorf("scope = %v, want all", result.Scope)
	}
}

func TestPermissionService_CheckRouteAccess_Denied(t *testing.T) {
	resource := &domain.Resource{ID: 10, Code: "admin.panel", Method: "GET", URLPattern: "/api/admin"}

	roleRepo := &mockRoleRepo{
		entityRoles: []*domain.Role{},
	}
	resourceRepo := &mockResourceRepo{
		findByMethod:    resource,
		roleResources:   []*domain.RoleResource{},
		entityResources: []*domain.EntityResource{},
	}
	svc := NewPermissionService(roleRepo, resourceRepo)

	result, err := svc.CheckRouteAccess(1, "GET", "/api/admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Allowed {
		t.Error("should deny access without any permissions")
	}
}

func TestPermissionService_CheckRouteAccess_OwnershipField(t *testing.T) {
	resource := &domain.Resource{
		ID: 10, Code: "leads.view", Method: "GET",
		URLPattern: "/api/leads/:id", OwnershipField: "created_by",
	}

	roleRepo := &mockRoleRepo{
		entityRoles: []*domain.Role{{ID: 100, Name: "sales"}},
	}
	resourceRepo := &mockResourceRepo{
		findByMethod: resource,
		roleResources: []*domain.RoleResource{
			{RoleID: 100, ResourceID: 10, Scope: domain.AccessScopeOwn},
		},
		entityResources: []*domain.EntityResource{},
	}
	svc := NewPermissionService(roleRepo, resourceRepo)

	result, err := svc.CheckRouteAccess(1, "GET", "/api/leads/:id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Allowed {
		t.Error("should allow access with own scope")
	}
	if result.OwnershipField != "created_by" {
		t.Errorf("OwnershipField = %v, want created_by", result.OwnershipField)
	}
}

func TestPermissionService_CanAccess(t *testing.T) {
	roleRepo := &mockRoleRepo{
		entityRoles: []*domain.Role{{ID: 100, Name: "admin"}},
	}
	resourceRepo := &mockResourceRepo{
		roleResources: []*domain.RoleResource{
			{RoleID: 100, ResourceID: 10, Scope: domain.AccessScopeAll},
		},
		entityResources: []*domain.EntityResource{},
	}
	svc := NewPermissionService(roleRepo, resourceRepo)

	can, err := svc.CanAccess(input.CheckPermissionInput{
		EntityID:   1,
		ResourceID: 10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !can {
		t.Error("should have access via role")
	}
}

func TestPermissionService_CanAccess_WithOwner(t *testing.T) {
	roleRepo := &mockRoleRepo{
		entityRoles: []*domain.Role{{ID: 100, Name: "sales"}},
	}
	resourceRepo := &mockResourceRepo{
		roleResources: []*domain.RoleResource{
			{RoleID: 100, ResourceID: 10, Scope: domain.AccessScopeOwn},
		},
		entityResources: []*domain.EntityResource{},
	}
	svc := NewPermissionService(roleRepo, resourceRepo)

	ownerID := uint(1)
	can, err := svc.CanAccess(input.CheckPermissionInput{
		EntityID:   1,
		ResourceID: 10,
		OwnerID:    &ownerID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !can {
		t.Error("should access own resource")
	}

	otherOwner := uint(999)
	can, err = svc.CanAccess(input.CheckPermissionInput{
		EntityID:   1,
		ResourceID: 10,
		OwnerID:    &otherOwner,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if can {
		t.Error("should not access other's resource with own scope")
	}
}

func TestPermissionService_GetScope(t *testing.T) {
	roleRepo := &mockRoleRepo{
		entityRoles: []*domain.Role{{ID: 100, Name: "viewer"}},
	}
	resourceRepo := &mockResourceRepo{
		roleResources: []*domain.RoleResource{
			{RoleID: 100, ResourceID: 10, Scope: domain.AccessScopeTeam},
		},
		entityResources: []*domain.EntityResource{},
	}
	svc := NewPermissionService(roleRepo, resourceRepo)

	scope, err := svc.GetScope(1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scope != domain.AccessScopeTeam {
		t.Errorf("scope = %v, want team", scope)
	}
}