package domain

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	sharedDomain "torque-dms/core/shared/domain"
)

type UserAccount struct {
	ID           uint
	EntityID     uint
	Username     string
	PasswordHash string
	LastLogin    time.Time
	Status       EntityStatus
	CreatedAt    time.Time
}

func NewUserAccount(entityID uint, username string, password string) (*UserAccount, error) {
	if err := sharedDomain.Validate("entity_id", fmt.Sprint(entityID)); err != nil {
		return nil, err
	}

	// Validar username
	if err := sharedDomain.Validate("username", username); err != nil {
		return nil, err
	}

	// Validar password
	if err := sharedDomain.Validate("password", password); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	return &UserAccount{
		EntityID:     entityID,
		Username:     username,
		PasswordHash: string(hash),
		Status:       EntityStatusActive,
		CreatedAt:    time.Now(),
	}, nil
}

func (u *UserAccount) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *UserAccount) ChangePassword(oldPassword string, newPassword string) error {
	if !u.CheckPassword(oldPassword) {
		return errors.New("incorrect current password")
	}

	if len(newPassword) < 8 {
		return errors.New("new password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	u.PasswordHash = string(hash)
	return nil
}

func (u *UserAccount) RecordLogin() {
	u.LastLogin = time.Now()
}

func (u *UserAccount) IsActive() bool {
	return u.Status == EntityStatusActive
}

func (u *UserAccount) Suspend() {
	u.Status = EntityStatusSuspended
}

func (u *UserAccount) Activate() {
	u.Status = EntityStatusActive
}