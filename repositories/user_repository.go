package repositories

import (
	"api-go-invitation/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
} //interface dibuat agar siapa pun yang implementasi UserRepository wajib punya fungsi Create, FindByEmail, FindByID.

type userRepository struct {
	db *gorm.DB //Struct yg nyimpen koneksi database GORM.
}

// Constructor function → cara bikin repository baru.
// Input: db (koneksi database), Output: UserRepository (interface).
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// Kenapa return UserRepository bukan *userRepository?
// → Supaya yang pakai cuma tahu “interface-nya”, bukan detail implementasinya.

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
