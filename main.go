package main

import (
	"api-go-invitation/config"
	"api-go-invitation/models"
	"api-go-invitation/routes"
	"os"

	// "api-go-invitation/utils"

	"github.com/gin-gonic/gin"
)

func InitializeApp() *gin.Engine {

	r := gin.Default()

	db := config.ConnectDatabase()

	// auto migrate
	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.Guest{}, &models.Payment{})

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
