package input

import "torque-dms/core/identity/domain"

type AssignRoleInput struct {
	EntityID uint
	RoleID   uint
}

type AssignResourceInput struct {
	EntityID   uint
	ResourceID uint
	Scope      string
	Reason     string
}

type CheckPermissionInput struct {
	EntityID   uint
	ResourceID uint
	OwnerID    *uint
}

type PermissionService interface {
	// Roles
	CreateRole(name string, description string) (*domain.Role, error)
	GetRoles() ([]*domain.Role, error)
	AssignRole(input AssignRoleInput) error
	RemoveRole(entityID uint, roleID uint) error
	GetEntityRoles(entityID uint) ([]*domain.Role, error)

	// Resources
	CreateResource(code string, name string, urlPattern string, method string, module string) (*domain.Resource, error)
	GetResources() ([]*domain.Resource, error)
	AssignResourceToRole(roleID uint, resourceID uint, scope string) error
	AssignResourceToEntity(input AssignResourceInput) error

	// Check
	CanAccess(input CheckPermissionInput) (bool, error)
	GetScope(entityID uint, resourceID uint) (domain.AccessScope, error)
}