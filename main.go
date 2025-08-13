package main

import (
	"api-go-test/config"
	"api-go-test/models"
	"api-go-test/routes"
	"os"

	// "api-go-test/utils"

	"github.com/gin-gonic/gin"
)

func InitializeApp() *gin.Engine {

	r := gin.Default()

	db := config.ConnectDatabase()

	// auto migrate
	db.AutoMigrate(&models.User{}, &models.Product{})
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
