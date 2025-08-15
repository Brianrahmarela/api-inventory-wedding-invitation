package services

import (
	"api-go-invitation/models"
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type PaymentService struct {
	DB *gorm.DB
}

func NewPaymentService(db *gorm.DB) *PaymentService {
	// Load .env ketika service dibuat
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on system env")
	}
	return &PaymentService{DB: db}
}

// UploadProof mengunggah bukti pembayaran ke Cloudinary
func (ps *PaymentService) UploadProof(orderID uint, file *multipart.FileHeader) (models.Payment, error) {
	// Cek order
	var order models.Order
	if err := ps.DB.First(&order, orderID).Error; err != nil {
		return models.Payment{}, fmt.Errorf("order not found")
	}

	// Ambil env Cloudinary
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return models.Payment{}, errors.New("cloudinary config is missing in .env")
	}

	// Init Cloudinary
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return models.Payment{}, fmt.Errorf("cloudinary init failed: %v", err)
	}

	// Buka file sebelum upload
	f, err := file.Open()
	if err != nil {
		return models.Payment{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	// Upload ke Cloudinary
	uploadResp, err := cld.Upload.Upload(
		context.Background(),
		f,
		uploader.UploadParams{
			Folder:   "payments",
			PublicID: fmt.Sprintf("order_%d", orderID),
		},
	)
	if err != nil {
		return models.Payment{}, fmt.Errorf("upload to cloudinary failed: %v", err)
	}

	// Simpan ke database
	payment := models.Payment{
		OrderID:  orderID,
		Amount:   0,
		ProofURL: uploadResp.SecureURL,
		Status:   "uploaded",
	}

	if err := ps.DB.Create(&payment).Error; err != nil {
		return models.Payment{}, err
	}

	return payment, nil
}

// VerifyPayment memverifikasi pembayaran oleh admin
func (ps *PaymentService) VerifyPayment(orderID uint, amount float64, adminID uint) error {
	var order models.Order
	if err := ps.DB.First(&order, orderID).Error; err != nil {
		return errors.New("invalid order_id")
	}

	if float64(order.TotalAmount) != amount {
		return errors.New("invalid amount")
	}

	if err := ps.DB.Model(&models.Payment{}).
		Where("order_id = ?", orderID).
		Updates(map[string]interface{}{
			"amount":      amount,
			"status":      "verified",
			"verified_by": adminID,
			"verified_at": time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("update payment failed: %v", err)
	}

	return ps.DB.Model(&order).Update("status", "paid").Error
}
