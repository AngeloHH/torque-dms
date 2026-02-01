package repositories

import (
	"gorm.io/gorm"
	"torque-dms/core/identity/ports/output"
	"torque-dms/models"
)

type phoneRepository struct {
	db *gorm.DB
}

func NewPhoneRepository(db *gorm.DB) output.PhoneRepository {
	return &phoneRepository{db: db}
}

func (r *phoneRepository) Save(phone *output.Phone) error {
	model := toPhoneModel(phone)
	result := r.db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	phone.ID = model.ID
	return nil
}

func (r *phoneRepository) Update(phone *output.Phone) error {
	model := toPhoneModel(phone)
	return r.db.Save(model).Error
}

func (r *phoneRepository) FindByID(id uint) (*output.Phone, error) {
	var model models.EntityPhone
	result := r.db.First(&model, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return toOutputPhone(&model), nil
}

func (r *phoneRepository) FindByEntityID(entityID uint) ([]*output.Phone, error) {
	var modelList []models.EntityPhone
	result := r.db.Where("entity_id = ?", entityID).Find(&modelList)
	if result.Error != nil {
		return nil, result.Error
	}

	phones := make([]*output.Phone, len(modelList))
	for i, model := range modelList {
		phones[i] = toOutputPhone(&model)
	}
	return phones, nil
}

func (r *phoneRepository) FindPrimaryByEntityID(entityID uint) (*output.Phone, error) {
	var model models.EntityPhone
	result := r.db.Where("entity_id = ? AND is_primary = ?", entityID, true).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toOutputPhone(&model), nil
}

func (r *phoneRepository) Delete(id uint) error {
	return r.db.Delete(&models.EntityPhone{}, id).Error
}

// Mappers

func toPhoneModel(p *output.Phone) *models.EntityPhone {
	return &models.EntityPhone{
		ID:        p.ID,
		EntityID:  p.EntityID,
		CountryID: p.CountryID,
		Number:    p.Number,
		Type:      models.PhoneType(p.Type),
		IsPrimary: p.IsPrimary,
		Verified:  p.Verified,
	}
}

func toOutputPhone(m *models.EntityPhone) *output.Phone {
	return &output.Phone{
		ID:        m.ID,
		EntityID:  m.EntityID,
		CountryID: m.CountryID,
		Number:    m.Number,
		Type:      string(m.Type),
		IsPrimary: m.IsPrimary,
		Verified:  m.Verified,
	}
}