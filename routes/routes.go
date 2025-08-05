package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")
	{
		SetupAuthRoutes(api, db)
		// SetupUserRoutes(api, db)
	}
}
