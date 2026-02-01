package repositories

import (
	"gorm.io/gorm"
	"torque-dms/core/identity/domain"
	"torque-dms/core/identity/ports/output"
	"torque-dms/models"
)

type resourceRepository struct {
	db *gorm.DB
}

func NewResourceRepository(db *gorm.DB) output.ResourceRepository {
	return &resourceRepository{db: db}
}

// Resource

func (r *resourceRepository) Save(resource *domain.Resource) error {
	model := toResourceModel(resource)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	resource.ID = model.ID
	return nil
}

func (r *resourceRepository) Update(resource *domain.Resource) error {
	model := toResourceModel(resource)
	return r.db.Save(model).Error
}

func (r *resourceRepository) FindByID(id uint) (*domain.Resource, error) {
	var model models.Resource
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainResource(&model), nil
}

func (r *resourceRepository) FindByCode(code string) (*domain.Resource, error) {
	var model models.Resource
	result := r.db.Where("code = ?", code).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainResource(&model), nil
}

func (r *resourceRepository) FindByMethodAndPattern(method string, urlPattern string) (*domain.Resource, error) {
	var model models.Resource
	result := r.db.Where("method = ? AND url_pattern = ?", method, urlPattern).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainResource(&model), nil
}

func (r *resourceRepository) FindAll() ([]*domain.Resource, error) {
	var modelList []models.Resource
	result := r.db.Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	resources := make([]*domain.Resource, len(modelList))
	for i, model := range modelList {
		resources[i] = toDomainResource(&model)
	}
	return resources, nil
}

func (r *resourceRepository) Delete(id uint) error {
	return r.db.Delete(&models.Resource{}, id).Error
}

// RoleResource

func (r *resourceRepository) AssignResourceToRole(roleResource *domain.RoleResource) error {
	model := &models.RoleResource{
		RoleID:     roleResource.RoleID,
		ResourceID: roleResource.ResourceID,
		Scope:      models.AccessScope(roleResource.Scope),
		CreatedAt:  roleResource.CreatedAt,
	}
	return r.db.Create(model).Error
}

func (r *resourceRepository) RemoveResourceFromRole(roleID uint, resourceID uint) error {
	return r.db.Where("role_id = ? AND resource_id = ?", roleID, resourceID).Delete(&models.RoleResource{}).Error
}

func (r *resourceRepository) FindResourcesByRoleID(roleID uint) ([]*domain.RoleResource, error) {
	var modelList []models.RoleResource
	result := r.db.Where("role_id = ?", roleID).Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	roleResources := make([]*domain.RoleResource, len(modelList))
	for i, model := range modelList {
		roleResources[i] = toDomainRoleResource(&model)
	}
	return roleResources, nil
}

// EntityResource

func (r *resourceRepository) AssignResourceToEntity(entityResource *domain.EntityResource) error {
	model := &models.EntityResource{
		EntityID:   entityResource.EntityID,
		ResourceID: entityResource.ResourceID,
		Scope:      models.AccessScope(entityResource.Scope),
		AssignedBy: entityResource.AssignedBy,
		Reason:     entityResource.Reason,
		ExpiresAt:  entityResource.ExpiresAt,
		CreatedAt:  entityResource.CreatedAt,
	}
	return r.db.Create(model).Error
}

func (r *resourceRepository) RemoveResourceFromEntity(entityID uint, resourceID uint) error {
	return r.db.Where("entity_id = ? AND resource_id = ?", entityID, resourceID).Delete(&models.EntityResource{}).Error
}

func (r *resourceRepository) FindResourcesByEntityID(entityID uint) ([]*domain.EntityResource, error) {
	var modelList []models.EntityResource
	result := r.db.Where("entity_id = ?", entityID).Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	entityResources := make([]*domain.EntityResource, len(modelList))
	for i, model := range modelList {
		entityResources[i] = toDomainEntityResource(&model)
	}
	return entityResources, nil
}

// Mappers

func toResourceModel(r *domain.Resource) *models.Resource {
	return &models.Resource{
		ID:             r.ID,
		Code:           r.Code,
		Name:           r.Name,
		URLPattern:     r.URLPattern,
		Method:         r.Method,
		Module:         r.Module,
		OwnershipField: r.OwnershipField,
		CreatedAt:      r.CreatedAt,
	}
}

func toDomainResource(m *models.Resource) *domain.Resource {
	return &domain.Resource{
		ID:             m.ID,
		Code:           m.Code,
		Name:           m.Name,
		URLPattern:     m.URLPattern,
		Method:         m.Method,
		Module:         m.Module,
		OwnershipField: m.OwnershipField,
		CreatedAt:      m.CreatedAt,
	}
}

func toDomainRoleResource(m *models.RoleResource) *domain.RoleResource {
	return &domain.RoleResource{
		ID:         m.ID,
		RoleID:     m.RoleID,
		ResourceID: m.ResourceID,
		Scope:      domain.AccessScope(m.Scope),
		CreatedAt:  m.CreatedAt,
	}
}

func toDomainEntityResource(m *models.EntityResource) *domain.EntityResource {
	return &domain.EntityResource{
		ID:         m.ID,
		EntityID:   m.EntityID,
		ResourceID: m.ResourceID,
		Scope:      domain.AccessScope(m.Scope),
		AssignedBy: m.AssignedBy,
		Reason:     m.Reason,
		ExpiresAt:  m.ExpiresAt,
		CreatedAt:  m.CreatedAt,
	}
}