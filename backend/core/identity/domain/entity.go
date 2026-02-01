package domain

import (
	"errors"
	"strings"
	"time"
	"regexp"
)

// Tipos
type EntityType string

const (
	EntityTypePerson       EntityType = "person"
	EntityTypeCompany      EntityType = "company"
	EntityTypeDealer       EntityType = "dealer"
	EntityTypeOrganization EntityType = "organization"
)

type EntityStatus string

const (
	EntityStatusActive    EntityStatus = "active"
	EntityStatusInactive  EntityStatus = "inactive"
	EntityStatusSuspended EntityStatus = "suspended"
)

// La entidad de dominio - representa qué ES un Entity en tu negocio
type Entity struct {
	ID             uint
	Type           EntityType
	FirstName      string
	LastName       string
	BusinessName   string
	TaxID          string
	Email          string
	Address        string
	City           string
	State          string
	Zip            string
	CountryID      *uint
	IsSystemUser   bool
	IsInternal     bool
	ParentEntityID *uint
	Status         EntityStatus
	CreatedAt      time.Time
	ModifiedAt     time.Time
}

// Constructor - crea un Entity validando las reglas de negocio
func NewEntity(entityType EntityType, phone string, email string) (*Entity, error) {

	// Si ambos están vacíos
	if email == "" && phone == "" {
		return nil, errors.New("either email or phone is required")
	}

	// Si hay email pero es inválido
	if email != "" && !isValidEmail(email) {
		return nil, errors.New("invalid email format")
	}

	// Si hay phone pero es inválido
	if phone != "" && !isValidPhone(phone) {
		return nil, errors.New("invalid phone format")
	}

	return &Entity{
		Type:      entityType,
		Email:     strings.ToLower(email),
		Status:    EntityStatusActive,
		CreatedAt: time.Now(),
	}, nil
}

func (e *Entity) SetField(field string, value string) error {
	fields := map[string]*string{
		"first_name":    &e.FirstName,
		"last_name":     &e.LastName,
		"business_name": &e.BusinessName,
		"tax_id":        &e.TaxID,
		"email":         &e.Email,
		"address":       &e.Address,
		"city":          &e.City,
		"state":         &e.State,
		"zip":           &e.Zip,
	}

	for key, ptr := range fields {
		if key == "email" && value != "" && !isValidEmail(value) {
			return errors.New("invalid email format")
		}

		if key == field {
			*ptr = value
			e.ModifiedAt = time.Now()
			return nil
		}
	}

	return errors.New("invalid field")
}

func (e *Entity) SetAsSystemUser() {
	e.IsSystemUser = true
	e.ModifiedAt = time.Now()
}

func (e *Entity) SetAsInternal() {
	e.IsInternal = true
	e.ModifiedAt = time.Now()
}

func (e *Entity) Suspend() error {
	if e.Status == EntityStatusSuspended {
		return errors.New("entity is already suspended")
	}
	e.Status = EntityStatusSuspended
	e.ModifiedAt = time.Now()
	return nil
}

func (e *Entity) Activate() error {
	if e.Status == EntityStatusActive {
		return errors.New("entity is already active")
	}
	e.Status = EntityStatusActive
	e.ModifiedAt = time.Now()
	return nil
}

func (e *Entity) Deactivate() error {
	if e.Status == EntityStatusInactive {
		return errors.New("entity is already inactive")
	}
	e.Status = EntityStatusInactive
	e.ModifiedAt = time.Now()
	return nil
}

// Consultas de estado

func (e *Entity) IsActive() bool {
	return e.Status == EntityStatusActive
}

func (e *Entity) IsPerson() bool {
	return e.Type == EntityTypePerson
}

func (e *Entity) IsCompany() bool {
	return e.Type == EntityTypeCompany
}

func (e *Entity) CanLogin() bool {
	return e.IsSystemUser && e.IsActive()
}

// Helpers privados
func isValidEmail(email string) bool {
	if len(email) > 254 {return false}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(strings.ToLower(email))
}

func isValidPhone(phone string) bool {
	if phone == "" {return false}
	phoneRegex := regexp.MustCompile(`^\+\d{1,4}\s\d{6,14}$`)
	return phoneRegex.MatchString(phone)
}