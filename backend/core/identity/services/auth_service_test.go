package services

import (
	"errors"
	"testing"

	"torque-dms/core/identity/domain"
	"torque-dms/core/identity/ports/input"
	"torque-dms/core/identity/ports/output"
	sharedDomain "torque-dms/core/shared/domain"
)

func init() {
	sharedDomain.LoadValidationRules("../../../settings/validation_rules.yml")
}

// --- Mocks ---

type mockEntityRepo struct {
	saved    *domain.Entity
	saveErr  error
	findByID *domain.Entity
	findErr  error
}

func (m *mockEntityRepo) Save(e *domain.Entity) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	e.ID = 1
	m.saved = e
	return nil
}
func (m *mockEntityRepo) Update(e *domain.Entity) error { return nil }
func (m *mockEntityRepo) FindByID(id uint) (*domain.Entity, error) {
	if m.findByID != nil {
		return m.findByID, nil
	}
	return nil, m.findErr
}
func (m *mockEntityRepo) FindByEmail(email string) (*domain.Entity, error)        { return nil, nil }
func (m *mockEntityRepo) FindAll(limit int, offset int) ([]*domain.Entity, error) { return nil, nil }
func (m *mockEntityRepo) Delete(id uint) error                                    { return nil }
func (m *mockEntityRepo) Exists(id uint) (bool, error)                            { return false, nil }

type mockUserRepo struct {
	saved       *domain.UserAccount
	saveErr     error
	existsVal   bool
	existsErr   error
	findByUser  *domain.UserAccount
	findByIDVal *domain.UserAccount
	findErr     error
	updated     *domain.UserAccount
}

func (m *mockUserRepo) Save(u *domain.UserAccount) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	u.ID = 1
	m.saved = u
	return nil
}
func (m *mockUserRepo) Update(u *domain.UserAccount) error {
	m.updated = u
	return nil
}
func (m *mockUserRepo) FindByID(id uint) (*domain.UserAccount, error) {
	if m.findByIDVal != nil {
		return m.findByIDVal, nil
	}
	return nil, errors.New("not found")
}
func (m *mockUserRepo) FindByEntityID(entityID uint) (*domain.UserAccount, error) { return nil, nil }
func (m *mockUserRepo) FindByUsername(username string) (*domain.UserAccount, error) {
	if m.findByUser != nil {
		return m.findByUser, nil
	}
	return nil, m.findErr
}
func (m *mockUserRepo) Delete(id uint) error { return nil }
func (m *mockUserRepo) Exists(username string) (bool, error) {
	return m.existsVal, m.existsErr
}

type mockPhoneRepo struct {
	saved *output.Phone
}

func (m *mockPhoneRepo) Save(p *output.Phone) error                                 { m.saved = p; return nil }
func (m *mockPhoneRepo) Update(p *output.Phone) error                               { return nil }
func (m *mockPhoneRepo) FindByID(id uint) (*output.Phone, error)                    { return nil, nil }
func (m *mockPhoneRepo) FindByEntityID(entityID uint) ([]*output.Phone, error)      { return nil, nil }
func (m *mockPhoneRepo) FindPrimaryByEntityID(entityID uint) (*output.Phone, error) { return nil, nil }
func (m *mockPhoneRepo) Delete(id uint) error                                       { return nil }

// --- Tests ---

