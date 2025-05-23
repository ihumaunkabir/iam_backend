package repository

import (
	"context"
	"errors"
	"time"

	database "iam_backend/db"
	models "iam_backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository handles database operations for users
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *database.Database) *UserRepository {
	return &UserRepository{
		collection: db.Database.Collection("users"),
	}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	// Check if username or email already exists
	existingUser, _ := r.FindByUsernameOrEmail(ctx, user.Username, user.Email)
	if existingUser != nil {
		return errors.New("username or email already exists")
	}

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID retrieves a user by their ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByUsernameOrEmail finds a user by username or email
func (r *UserRepository) FindByUsernameOrEmail(ctx context.Context, username, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": username},
			{"email": email},
		},
	}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	update := bson.M{"$set": user}

	_, err := r.collection.UpdateByID(ctx, user.ID, update)
	return err
}

// Delete removes a user from the database
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// UpdateLastLogin updates the last login time for a user
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID primitive.ObjectID) error {
	now := time.Now()
	update := bson.M{"$set": bson.M{"last_login": now}}

	_, err := r.collection.UpdateByID(ctx, userID, update)
	return err
}
