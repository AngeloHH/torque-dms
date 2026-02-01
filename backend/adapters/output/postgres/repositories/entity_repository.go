package repositories

import (
	"gorm.io/gorm"
	"torque-dms/core/identity/domain"
	"torque-dms/core/identity/ports/output"
	"torque-dms/models"
)

type entityRepository struct {
	db *gorm.DB
}

func NewEntityRepository(db *gorm.DB) output.EntityRepository {
	return &entityRepository{db: db}
}

func (r *entityRepository) Save(entity *domain.Entity) error {
	model := toEntityModel(entity)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	entity.ID = model.ID
	return nil
}

func (r *entityRepository) Update(entity *domain.Entity) error {
	model := toEntityModel(entity)
	return r.db.Save(model).Error
}

func (r *entityRepository) FindByID(id uint) (*domain.Entity, error) {
	var model models.Entity
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainEntity(&model), nil
}

func (r *entityRepository) FindByEmail(email string) (*domain.Entity, error) {
	var model models.Entity
	result := r.db.Where("email = ?", email).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainEntity(&model), nil
}

func (r *entityRepository) FindAll(limit int, offset int) ([]*domain.Entity, error) {
	var modelList []models.Entity
	result := r.db.Limit(limit).Offset(offset).Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	entities := make([]*domain.Entity, len(modelList))
	for i, model := range modelList {
		entities[i] = toDomainEntity(&model)
	}
	return entities, nil
}

func (r *entityRepository) Delete(id uint) error {
	return r.db.Delete(&models.Entity{}, id).Error
}

func (r *entityRepository) Exists(id uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.Entity{}).Where("id = ?", id).Count(&count)
	return count > 0, result.Error
}

// Mappers

func toEntityModel(e *domain.Entity) *models.Entity {
	return &models.Entity{
		ID:             e.ID,
		Type:           models.EntityType(e.Type),
		FirstName:      e.FirstName,
		LastName:       e.LastName,
		BusinessName:   e.BusinessName,
		TaxID:          e.TaxID,
		Email:          e.Email,
		Address:        e.Address,
		City:           e.City,
		State:          e.State,
		Zip:            e.Zip,
		CountryID:      e.CountryID,
		IsSystemUser:   e.IsSystemUser,
		IsInternal:     e.IsInternal,
		ParentEntityID: e.ParentEntityID,
		Status:         models.EntityStatus(e.Status),
		CreatedAt:      e.CreatedAt,
		ModifiedAt:     e.ModifiedAt,
	}
}

func toDomainEntity(m *models.Entity) *domain.Entity {
	return &domain.Entity{
		ID:             m.ID,
		Type:           domain.EntityType(m.Type),
		FirstName:      m.FirstName,
		LastName:       m.LastName,
		BusinessName:   m.BusinessName,
		TaxID:          m.TaxID,
		Email:          m.Email,
		Address:        m.Address,
		City:           m.City,
		State:          m.State,
		Zip:            m.Zip,
		CountryID:      m.CountryID,
		IsSystemUser:   m.IsSystemUser,
		IsInternal:     m.IsInternal,
		ParentEntityID: m.ParentEntityID,
		Status:         domain.EntityStatus(m.Status),
		CreatedAt:      m.CreatedAt,
		ModifiedAt:     m.ModifiedAt,
	}
}