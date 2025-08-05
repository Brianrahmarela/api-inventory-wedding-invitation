package models

import (
	"golang.org/x/crypto/bcrypt"
	// "gorm.io/gorm"
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Name         string    `json:"name" gorm:"size:150;not null"`
	Email        string    `json:"email" gorm:"size:150;uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"size:255;not null"`
	Role         string    `json:"role" gorm:"type:enum('admin','customer');default:'customer'"`
	// Profile      *Profile       `json:"profile,omitempty" gorm:"foreignKey:UserID"`
}

// Hash dan Check tetap pakai PasswordHash
func (u *User) HashPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashed)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

// DTO untuk binding register & login
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"omitempty,oneof=admin customer"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
