package main

import (
	"log"
	"os"

	database "iam_backend/db"
	controllers "iam_backend/jwork"
	repository "iam_backend/repo"
	"iam_backend/router"
)

func main() {
	// Load environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "iam_database"
	}

	// Establish database connection
	db, err := database.NewMongoConnection(mongoURI, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Disconnect()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize controllers
	userController := controllers.NewUserController(userRepo)

	// Setup router
	r := router.SetupRouter(userController)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on :%s", port)
	log.Fatal(r.Run(":" + port))
}
