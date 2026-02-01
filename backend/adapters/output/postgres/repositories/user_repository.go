package repositories

import (
	"gorm.io/gorm"
	"torque-dms/core/identity/domain"
	"torque-dms/core/identity/ports/output"
	"torque-dms/models"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) output.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(user *domain.UserAccount) error {
	model := toUserModel(user)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	user.ID = model.ID
	return nil
}

func (r *userRepository) Update(user *domain.UserAccount) error {
	model := toUserModel(user)
	return r.db.Save(model).Error
}

func (r *userRepository) FindByID(id uint) (*domain.UserAccount, error) {
	var model models.UserAccount
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainUser(&model), nil
}

func (r *userRepository) FindByEntityID(entityID uint) (*domain.UserAccount, error) {
	var model models.UserAccount
	result := r.db.Where("entity_id = ?", entityID).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainUser(&model), nil
}

func (r *userRepository) FindByUsername(username string) (*domain.UserAccount, error) {
	var model models.UserAccount
	result := r.db.Where("username = ?", username).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomainUser(&model), nil
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserAccount{}, id).Error
}

func (r *userRepository) Exists(username string) (bool, error) {
	var count int64
	result := r.db.Model(&models.UserAccount{}).Where("username = ?", username).Count(&count)
	return count > 0, result.Error
}

// Mappers

func toUserModel(u *domain.UserAccount) *models.UserAccount {
	return &models.UserAccount{
		ID:           u.ID,
		EntityID:     u.EntityID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		LastLogin:    u.LastLogin,
		Status:       models.EntityStatus(u.Status),
		CreatedAt:    u.CreatedAt,
	}
}

func toDomainUser(m *models.UserAccount) *domain.UserAccount {
	return &domain.UserAccount{
		ID:           m.ID,
		EntityID:     m.EntityID,
		Username:     m.Username,
		PasswordHash: m.PasswordHash,
		LastLogin:    m.LastLogin,
		Status:       domain.EntityStatus(m.Status),
		CreatedAt:    m.CreatedAt,
	}
}