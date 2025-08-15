package models

import "time"

type Payment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderID    uint      `gorm:"index;not null" json:"order_id"`
	Amount     int64     `json:"amount"`                                   // jumlah yang dikirim (optional dari user)
	ProofURL   string    `gorm:"size:512" json:"proof_url"`                // URL file bukti di Cloudinary
	Status     string    `gorm:"size:50;default:'uploaded'" json:"status"` // uploaded, verified, rejected
	VerifiedBy uint      `json:"verified_by"`                              // admin user id (jika ada)
	VerifiedAt time.Time `json:"verified_at"`
	CreatedAt  time.Time `json:"created_at"`
}
