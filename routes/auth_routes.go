package routes

import (
	"api-go-test/controllers"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func SetupAuthRoutes(router *gin.RouterGroup, db *gorm.DB) {
	// Buat object bernama authController
	// NewAuthController: constructor (pembuat objek)
	// Parameternya db: dikirim agar bisa disimpan di field DB struct AuthService
	authController := controllers.NewAuthController(db)
	// Membuat group route
	// Bisa tambahkan middleware di sini, misalnya autentikasi
	protected := router.Group("/")
	{
		//Saat user POST ke /api/register Jalankan method Register() di controller authController
		protected.POST("/register", authController.Register)
		protected.POST("/login", authController.Login)
	}
}
