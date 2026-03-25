package service

import (
	"errors"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/service/dto"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

/*
UNIT TESTING CONCEPTS (BACKEND - SERVICE LAYER)

1) What is unit testing in backend?
   - Unit testing focuses on one business unit (here: UserService methods) in isolation.
   - It avoids external systems so behavior is fast, repeatable, and easy to debug.

2) Why service layer tests first?
   - Service layer holds business rules like password verification, role normalization,
     and validation decisions.
   - Controller layer mostly handles HTTP transport concerns.
   - Real DB behavior is better covered in integration tests.

3) What is mocking and why use it?
   - Mocking means replacing real repository with a fake struct that returns controlled data.
   - This lets us simulate success, edge, and failure paths without using a real database.

4) Unit vs Integration vs E2E
   - Unit: test one function with mocked dependencies.
   - Integration: test multiple modules together (often with real DB).
   - E2E: test complete workflow from client request to persistence and response.
*/

type mockUserRepository struct {
	getUserResult domain.User
	getUserErr    error
	requestedName string

	updateRoleErr      error
	capturedUserID     int64
	capturedRole       string
	updateUserRoleCall int
}

func (mockRepo *mockUserRepository) EnsureUserSchema() error {
	return nil
}

func (mockRepo *mockUserRepository) GetUserByUsername(username string) (domain.User, error) {
	mockRepo.requestedName = username
	return mockRepo.getUserResult, mockRepo.getUserErr
}

func (mockRepo *mockUserRepository) SignUp(user domain.User) error {
	return nil
}

func (mockRepo *mockUserRepository) GetAllUsers() ([]dto.UserSummary, error) {
	return []dto.UserSummary{}, nil
}

func (mockRepo *mockUserRepository) UpdateUserRole(userID int64, role string) error {
	mockRepo.updateUserRoleCall++
	mockRepo.capturedUserID = userID
	mockRepo.capturedRole = role
	return mockRepo.updateRoleErr
}

func (mockRepo *mockUserRepository) EnsureAdminExists() error {
	return nil
}

// TestUserService_Login_Success verifies valid credentials produce token and normalized role.
// Why important: login is a critical security workflow and must issue auth tokens correctly.
// Validates: password verification logic, token creation, and role normalization in service layer.
// Setup: create bcrypt hash for known password, mock repo returns user, set JWT_KEY env.
// Execution: call service.Login with correct credentials.
// Assertion: check no error, token exists, role is lowercase, and username lookup was correct.
func TestUserService_Login_Success(t *testing.T) {
	t.Setenv("JWT_KEY", "unit-test-secret")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Test1234"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password for setup: %v", err)
	}

	mockRepo := &mockUserRepository{
		getUserResult: domain.User{
			Id:       10,
			Username: "erkin",
			Password: string(hashedPassword),
			Role:     "Admin",
		},
	}
	service := NewUserService(mockRepo)

	result, err := service.Login("erkin", "Test1234")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatalf("expected non-nil login result")
	}
	if result.Token == "" {
		t.Fatalf("expected non-empty JWT token")
	}
	if result.Role != "admin" {
		t.Fatalf("expected role admin, got %s", result.Role)
	}
	if mockRepo.requestedName != "erkin" {
		t.Fatalf("expected repository lookup for erkin, got %s", mockRepo.requestedName)
	}
}

// TestUserService_Login_EdgeCase_UserNotFound validates user-not-found handling.
// Why important: clear failure behavior is essential for auth flows and UI messaging.
// Validates: service error branch when repository cannot find user.
// Setup: mock repo returns error on username lookup.
// Execution: call service.Login with unknown username.
// Assertion: expect an error and no result object.
func TestUserService_Login_EdgeCase_UserNotFound(t *testing.T) {
	t.Setenv("JWT_KEY", "unit-test-secret")

	mockRepo := &mockUserRepository{getUserErr: errors.New("not found")}
	service := NewUserService(mockRepo)

	result, err := service.Login("missing-user", "anything")

	if err == nil {
		t.Fatalf("expected error when user is not found")
	}
	if result != nil {
		t.Fatalf("expected nil result on user-not-found")
	}
	if !strings.Contains(err.Error(), "no user found") {
		t.Fatalf("expected error to mention no user found, got: %v", err)
	}
}

