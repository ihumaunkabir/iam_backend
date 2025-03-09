package jwork

import (
	"context"
	"time"

	models "iam_backend/models"
	repository "iam_backend/repo"
)

// UserController handles business logic for user operations
type UserController struct {
	userRepo *repository.UserRepository
}

// NewUserController creates a new instance of UserController
func NewUserController(userRepo *repository.UserRepository) *UserController {
	return &UserController{
		userRepo: userRepo,
	}
}

// RegisterUser handles user registration
func (c *UserController) RegisterUser(ctx context.Context, username, email, password string) (*models.User, error) {
	// Create a new user
	user, err := models.NewUser(username, email, password)
	if err != nil {
		return nil, err
	}

	// Save the user to the database
	err = c.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateUser handles user login
func (c *UserController) AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	// Find the user by username
	user, err := c.userRepo.FindByUsernameOrEmail(ctx, username, username)
	if err != nil {
		return nil, err
	}

	// Check password
	if !user.CheckPasswordHash(password) {
		return nil, err
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	err = c.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUserRoles updates roles for a user
func (c *UserController) UpdateUserRoles(ctx context.Context, userID string, roles []string) error {
	user, err := c.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Roles = roles
	return c.userRepo.Update(ctx, user)
}
