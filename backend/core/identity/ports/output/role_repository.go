package output

import "torque-dms/core/identity/domain"

type RoleRepository interface {
	// Role
	Save(role *domain.Role) error
	Update(role *domain.Role) error
	FindByID(id uint) (*domain.Role, error)
	FindByName(name string) (*domain.Role, error)
	FindAll() ([]*domain.Role, error)
	Delete(id uint) error

	// EntityRole
	AssignRoleToEntity(entityRole *domain.EntityRole) error
	RemoveRoleFromEntity(entityID uint, roleID uint) error
	FindRolesByEntityID(entityID uint) ([]*domain.Role, error)
	FindEntitiesByRoleID(roleID uint) ([]uint, error)
}