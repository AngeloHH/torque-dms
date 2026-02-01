package output

import "torque-dms/core/identity/domain"

type UserRepository interface {
	Save(user *domain.UserAccount) error
	Update(user *domain.UserAccount) error
	FindByID(id uint) (*domain.UserAccount, error)
	FindByEntityID(entityID uint) (*domain.UserAccount, error)
	FindByUsername(username string) (*domain.UserAccount, error)
	Delete(id uint) error
	Exists(username string) (bool, error)
}