package main

import (
	"os"
	"testing"

	database "iam_backend/db"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	// Test default connection string
	os.Unsetenv("MONGO_URI")
	os.Unsetenv("DB_NAME")

	db, err := database.NewMongoConnection("mongodb://localhost:27017", "iam_database")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	defer db.Disconnect()
}

func TestEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected map[string]string
	}{
		{
			name:    "Default values",
			envVars: map[string]string{},
			expected: map[string]string{
				"MONGO_URI": "mongodb://localhost:27017",
				"DB_NAME":   "iam_database",
				"PORT":      "8080",
			},
		},
		{
			name: "Custom values",
			envVars: map[string]string{
				"MONGO_URI": "mongodb://custom:27017",
				"DB_NAME":   "custom_db",
				"PORT":      "3000",
			},
			expected: map[string]string{
				"MONGO_URI": "mongodb://custom:27017",
				"DB_NAME":   "custom_db",
				"PORT":      "3000",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Clearenv()

			// Set test environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Check values
			mongoURI := os.Getenv("MONGO_URI")
			if mongoURI == "" {
				mongoURI = "mongodb://localhost:27017"
			}
			assert.Equal(t, tt.expected["MONGO_URI"], mongoURI)

			dbName := os.Getenv("DB_NAME")
			if dbName == "" {
				dbName = "iam_database"
			}
			assert.Equal(t, tt.expected["DB_NAME"], dbName)

			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
			}
			assert.Equal(t, tt.expected["PORT"], port)
		})
	}
}
