package routes

import (
	"api-go-invitation/controllers"
	"api-go-invitation/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupNotificationRoutes(router *gin.RouterGroup, db *gorm.DB) {
	nc := controllers.NewNotificationController(db)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Kirim email link undangan (kalau status payment verified & status order paid)
		protected.GET("/notify/:userId", nc.SendGuestLinks)
	}
}
