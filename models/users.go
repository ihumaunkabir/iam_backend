package users

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User represents the user model
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username" validate:"required,min=3,max=50"`
	Email        string             `bson:"email" json:"email" validate:"required,email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Roles        []string           `bson:"roles" json:"roles"`
	Active       bool               `bson:"active" json:"active"`
	LastLogin    *time.Time         `bson:"last_login" json:"last_login"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

// HashPassword generates a bcrypt hash of the password
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.PasswordHash = string(bytes)
	return nil
}

// CheckPasswordHash checks if the provided password matches the stored hash
func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// NewUser creates a new user with default values
func NewUser(username, email, password string) (*User, error) {
	now := time.Now()
	user := &User{
		Username:  username,
		Email:     email,
		Active:    true,
		Roles:     []string{"user"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := user.HashPassword(password); err != nil {
		return nil, err
	}

	return user, nil
}
