package input

import "torque-dms/core/identity/domain"

type CreateEntityInput struct {
	Type         string
	FirstName    string
	LastName     string
	BusinessName string
	TaxID        string
	Email        string
	Phone        string
	Address      string
	City         string
	State        string
	Zip          string
	CountryID    *uint
}

type UpdateEntityInput struct {
	Field string
	Value string
}

type EntityService interface {
	Create(input CreateEntityInput) (*domain.Entity, error)
	GetByID(id uint) (*domain.Entity, error)
	GetByEmail(email string) (*domain.Entity, error)
	Update(id uint, input UpdateEntityInput) (*domain.Entity, error)
	Delete(id uint) error
	List(limit int, offset int) ([]*domain.Entity, error)
	Suspend(id uint) error
	Activate(id uint) error
}