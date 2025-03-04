package dbq

import (
	"context"
	"errors"
	"time"

	"iam_backend/db"
	users "iam_backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserDBQ handles database operations for users
type UserDBQ struct {
	collection *mongo.Collection
}

// NewUserDBQ creates a new instance of UserDBQ
func NewUserDBQ(db *db.Database) *UserDBQ {
	return &UserDBQ{
		collection: db.Database.Collection("users"),
	}
}

// Create inserts a new user into the database
func (r *UserDBQ) Create(ctx context.Context, user *users.User) error {
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
func (r *UserDBQ) FindByID(ctx context.Context, id string) (*users.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user users.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByUsernameOrEmail finds a user by username or email
func (r *UserDBQ) FindByUsernameOrEmail(ctx context.Context, username, email string) (*users.User, error) {
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
func (r *UserDBQ) Update(ctx context.Context, user *users.User) error {
	user.UpdatedAt = time.Now()
	update := bson.M{"$set": user}

	_, err := r.collection.UpdateByID(ctx, user.ID, update)
	return err
}

// Delete removes a user from the database
func (r *UserDBQ) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// UpdateLastLogin updates the last login time for a user
func (r *UserDBQ) UpdateLastLogin(ctx context.Context, userID primitive.ObjectID) error {
	now := time.Now()
	update := bson.M{"$set": bson.M{"last_login": now}}

	_, err := r.collection.UpdateByID(ctx, userID, update)
	return err
}
