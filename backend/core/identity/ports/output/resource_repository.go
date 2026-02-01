package output

import "torque-dms/core/identity/domain"

type ResourceRepository interface {
	// Resource
	Save(resource *domain.Resource) error
	Update(resource *domain.Resource) error
	FindByID(id uint) (*domain.Resource, error)
	FindByCode(code string) (*domain.Resource, error)
	FindByMethodAndPattern(method string, urlPattern string) (*domain.Resource, error)
	FindAll() ([]*domain.Resource, error)
	Delete(id uint) error

	// RoleResource
	AssignResourceToRole(roleResource *domain.RoleResource) error
	RemoveResourceFromRole(roleID uint, resourceID uint) error
	FindResourcesByRoleID(roleID uint) ([]*domain.RoleResource, error)

	// EntityResource
	AssignResourceToEntity(entityResource *domain.EntityResource) error
	RemoveResourceFromEntity(entityID uint, resourceID uint) error
	FindResourcesByEntityID(entityID uint) ([]*domain.EntityResource, error)
}