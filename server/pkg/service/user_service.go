package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/repository"
	"ims-intro/pkg/service/dto"
	"os"
	"time"
)

type LoginResult struct {
	Token string
	Role  string
}

type IUserService interface {
	Login(username, password string) (string, error)
	SignUp(user dto.UserCreate) error
	GetAllUsers() ([]dto.UserSummary, error)
	UpdateUserRole(userID int64, role string) error
}

type UserService struct {
	userRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) IUserService {
	return &UserService{userRepository}
}

func (service *UserService) Login(username, password string) (string, error) {
	jwtKey := os.Getenv("JWT_KEY")

	user, err := service.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("no user found with the username: " + username)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &domain.Claims{
		UserID: user.Id,
		Role:   strings.ToLower(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", errors.New("error signing the token: " + err.Error())
	}

	return &LoginResult{Token: tokenString, Role: strings.ToLower(user.Role)}, nil
}

func (service *UserService) SignUp(userCreate dto.UserCreate) error {
	userCreate.Role = "user"
	err := validateUserCreate(userCreate)
	if err != nil {
		return err
	}

	user := userCreateToUser(userCreate)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error while creating password hash")
	}

	user.Password = string(hashedPassword)

	return service.userRepository.SignUp(user)
}

func validateUserCreate(u dto.UserCreate) error {
	if u.Username == "" {
		return errors.New("username can't be empty")
	}
	if u.Password == "" {
		return errors.New("category can't be empty")
	}
	return nil
}

func userCreateToUser(userCreate dto.UserCreate) domain.User {
	return domain.User{
		Username: userCreate.Username,
		Email:    userCreate.Email,
		Password: userCreate.Password,
		Role:     userCreate.Role,
	}
}
