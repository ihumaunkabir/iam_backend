package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserController is a mock implementation of the UserController
type MockUserController struct {
	mock.Mock
}

func (m *MockUserController) Register(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserController) Login(email, password string) (*models.User, error) {
	args := m.Called(email, password)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestSetupRouter(t *testing.T) {
	mockController := new(MockUserController)
	router := SetupRouter(mockController)

	assert.NotNil(t, router, "Router should not be nil")

	// Test register endpoint
	t.Run("Register Endpoint", func(t *testing.T) {
		user := models.User{
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonValue, _ := json.Marshal(user)

		mockController.On("Register", mock.AnythingOfType("*models.User")).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test login endpoint
	t.Run("Login Endpoint", func(t *testing.T) {
		credentials := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonValue, _ := json.Marshal(credentials)

		mockUser := &models.User{
			Email: "test@example.com",
		}
		mockController.On("Login", "test@example.com", "password123").Return(mockUser, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test non-existent endpoint
	t.Run("Non-existent Endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/nonexistent", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
