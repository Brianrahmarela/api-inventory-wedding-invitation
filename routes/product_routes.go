package routes

import (
	"api-go-invitation/controllers"
	"api-go-invitation/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupProductRoutes(router *gin.RouterGroup, db *gorm.DB) {
	//bikin object ProductController baru. istilah di oop adalah "instance"(hasil nyata dari cetakan struct)
	pc := controllers.NewProductController(db)
	//hasil pc hanya  alamat pointer saja -> pc = *ProductController (0x1400EEFF),
	// Saat request datang, handler memakai pc.ProductService untuk akses DB melalui ProductService
	fmt.Println("instance pc", pc)

	// Public
	router.GET("/products", pc.GetAll)
	router.GET("/products/:id", pc.GetByID)
	router.GET("/products/slug/:slug", pc.GetBySlug)

	// Protected (admin)
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/products", pc.Create)
		protected.PUT("/products/:id", pc.Update)
		protected.DELETE("/products/:id", pc.Delete)
	}
}
