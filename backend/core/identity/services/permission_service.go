package services

import (
	// "errors"

	"torque-dms/core/identity/domain"
	"torque-dms/core/identity/ports/input"
	"torque-dms/core/identity/ports/output"
)

type permissionService struct {
	roleRepo     output.RoleRepository
	resourceRepo output.ResourceRepository
}

func NewPermissionService(
	roleRepo output.RoleRepository,
	resourceRepo output.ResourceRepository,
) input.PermissionService {
	return &permissionService{
		roleRepo:     roleRepo,
		resourceRepo: resourceRepo,
	}
}

// Roles

func (s *permissionService) CreateRole(name string, description string) (*domain.Role, error) {
	role, err := domain.NewRole(name, description)
	if err != nil {
		return nil, err
	}

	if err := s.roleRepo.Save(role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *permissionService) GetRoles() ([]*domain.Role, error) {
	return s.roleRepo.FindAll()
}

func (s *permissionService) AssignRole(inp input.AssignRoleInput) error {
	entityRole, err := domain.NewEntityRole(inp.EntityID, inp.RoleID)
	if err != nil {
		return err
	}

	return s.roleRepo.AssignRoleToEntity(entityRole)
}

func (s *permissionService) RemoveRole(entityID uint, roleID uint) error {
	return s.roleRepo.RemoveRoleFromEntity(entityID, roleID)
}

func (s *permissionService) GetEntityRoles(entityID uint) ([]*domain.Role, error) {
	return s.roleRepo.FindRolesByEntityID(entityID)
}

// Resources

func (s *permissionService) CreateResource(code string, name string, urlPattern string, method string, module string) (*domain.Resource, error) {
	resource, err := domain.NewResource(code, name, urlPattern, method, module)
	if err != nil {
		return nil, err
	}

	if err := s.resourceRepo.Save(resource); err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *permissionService) GetResources() ([]*domain.Resource, error) {
	return s.resourceRepo.FindAll()
}

func (s *permissionService) AssignResourceToRole(roleID uint, resourceID uint, scope string) error {
	roleResource, err := domain.NewRoleResource(roleID, resourceID, domain.AccessScope(scope))
	if err != nil {
		return err
	}

	return s.resourceRepo.AssignResourceToRole(roleResource)
}

func (s *permissionService) AssignResourceToEntity(inp input.AssignResourceInput) error {
	entityResource, err := domain.NewEntityResource(
		inp.EntityID,
		inp.ResourceID,
		domain.AccessScope(inp.Scope),
		inp.EntityID, // assigned by self for now
		inp.Reason,
	)
	if err != nil {
		return err
	}

	return s.resourceRepo.AssignResourceToEntity(entityResource)
}

// Check

func (s *permissionService) CanAccess(inp input.CheckPermissionInput) (bool, error) {
	// Obtener roles del entity
	roles, err := s.roleRepo.FindRolesByEntityID(inp.EntityID)
	if err != nil {
		return false, err
	}

	// Obtener permisos directos
	entityResources, err := s.resourceRepo.FindResourcesByEntityID(inp.EntityID)
	if err != nil {
		return false, err
	}

	// Obtener permisos por rol
	var roleResources []*domain.RoleResource
	for _, role := range roles {
		rr, err := s.resourceRepo.FindResourcesByRoleID(role.ID)
		if err != nil {
			return false, err
		}
		roleResources = append(roleResources, rr...)
	}

	// Convertir a domain types para el checker
	var entityRoles []domain.EntityRole
	for _, role := range roles {
		entityRoles = append(entityRoles, domain.EntityRole{
			EntityID: inp.EntityID,
			RoleID:   role.ID,
		})
	}
	var domainRoleResources []domain.RoleResource
	for _, rr := range roleResources {
		domainRoleResources = append(domainRoleResources, *rr)
	}

	var domainEntityResources []domain.EntityResource
	for _, er := range entityResources {
		domainEntityResources = append(domainEntityResources, *er)
	}

	// Crear checker y verificar
	checker := domain.NewPermissionChecker(entityRoles, domainRoleResources, domainEntityResources)

	if inp.OwnerID != nil {
		return checker.CanAccessOwn(inp.EntityID, inp.ResourceID, *inp.OwnerID), nil
	}

	return checker.CanAccess(inp.EntityID, inp.ResourceID), nil
}

func (s *permissionService) GetScope(entityID uint, resourceID uint) (domain.AccessScope, error) {
	// Similar a CanAccess pero retorna el scope
	roles, err := s.roleRepo.FindRolesByEntityID(entityID)
	if err != nil {
		return domain.AccessScopeNone, err
	}

	entityResources, err := s.resourceRepo.FindResourcesByEntityID(entityID)
	if err != nil {
		return domain.AccessScopeNone, err
	}

	var roleResources []*domain.RoleResource
	for _, role := range roles {
		rr, err := s.resourceRepo.FindResourcesByRoleID(role.ID)
		if err != nil {
			return domain.AccessScopeNone, err
		}
		roleResources = append(roleResources, rr...)
	}

	var entityRoles []domain.EntityRole
	for _, role := range roles {
		entityRoles = append(entityRoles, domain.EntityRole{
			EntityID: entityID,
			RoleID:   role.ID,
		})
	}

	var domainRoleResources []domain.RoleResource
	for _, rr := range roleResources {
		domainRoleResources = append(domainRoleResources, *rr)
	}

	var domainEntityResources []domain.EntityResource
	for _, er := range entityResources {
		domainEntityResources = append(domainEntityResources, *er)
	}

	checker := domain.NewPermissionChecker(entityRoles, domainRoleResources, domainEntityResources)
	return checker.GetScope(entityID, resourceID), nil
}