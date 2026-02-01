package services

import (
	"errors"

	"torque-dms/core/identity/domain"
	"torque-dms/core/identity/ports/input"
	"torque-dms/core/identity/ports/output"
)

type entityService struct {
	entityRepo output.EntityRepository
	phoneRepo  output.PhoneRepository
}

func NewEntityService(entityRepo output.EntityRepository, phoneRepo output.PhoneRepository) input.EntityService {
	return &entityService{
		entityRepo: entityRepo,
		phoneRepo:  phoneRepo,
	}
}

func (s *entityService) Create(inp input.CreateEntityInput) (*domain.Entity, error) {
	// Crear entity en domain
	entity, err := domain.NewEntity(domain.EntityType(inp.Type), inp.Phone, inp.Email)
	if err != nil {
		return nil, err
	}

	// Setear campos opcionales
	if inp.FirstName != "" {
		entity.SetField("first_name", inp.FirstName)
	}
	if inp.LastName != "" {
		entity.SetField("last_name", inp.LastName)
	}
	if inp.BusinessName != "" {
		entity.SetField("business_name", inp.BusinessName)
	}
	if inp.TaxID != "" {
		entity.SetField("tax_id", inp.TaxID)
	}
	if inp.Address != "" {
		entity.SetField("address", inp.Address)
	}
	if inp.City != "" {
		entity.SetField("city", inp.City)
	}
	if inp.State != "" {
		entity.SetField("state", inp.State)
	}
	if inp.Zip != "" {
		entity.SetField("zip", inp.Zip)
	}
	if inp.CountryID != nil {
		entity.CountryID = inp.CountryID
	}

	// Guardar entity
	if err := s.entityRepo.Save(entity); err != nil {
		return nil, err
	}

	// Si hay teléfono, guardarlo
	if inp.Phone != "" {
		phone := &output.Phone{
			EntityID:  entity.ID,
			Number:    inp.Phone,
			IsPrimary: true,
		}
		if err := s.phoneRepo.Save(phone); err != nil {
			return nil, err
		}
	}

	return entity, nil
}

func (s *entityService) GetByID(id uint) (*domain.Entity, error) {
	return s.entityRepo.FindByID(id)
}

func (s *entityService) GetByEmail(email string) (*domain.Entity, error) {
	return s.entityRepo.FindByEmail(email)
}

func (s *entityService) Update(id uint, inp input.UpdateEntityInput) (*domain.Entity, error) {
	entity, err := s.entityRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("entity not found")
	}

	if err := entity.SetField(inp.Field, inp.Value); err != nil {
		return nil, err
	}

	if err := s.entityRepo.Update(entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *entityService) Delete(id uint) error {
	exists, err := s.entityRepo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("entity not found")
	}

	return s.entityRepo.Delete(id)
}

func (s *entityService) List(limit int, offset int) ([]*domain.Entity, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.entityRepo.FindAll(limit, offset)
}

func (s *entityService) Suspend(id uint) error {
	entity, err := s.entityRepo.FindByID(id)
	if err != nil {
		return errors.New("entity not found")
	}

	if err := entity.Suspend(); err != nil {
		return err
	}

	return s.entityRepo.Update(entity)
}

func (s *entityService) Activate(id uint) error {
	entity, err := s.entityRepo.FindByID(id)
	if err != nil {
		return errors.New("entity not found")
	}

	if err := entity.Activate(); err != nil {
		return err
	}

	return s.entityRepo.Update(entity)
}