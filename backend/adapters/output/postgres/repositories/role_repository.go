package repositories

import (
	"gorm.io/gorm"
	"torque-dms/core/identity/domain"
	"torque-dms/core/identity/ports/output"
	"torque-dms/models"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) output.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Save(role *domain.Role) error {
	model := toRoleModel(role)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	role.ID = model.ID
	return nil
}

func (r *roleRepository) Update(role *domain.Role) error {
	model := toRoleModel(role)
	return r.db.Save(model).Error
}

func (r *roleRepository) FindByID(id uint) (*domain.Role, error) {
	var model models.Role
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainRole(&model), nil
}

func (r *roleRepository) FindByName(name string) (*domain.Role, error) {
	var model models.Role
	result := r.db.Where("name = ?", name).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainRole(&model), nil
}

func (r *roleRepository) FindAll() ([]*domain.Role, error) {
	var modelList []models.Role
	result := r.db.Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	roles := make([]*domain.Role, len(modelList))
	for i, model := range modelList {
		roles[i] = toDomainRole(&model)
	}
	return roles, nil
}

func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

func (r *roleRepository) AssignRoleToEntity(entityRole *domain.EntityRole) error {
	model := &models.EntityRole{
		EntityID:  entityRole.EntityID,
		RoleID:    entityRole.RoleID,
		CreatedAt: entityRole.CreatedAt,
	}
	return r.db.Create(model).Error
}

func (r *roleRepository) RemoveRoleFromEntity(entityID uint, roleID uint) error {
	return r.db.Where("entity_id = ? AND role_id = ?", entityID, roleID).Delete(&models.EntityRole{}).Error
}

func (r *roleRepository) FindRolesByEntityID(entityID uint) ([]*domain.Role, error) {
	var roleList []models.Role
	result := r.db.
		Joins("JOIN entity_role ON entity_role.role_id = role.id").
		Where("entity_role.entity_id = ?", entityID).
		Find(&roleList)
	if result.Error != nil {
		return nil, result.Error
	}

	roles := make([]*domain.Role, len(roleList))
	for i, model := range roleList {
		roles[i] = toDomainRole(&model)
	}
	return roles, nil
}

func (r *roleRepository) FindEntitiesByRoleID(roleID uint) ([]uint, error) {
	var entityIDs []uint
	result := r.db.Model(&models.EntityRole{}).Where("role_id = ?", roleID).Pluck("entity_id", &entityIDs)
	return entityIDs, result.Error
}

// Mappers

func toRoleModel(role *domain.Role) *models.Role {
	return &models.Role{
		ID:           role.ID,
		Name:         role.Name,
		Description:  role.Description,
		IsSystemRole: role.IsSystemRole,
		CreatedAt:    role.CreatedAt,
	}
}

func toDomainRole(m *models.Role) *domain.Role {
	return &domain.Role{
		ID:           m.ID,
		Name:         m.Name,
		Description:  m.Description,
		IsSystemRole: m.IsSystemRole,
		CreatedAt:    m.CreatedAt,
	}
}