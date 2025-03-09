package jwork

import (
	"context"
	"errors"
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

// GetUserByID retrieves a user by their ID
func (c *UserController) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	// Validate input
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get the user from the repository
	user, err := c.userRepo.FindByID(ctx, userID)
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

// DeactivateUser deactivates a user account
func (c *UserController) DeactivateUser(ctx context.Context, userID string) error {
	user, err := c.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Active = false
	return c.userRepo.Update(ctx, user)
}

// ReactivateUser reactivates a deactivated user account
func (c *UserController) ReactivateUser(ctx context.Context, userID string) error {
	user, err := c.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Active = true
	return c.userRepo.Update(ctx, user)
}

// ChangePassword handles password changes
func (c *UserController) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	// Get the user
	user, err := c.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if !user.CheckPasswordHash(oldPassword) {
		return errors.New("incorrect password")
	}

	// Update with new password
	err = user.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return c.userRepo.Update(ctx, user)
}
