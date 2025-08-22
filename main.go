package main

import (
	"api-go-invitation/config"
	"api-go-invitation/models"
	"api-go-invitation/routes"
	"api-go-invitation/utils"
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func InitializeApp() *gin.Engine {
	r := gin.Default()

	db := config.ConnectDatabase()

	// Auto migrate tabel
	db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.Guest{},
		&models.Payment{},
	)

	routes.SetupRoutes(r, db)

	return r
}

func main() {
	// Flag untuk generate hash password
	// di terminal jlnkan: go run main.go --hash (password yg mau di register ke new admin)
	hashFlag := flag.String("hash", "", "Generate bcrypt hash for a password")
	flag.Parse()

	if *hashFlag != "" {
		utils.PrintHashedPassword(*hashFlag)
		return
	}
	//query db create admin user di dbeaver
	// INSERT INTO users (name, email, password_hash, role, created_at, updated_at)
	// VALUES (
	// 'Admin',
	// 'admin@gmail.com',
	// '$2a$10$JA19QKIBQtpaQxQfXnbvUe6fJh1JBvnq3cfebuuCljgaeI7qvH2Oq', -- hasil dari Go
	// 'admin',
	// NOW(),
	// NOW()
	// );

	// Normal start server
	app := InitializeApp()
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // fallbacks
	}
	fmt.Println("Server running on port", port)
	app.Run(":" + port)
}
