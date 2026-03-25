package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/log"
	"ims-intro/pkg/domain"
	"ims-intro/pkg/service/dto"
	"strings"
)

type IUserRepository interface {
	EnsureUserSchema() error
	GetUserByUsername(username string) (domain.User, error)
	SignUp(user domain.User) error
	GetAllUsers() ([]dto.UserSummary, error)
	UpdateUserRole(userID int64, role string) error
	EnsureAdminExists() error
}

type UserRepository struct {
	dbPool *pgxpool.Pool
}

func NewUserRepository(dbPool *pgxpool.Pool) IUserRepository {
	return &UserRepository{dbPool}
}

func (repository *UserRepository) EnsureUserSchema() error {
	ctx := context.Background()

	_, err := repository.dbPool.Exec(ctx, "ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(255)")
	if err != nil {
		return err
	}

	_, err = repository.dbPool.Exec(ctx, "ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20)")
	if err != nil {
		return err
	}

	_, err = repository.dbPool.Exec(ctx, "ALTER TABLE users ALTER COLUMN role SET DEFAULT 'user'")
	if err != nil {
		return err
	}

	_, err = repository.dbPool.Exec(ctx, "UPDATE users SET role = 'user' WHERE role IS NULL OR role = ''")
	if err != nil {
		return err
	}

	return nil
}

func (repository *UserRepository) GetUserByUsername(username string) (domain.User, error) {
	ctx := context.Background()

	var user domain.User

	selectStatement := "SELECT id, username, COALESCE(email, ''), password, role FROM users WHERE username = $1"
	userRow := repository.dbPool.QueryRow(ctx, selectStatement, username)

	err := userRow.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Role)
	user.Role = strings.ToLower(user.Role)

	if err != nil && err.Error() == "no rows in result set" {
		return domain.User{}, errors.New("error while finding user")
	}

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (repository *UserRepository) SignUp(user domain.User) error {
	ctx := context.Background()

	insertStatement := "INSERT INTO users(username, email, password, role) VALUES ($1, $2, $3, $4)"

	role := strings.ToLower(strings.TrimSpace(user.Role))
	if role == "" {
		role = "user"
	}

	addNewUser, err := repository.dbPool.Exec(ctx, insertStatement, user.Username, user.Email, user.Password, role)
	if err != nil {
		log.Errorf("error while adding new user: %v", err)
		return err
	}

	log.Info(fmt.Sprint("User added successfully: %v", addNewUser))
	return nil
}

func (repository *UserRepository) GetAllUsers() ([]dto.UserSummary, error) {
	ctx := context.Background()

	rows, err := repository.dbPool.Query(ctx, "SELECT id, username, COALESCE(email, ''), role FROM users ORDER BY id ASC")
	if err != nil {
		log.Errorf("error while listing users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []dto.UserSummary
	for rows.Next() {
		var user dto.UserSummary
		err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role)
		if err != nil {
			return nil, err
		}
		user.Role = strings.ToLower(user.Role)
		users = append(users, user)
	}

	return users, nil
}

func (repository *UserRepository) UpdateUserRole(userID int64, role string) error {
	ctx := context.Background()

	_, err := repository.dbPool.Exec(ctx, "UPDATE users SET role = $1 WHERE id = $2", strings.ToLower(role), userID)
	if err != nil {
		log.Errorf("error while updating user role: %v", err)
		return err
	}

	return nil
}

func (repository *UserRepository) EnsureAdminExists() error {
	ctx := context.Background()

	const adminUsername = "admin"
	const adminEmail = "admin@inventory.com"
	const adminRole = "admin"
	const adminPasswordHash = "$2a$10$BvZM.DEP8RSRg/yurAqlVuPlXpamUZLDWO0TQSsJPPt0lMyUtX6OW"

	var exists bool
	err := repository.dbPool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", adminUsername).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = repository.dbPool.Exec(ctx, "INSERT INTO users(username, email, password, role) VALUES ($1, $2, $3, $4)", adminUsername, adminEmail, adminPasswordHash, adminRole)
	if err != nil {
		log.Errorf("error while ensuring admin user: %v", err)
		return err
	}

	return nil
}
