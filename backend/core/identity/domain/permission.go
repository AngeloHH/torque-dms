package domain

import (
	"errors"
	"time"
)

type AccessScope string

const (
	AccessScopeAll  AccessScope = "all"
	AccessScopeOwn  AccessScope = "own"
	AccessScopeTeam AccessScope = "team"
	AccessScopeNone AccessScope = "none"
)

// Resource - representa una URL/acción del sistema
type Resource struct {
	ID             uint
	Code           string
	Name           string
	URLPattern     string
	Method         string
	Module         string
	OwnershipField string
	CreatedAt      time.Time
}

func NewResource(code string, name string, urlPattern string, method string, module string) (*Resource, error) {
	if code == "" {
		return nil, errors.New("code is required")
	}
	if urlPattern == "" {
		return nil, errors.New("url pattern is required")
	}
	if method == "" {
		return nil, errors.New("method is required")
	}

	return &Resource{
		Code:       code,
		Name:       name,
		URLPattern: urlPattern,
		Method:     method,
		Module:     module,
		CreatedAt:  time.Now(),
	}, nil
}

func (r *Resource) SetOwnershipField(field string) {
	r.OwnershipField = field
}

func (r *Resource) RequiresOwnership() bool {
	return r.OwnershipField != ""
}

// Role - agrupa permisos
type Role struct {
	ID           uint
	Name         string
	Description  string
	IsSystemRole bool
	CreatedAt    time.Time
}

func NewRole(name string, description string) (*Role, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	return &Role{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}, nil
}

func (r *Role) SetAsSystemRole() {
	r.IsSystemRole = true
}

// RoleResource - qué scope tiene un rol sobre un recurso
type RoleResource struct {
	ID         uint
	RoleID     uint
	ResourceID uint
	Scope      AccessScope
	CreatedAt  time.Time
}

func NewRoleResource(roleID uint, resourceID uint, scope AccessScope) (*RoleResource, error) {
	if roleID == 0 {
		return nil, errors.New("role is required")
	}
	if resourceID == 0 {
		return nil, errors.New("resource is required")
	}
	if !isValidScope(scope) {
		return nil, errors.New("invalid scope")
	}

	return &RoleResource{
		RoleID:     roleID,
		ResourceID: resourceID,
		Scope:      scope,
		CreatedAt:  time.Now(),
	}, nil
}

// EntityResource - permiso directo de una entity sobre un recurso
type EntityResource struct {
	ID         uint
	EntityID   uint
	ResourceID uint
	Scope      AccessScope
	AssignedBy uint
	Reason     string
	ExpiresAt  *time.Time
	CreatedAt  time.Time
}

func NewEntityResource(entityID uint, resourceID uint, scope AccessScope, assignedBy uint, reason string) (*EntityResource, error) {
	if entityID == 0 {
		return nil, errors.New("entity is required")
	}
	if resourceID == 0 {
		return nil, errors.New("resource is required")
	}
	if !isValidScope(scope) {
		return nil, errors.New("invalid scope")
	}
	if assignedBy == 0 {
		return nil, errors.New("assigned_by is required")
	}

	return &EntityResource{
		EntityID:   entityID,
		ResourceID: resourceID,
		Scope:      scope,
		AssignedBy: assignedBy,
		Reason:     reason,
		CreatedAt:  time.Now(),
	}, nil
}

func (er *EntityResource) SetExpiration(expiresAt time.Time) {
	er.ExpiresAt = &expiresAt
}

func (er *EntityResource) IsExpired() bool {
	if er.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*er.ExpiresAt)
}

// EntityRole - asigna un rol a una entity
type EntityRole struct {
	ID        uint
	EntityID  uint
	RoleID    uint
	CreatedAt time.Time
}

func NewEntityRole(entityID uint, roleID uint) (*EntityRole, error) {
	if entityID == 0 {
		return nil, errors.New("entity is required")
	}
	if roleID == 0 {
		return nil, errors.New("role is required")
	}

	return &EntityRole{
		EntityID:  entityID,
		RoleID:    roleID,
		CreatedAt: time.Now(),
	}, nil
}

// Helper
func isValidScope(scope AccessScope) bool {
	return scope == AccessScopeAll ||
		scope == AccessScopeOwn ||
		scope == AccessScopeTeam ||
		scope == AccessScopeNone
}

// PermissionChecker - lógica para verificar permisos
type PermissionChecker struct {
	entityRoles     []EntityRole
	roleResources   []RoleResource
	entityResources []EntityResource
}

func NewPermissionChecker(entityRoles []EntityRole, roleResources []RoleResource, entityResources []EntityResource) *PermissionChecker {
	return &PermissionChecker{
		entityRoles:     entityRoles,
		roleResources:   roleResources,
		entityResources: entityResources,
	}
}

func (pc *PermissionChecker) GetScope(entityID uint, resourceID uint) AccessScope {
	// 1. Primero revisar permisos directos (entity_resource)
	for _, er := range pc.entityResources {
		if er.EntityID == entityID && er.ResourceID == resourceID {
			if !er.IsExpired() {
				return er.Scope
			}
		}
	}

	// 2. Buscar roles de la entity
	var roleIDs []uint
	for _, er := range pc.entityRoles {
		if er.EntityID == entityID {
			roleIDs = append(roleIDs, er.RoleID)
		}
	}

	// 3. Buscar el scope más permisivo de sus roles
	bestScope := AccessScopeNone
	for _, rr := range pc.roleResources {
		if rr.ResourceID == resourceID && containsUint(roleIDs, rr.RoleID) {
			if scopePriority(rr.Scope) > scopePriority(bestScope) {
				bestScope = rr.Scope
			}
		}
	}

	return bestScope
}

func (pc *PermissionChecker) CanAccess(entityID uint, resourceID uint) bool {
	scope := pc.GetScope(entityID, resourceID)
	return scope != AccessScopeNone
}

func (pc *PermissionChecker) CanAccessOwn(entityID uint, resourceID uint, ownerID uint) bool {
	scope := pc.GetScope(entityID, resourceID)

	if scope == AccessScopeAll {
		return true
	}
	if scope == AccessScopeOwn && entityID == ownerID {
		return true
	}

	return false
}

// Helpers
func containsUint(slice []uint, val uint) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func scopePriority(scope AccessScope) int {
	priorities := map[AccessScope]int{
		AccessScopeNone: 0,
		AccessScopeOwn:  1,
		AccessScopeTeam: 2,
		AccessScopeAll:  3,
	}
	return priorities[scope]
}