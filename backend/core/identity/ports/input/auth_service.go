package input

import "torque-dms/core/identity/domain"

type RegisterInput struct {
	Type         string
	FirstName    string
	LastName     string
	BusinessName string
	Email        string
	Phone        string
	Username     string
	Password     string
}

type LoginInput struct {
	Username string
	Password string
}

type LoginOutput struct {
	User  *domain.UserAccount
	Token string
}

type ChangePasswordInput struct {
	UserID      uint
	OldPassword string
	NewPassword string
}

type AuthService interface {
	Register(input RegisterInput) (*domain.Entity, *domain.UserAccount, error)
	Login(input LoginInput) (*LoginOutput, error)
	ChangePassword(input ChangePasswordInput) error
	Logout(userID uint) error
}