package models

import "time"

type Order struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	UserID      uint   `gorm:"not null" json:"user_id"`
	ProductID   uint   `gorm:"not null" json:"product_id"`
	TotalAmount int64  `json:"total_amount"`
	Status      string `gorm:"size:50;default:'pending'" json:"status"`
	GroomName   string `gorm:"size:100;not null" json:"groom_name"` // nama mempelai pria
	BrideName   string `gorm:"size:100;not null" json:"bride_name"` // nama mempelai wanita

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Guests  []Guest `gorm:"foreignKey:OrderID" json:"guests"`
	Payment Payment `gorm:"foreignKey:OrderID" json:"payment"`
}
