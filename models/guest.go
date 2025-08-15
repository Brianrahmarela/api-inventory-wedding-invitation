package models

import "time"

type Guest struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `gorm:"index;not null" json:"order_id"`
	Name      string    `gorm:"size:255" json:"name"`
	Partner   string    `gorm:"size:255" json:"partner"`
	Email     string    `gorm:"size:255" json:"email"`
	Phone     string    `gorm:"size:50" json:"phone"`
	Link      string    `gorm:"size:512" json:"link"` // link undangan unik
	CreatedAt time.Time `json:"created_at"`
}
