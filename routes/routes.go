package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Grup /api
	api := router.Group("/api")
	{
		// Endpoint test: /api/about
		api.GET("/about", func(c *gin.Context) {
			c.String(http.StatusOK, "hello!")
		})

		SetupAuthRoutes(api, db)
		SetupProductRoutes(api, db)
		SetupOrderRoutes(api, db)
		SetupPaymentRoutes(api, db)
	}
}
