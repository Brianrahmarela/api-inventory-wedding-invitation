package routes

import (
	"api-go-invitation/controllers"
	"api-go-invitation/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupPaymentRoutes(r *gin.RouterGroup, db *gorm.DB) {
	paymentService := services.NewPaymentService(db)
	paymentController := controllers.NewPaymentController(paymentService)

	// Customer
	r.POST("/payments/upload", paymentController.UploadProof)

	// Admin
	r.POST("/payments/verify", paymentController.VerifyPayment)
}
