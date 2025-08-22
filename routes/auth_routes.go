package routes

import (
	"api-go-invitation/controllers"
	"api-go-invitation/repositories"
	"api-go-invitation/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// *db gorm.DB → menerima koneksi database, agar bisa dipakai utk bikin repository.
func SetupAuthRoutes(router *gin.RouterGroup, db *gorm.DB) {
	// Buat repository
	userRepo := repositories.NewUserRepository(db)
	//repositories.NewUserRepository(db) → bikin lapisan repository yang khusus handle query database untuk user.
	//userRepo skrg bisa dipakai untuk CRUD user, tanpa harus tulis query manual di service/controller.

	// Buat service
	authService := services.NewAuthService(userRepo)
	//Service ini bergantung ke repository (userRepo) untuk ambil/simpan data ke DB.

	// Buat controller
	authController := controllers.NewAuthController(authService)
	//controllers.NewAuthController(authService) → bikin controller yang handle request/response API.

	// Routes
	protected := router.Group("/")
	{
		protected.POST("/register", authController.Register)
		protected.POST("/login", authController.Login)
	}
}
