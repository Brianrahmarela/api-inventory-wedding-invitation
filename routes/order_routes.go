package routes

import (
	"api-go-invitation/controllers"
	"api-go-invitation/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupOrderRoutes(router *gin.RouterGroup, db *gorm.DB) {
	oc := controllers.NewOrderController(db)
	fmt.Println("instance oc", oc)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Create order + upload tamu
		protected.POST("/orders", oc.CreateOrderHandler)
	}
}
