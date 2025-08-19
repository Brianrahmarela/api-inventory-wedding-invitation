package controllers

import (
	"api-go-invitation/models"
	"api-go-invitation/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NotificationController struct {
	DB *gorm.DB
}

func NewNotificationController(db *gorm.DB) *NotificationController {
	return &NotificationController{DB: db}
}

// Kirim email link undangan setelah pembayaran verified
func (nc *NotificationController) SendGuestLinks(c *gin.Context) {
	userID := c.Param("userId") // misal: /notify/2

	// Cek order yang statusnya paid dan payment verified
	var order models.Order
	if err := nc.DB.Where("user_id = ? AND status = ?", userID, "paid").First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order not found"})
		return
	}

	var payment models.Payment
	if err := nc.DB.Where("order_id = ? AND status = ?", order.ID, "verified").First(&payment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment not verified"})
		return
	}

	// Ambil semua guest link
	var guests []models.Guest
	if err := nc.DB.Where("order_id = ?", order.ID).Find(&guests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch guests"})
		return
	}

	links := []string{}
	for _, g := range guests {
		links = append(links, g.Link)
	}

	// Ambil email user pemesan
	var user models.User
	if err := nc.DB.First(&user, order.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	// Format email
	subject := "Link Undangan Pernikahan Anda"
	body := utils.FormatGuestLinks(links)

	// Kirim email
	if err := utils.SendEmail(user.Email, subject, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send email", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email sent successfully"})
}
