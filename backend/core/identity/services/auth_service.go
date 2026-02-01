package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"torque-dms/core/identity/domain"
	"torque-dms/core/identity/ports/input"
	"torque-dms/core/identity/ports/output"
)

type authService struct {
	entityRepo output.EntityRepository
	userRepo   output.UserRepository
	phoneRepo  output.PhoneRepository
	jwtSecret  string
}

func NewAuthService(
	entityRepo output.EntityRepository,
	userRepo output.UserRepository,
	phoneRepo output.PhoneRepository,
	jwtSecret string,
) input.AuthService {
	return &authService{
		entityRepo: entityRepo,
		userRepo:   userRepo,
		phoneRepo:  phoneRepo,
		jwtSecret:  jwtSecret,
	}
}

func (s *authService) Register(inp input.RegisterInput) (*domain.Entity, *domain.UserAccount, error) {
	// Verificar que username no exista
	exists, err := s.userRepo.Exists(inp.Username)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, errors.New("username already exists")
	}

	// Crear entity
	entity, err := domain.NewEntity(domain.EntityType(inp.Type), inp.Phone, inp.Email)
	if err != nil {
		return nil, nil, err
	}

	entity.SetField("first_name", inp.FirstName)
	entity.SetField("last_name", inp.LastName)
	entity.SetField("business_name", inp.BusinessName)
	entity.SetAsSystemUser()

	if err := s.entityRepo.Save(entity); err != nil {
		return nil, nil, err
	}

	// Si hay teléfono, guardarlo
	if inp.Phone != "" {
		phone := &output.Phone{
			EntityID:  entity.ID,
			Number:    inp.Phone,
			IsPrimary: true,
		}
		if err := s.phoneRepo.Save(phone); err != nil {
			return nil, nil, err
		}
	}

	// Crear user account
	user, err := domain.NewUserAccount(entity.ID, inp.Username, inp.Password)
	if err != nil {
		return nil, nil, err
	}

	if err := s.userRepo.Save(user); err != nil {
		return nil, nil, err
	}

	return entity, user, nil
}

func (s *authService) Login(inp input.LoginInput) (*input.LoginOutput, error) {
	user, err := s.userRepo.FindByUsername(inp.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.CheckPassword(inp.Password) {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive() {
		return nil, errors.New("account is not active")
	}

	// Generar JWT
	token, err := s.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Registrar login
	user.RecordLogin()
	s.userRepo.Update(user)

	return &input.LoginOutput{
		User:  user,
		Token: token,
	}, nil
}

func (s *authService) ChangePassword(inp input.ChangePasswordInput) error {
	user, err := s.userRepo.FindByID(inp.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := user.ChangePassword(inp.OldPassword, inp.NewPassword); err != nil {
		return err
	}

	return s.userRepo.Update(user)
}

func (s *authService) Logout(userID uint) error {
	// Por ahora solo retornamos nil
	// En el futuro podrías invalidar el token en una blacklist
	return nil
}

func (s *authService) generateToken(user *domain.UserAccount) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"entity_id": user.EntityID,
		"username":  user.Username,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}