package output

import "torque-dms/core/identity/domain"

type EntityRepository interface {
	Save(entity *domain.Entity) error
	Update(entity *domain.Entity) error
	FindByID(id uint) (*domain.Entity, error)
	FindByEmail(email string) (*domain.Entity, error)
	FindAll(limit int, offset int) ([]*domain.Entity, error)
	Delete(id uint) error
	Exists(id uint) (bool, error)
}