func TestAuthService_Register_Success(t *testing.T) {
	entityRepo := &mockEntityRepo{}
	userRepo := &mockUserRepo{existsVal: false}
	phoneRepo := &mockPhoneRepo{}

	svc := NewAuthService(entityRepo, userRepo, phoneRepo, "test-secret")

	entity, user, err := svc.Register(input.RegisterInput{
		Type:      "person",
		FirstName: "Juan",
		LastName:  "Perez",
		Email:     "juan@test.com",
		Username:  "juanperez",
		Password:  "Strong@123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entity == nil || user == nil {
		t.Fatal("entity and user should not be nil")
	}
	if entity.FirstName != "Juan" {
		t.Errorf("FirstName = %v, want Juan", entity.FirstName)
	}
	if !entity.IsSystemUser {
		t.Error("registered entity should be system user")
	}
}

func TestAuthService_Register_DuplicateUsername(t *testing.T) {
	entityRepo := &mockEntityRepo{}
	userRepo := &mockUserRepo{existsVal: true}
	phoneRepo := &mockPhoneRepo{}

	svc := NewAuthService(entityRepo, userRepo, phoneRepo, "test-secret")

	_, _, err := svc.Register(input.RegisterInput{
		Type:     "person",
		Email:    "juan@test.com",
		Username: "juanperez",
		Password: "Strong@123",
	})

	if err == nil {
		t.Error("expected error for duplicate username")
	}
}

func TestAuthService_Register_InvalidEmail(t *testing.T) {
	entityRepo := &mockEntityRepo{}
	userRepo := &mockUserRepo{existsVal: false}
	phoneRepo := &mockPhoneRepo{}

	svc := NewAuthService(entityRepo, userRepo, phoneRepo, "test-secret")

	_, _, err := svc.Register(input.RegisterInput{
		Type:     "person",
		Email:    "invalid",
		Username: "juanperez",
		Password: "Strong@123",
	})

	if err == nil {
		t.Error("expected error for invalid email")
	}
}

func TestAuthService_Register_WithPhone(t *testing.T) {
	entityRepo := &mockEntityRepo{}
	userRepo := &mockUserRepo{existsVal: false}
	phoneRepo := &mockPhoneRepo{}

	svc := NewAuthService(entityRepo, userRepo, phoneRepo, "test-secret")

	_, _, err := svc.Register(input.RegisterInput{
		Type:     "person",
		Phone:    "+1 7875551234",
		Username: "juanperez",
		Password: "Strong@123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if phoneRepo.saved == nil {
		t.Error("phone should have been saved")
	}
	if !phoneRepo.saved.IsPrimary {
		t.Error("phone should be primary")
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	user, _ := domain.NewUserAccount(1, "juanperez", "Strong@123")
	user.ID = 1

	entityRepo := &mockEntityRepo{}
	userRepo := &mockUserRepo{findByUser: user}
	phoneRepo := &mockPhoneRepo{}

	svc := NewAuthService(entityRepo, userRepo, phoneRepo, "test-secret")

	result, err := svc.Login(input.LoginInput{
		Username: "juanperez",
		Password: "Strong@123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Token == "" {
		t.Error("token should not be empty")
	}
	if result.User.ID != 1 {
		t.Errorf("user ID = %v, want 1", result.User.ID)
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	user, _ := domain.NewUserAccount(1, "juanperez", "Strong@123")

	entityRepo := &mockEntityRepo{}
	userRepo := &mockUserRepo{findByUser: user}
	phoneRepo := &mockPhoneRepo{}

	svc := NewAuthService(entityRepo, userRepo, phoneRepo, "test-secret")

	_, err := svc.Login(input.LoginInput{
		Username: "juanperez",
		Password: "WrongPass@1",
	})

	if err == nil {
		t.Error("expected error for wrong password")
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	entityRepo := &mockEntityRepo{}
	userRepo := &mockUserRepo{findErr: errors.New("not found")}
	phoneRepo := &mockPhoneRepo{}

	svc := NewAuthService(entityRepo, userRepo, phoneRepo, "test-secret")

	_, err := svc.Login(input.LoginInput{
		Username: "nonexistent",
		Password: "Strong@123",
	})

	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}

func TestAuthService_Login_SuspendedAccount(t *testing.T) {
	user, _ := domain.NewUserAccount(1, "juanperez", "Strong@123")
	user.Suspend()

	entityRepo := &mockEntityRepo{}
	userRepo := &mockUserRepo{findByUser: user}
	phoneRepo := &mockPhoneRepo{}

	svc := NewAuthService(entityRepo, userRepo, phoneRepo, "test-secret")

	_, err := svc.Login(input.LoginInput{
		Username: "juanperez",
		Password: "Strong@123",
	})

	if err == nil {
		t.Error("expected error for suspended account")
	}
}

func TestAuthService_ChangePassword_Success(t *testing.T) {
	user, _ := domain.NewUserAccount(1, "juanperez", "Strong@123")
	user.ID = 5

	entityRepo := &mockEntityRepo{}
	userRepo := &mockUserRepo{findByIDVal: user}
	phoneRepo := &mockPhoneRepo{}

	svc := NewAuthService(entityRepo, userRepo, phoneRepo, "test-secret")

	err := svc.ChangePassword(input.ChangePasswordInput{
		UserID:      5,
		OldPassword: "Strong@123",
		NewPassword: "NewStrong@1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuthService_ChangePassword_WrongOld(t *testing.T) {
	user, _ := domain.NewUserAccount(1, "juanperez", "Strong@123")
	user.ID = 5

	entityRepo := &mockEntityRepo{}
	userRepo := &mockUserRepo{findByIDVal: user}
	phoneRepo := &mockPhoneRepo{}

	svc := NewAuthService(entityRepo, userRepo, phoneRepo, "test-secret")

	err := svc.ChangePassword(input.ChangePasswordInput{
		UserID:      5,
		OldPassword: "WrongOld@1",
		NewPassword: "NewStrong@1",
	})

	if err == nil {
		t.Error("expected error for wrong old password")
	}
}
