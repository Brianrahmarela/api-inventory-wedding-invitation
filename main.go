package main

import (
	"api-go-test/config"
	"api-go-test/models"
	"api-go-test/routes"
	"os"

	// "api-go-test/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func InitializeApp() *gin.Engine {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading ENV")
	}

	r := gin.Default()

	db := config.ConnectDatabase()

	// auto migrate
	db.AutoMigrate(&models.User{})
	// db.AutoMigrate(&models.User{}, &models.Profile{})

	routes.SetupRoutes(r, db)

	return r
}

func main() {
	//create hash password to make admin manually from mysql
	// utils.PrintHashedPassword("123456")
	app := InitializeApp()
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // fallback
	}
	app.Run(":" + port)
}