// TestUserService_Login_Failure_InvalidPassword validates password mismatch behavior.
// Why important: ensures security rule rejects invalid credentials.
// Validates: bcrypt comparison branch in service logic.
// Setup: mock repo returns valid user with hash for a different password.
// Execution: call Login with wrong password.
// Assertion: expect invalid password error and nil login result.
func TestUserService_Login_Failure_InvalidPassword(t *testing.T) {
	t.Setenv("JWT_KEY", "unit-test-secret")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("CorrectPassword"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password for setup: %v", err)
	}

	mockRepo := &mockUserRepository{
		getUserResult: domain.User{
			Id:       11,
			Username: "erkin",
			Password: string(hashedPassword),
			Role:     "user",
		},
	}
	service := NewUserService(mockRepo)

	result, err := service.Login("erkin", "WrongPassword")

	if err == nil {
		t.Fatalf("expected invalid password error, got nil")
	}
	if result != nil {
		t.Fatalf("expected nil result when password is invalid")
	}
	if !strings.Contains(err.Error(), "invalid password") {
		t.Fatalf("expected invalid password message, got: %v", err)
	}
}

// TestUserService_UpdateUserRole_Success_NormalizedRole verifies role normalization before repository call.
// Why important: guarantees consistent role values in storage (admin/user lowercase).
// Validates: role business rule (trim + lowercase + allowed values).
// Setup: mock repository that records input role.
// Execution: call UpdateUserRole with "  ADMIN  ".
// Assertion: expect no error, repository called once, role passed as "admin".
func TestUserService_UpdateUserRole_Success_NormalizedRole(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewUserService(mockRepo)

	err := service.UpdateUserRole(42, "  ADMIN  ")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if mockRepo.updateUserRoleCall != 1 {
		t.Fatalf("expected UpdateUserRole repository call once, got %d", mockRepo.updateUserRoleCall)
	}
	if mockRepo.capturedUserID != 42 {
		t.Fatalf("expected userID 42, got %d", mockRepo.capturedUserID)
	}
	if mockRepo.capturedRole != "admin" {
		t.Fatalf("expected normalized role admin, got %s", mockRepo.capturedRole)
	}
}

// TestUserService_UpdateUserRole_EdgeCase_InvalidRole checks validation for unsupported roles.
// Why important: prevents invalid authorization states from being stored.
// Validates: service-side role guard clause.
// Setup: mock repository (should not be called).
// Execution: call UpdateUserRole with unsupported role "manager".
// Assertion: expect error and zero repository calls.
func TestUserService_UpdateUserRole_EdgeCase_InvalidRole(t *testing.T) {
	mockRepo := &mockUserRepository{}
	service := NewUserService(mockRepo)

	err := service.UpdateUserRole(42, "manager")

	if err == nil {
		t.Fatalf("expected invalid role error, got nil")
	}
	if mockRepo.updateUserRoleCall != 0 {
		t.Fatalf("expected repository not to be called for invalid role")
	}
}

// TestUserService_UpdateUserRole_Failure_RepositoryError validates propagation of repo failure.
// Why important: service should expose persistence failure so caller can handle it.
// Validates: failure handling for role update path.
// Setup: mock repository returns error.
// Execution: call UpdateUserRole with valid role.
// Assertion: expect error from service.
func TestUserService_UpdateUserRole_Failure_RepositoryError(t *testing.T) {
	mockRepo := &mockUserRepository{updateRoleErr: errors.New("update failed")}
	service := NewUserService(mockRepo)

	err := service.UpdateUserRole(42, "user")

	if err == nil {
		t.Fatalf("expected repository error, got nil")
	}
}
