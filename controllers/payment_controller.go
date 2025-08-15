package controllers

import (
	"api-go-invitation/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	Service *services.PaymentService
}

func NewPaymentController(service *services.PaymentService) *PaymentController {
	return &PaymentController{Service: service}
}

// Upload bukti pembayaran (Cloudinary upload)
func (pc *PaymentController) UploadProof(c *gin.Context) {
	orderIDStr := c.PostForm("order_id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	payment, err := pc.Service.UploadProof(uint(orderID), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// Verifikasi pembayaran oleh admin
func (pc *PaymentController) VerifyPayment(c *gin.Context) {
	var req struct {
		OrderID string `json:"order_id"`
		Amount  string `json:"amount"`
		AdminID string `json:"admin_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID, err := strconv.ParseUint(req.OrderID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order_id"})
		return
	}

	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	}

	adminID, err := strconv.ParseUint(req.AdminID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin_id"})
		return
	}

	if err := pc.Service.VerifyPayment(uint(orderID), amount, uint(adminID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment verified successfully"})
}
