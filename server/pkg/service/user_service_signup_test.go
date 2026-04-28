package service

import (
	"errors"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/service/dto"
)

type signupUserRepositoryMock struct {
	signUpCalled int
	capturedUser domain.User
	signUpErr    error
}

func (m *signupUserRepositoryMock) EnsureUserSchema() error { return nil }
func (m *signupUserRepositoryMock) GetUserByUsername(username string) (domain.User, error) {
	return domain.User{}, errors.New("not used")
}
func (m *signupUserRepositoryMock) SignUp(user domain.User) error {
	m.signUpCalled++
	m.capturedUser = user
	return m.signUpErr
}
func (m *signupUserRepositoryMock) GetAllUsers() ([]dto.UserSummary, error) { return nil, nil }
func (m *signupUserRepositoryMock) UpdateUserRole(userID int64, role string) error { return nil }
func (m *signupUserRepositoryMock) EnsureAdminExists() error { return nil }

func TestUserService_SignUp_Success_ShouldHashPasswordAndForceUserRole(t *testing.T) {
	mockRepo := &signupUserRepositoryMock{}
	serviceInstance := NewUserService(mockRepo)

	err := serviceInstance.SignUp(dto.UserCreate{Username: "erkin", Email: "erkin@example.com", Password: "Test1234", Role: "admin"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if mockRepo.signUpCalled != 1 {
		t.Fatalf("expected repository SignUp call once, got %d", mockRepo.signUpCalled)
	}
	if mockRepo.capturedUser.Role != "user" {
		t.Fatalf("expected role user, got %s", mockRepo.capturedUser.Role)
	}
	if mockRepo.capturedUser.Password == "Test1234" {
		t.Fatalf("expected password to be hashed")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(mockRepo.capturedUser.Password), []byte("Test1234")); err != nil {
		t.Fatalf("expected hashed password to verify, got %v", err)
	}
}

func TestUserService_SignUp_DuplicateUser_ShouldFail(t *testing.T) {
	mockRepo := &signupUserRepositoryMock{signUpErr: errors.New("duplicate user")}
	serviceInstance := NewUserService(mockRepo)

	err := serviceInstance.SignUp(dto.UserCreate{Username: "erkin", Email: "erkin@example.com", Password: "Test1234"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestUserService_SignUp_MissingUsername_ShouldFail(t *testing.T) {
	mockRepo := &signupUserRepositoryMock{}
	serviceInstance := NewUserService(mockRepo)

	err := serviceInstance.SignUp(dto.UserCreate{Email: "erkin@example.com", Password: "Test1234"})
	if err == nil || !strings.Contains(err.Error(), "username") {
		t.Fatalf("expected username validation error, got %v", err)
	}
	if mockRepo.signUpCalled != 0 {
		t.Fatalf("expected repository not to be called")
	}
}

func TestUserService_SignUp_MissingPassword_ShouldFail(t *testing.T) {
	mockRepo := &signupUserRepositoryMock{}
	serviceInstance := NewUserService(mockRepo)

	err := serviceInstance.SignUp(dto.UserCreate{Username: "erkin", Email: "erkin@example.com"})
	if err == nil || !strings.Contains(err.Error(), "password") {
		t.Fatalf("expected password validation error, got %v", err)
	}
	if mockRepo.signUpCalled != 0 {
		t.Fatalf("expected repository not to be called")
	}
}
