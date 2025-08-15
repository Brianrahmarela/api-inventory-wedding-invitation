package services

import (
	"api-go-invitation/models"
	"api-go-invitation/utils"

	// "encoding/json"
	"errors"
	// "fmt"
	"net/http"

	"strings"

	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (as *AuthService) Register(req *models.RegisterRequest) (string, int, error) {
	// cust role not allowed register to admin
	if strings.ToLower(req.Role) == "admin" {
		return "", http.StatusForbidden, errors.New("you cannot register as admin")
	}

	// check duplicate email
	var exist models.User
	if err := as.DB.Where("email = ?", req.Email).First(&exist).Error; err == nil {
		return "", http.StatusConflict, errors.New("email has been registered")
	}

	// prepare user
	user := models.User{
		Name:  req.Name,
		Email: req.Email,
	}
	// set role atau default ke customer
	if req.Role == "admin" {
		user.Role = "admin"
	} else {
		user.Role = "customer"
	}

	// hash password
	if err := user.HashPassword(req.Password); err != nil {
		return "", http.StatusInternalServerError, errors.New("error hashing password")
	}

	// simpan
	if err := as.DB.Create(&user).Error; err != nil {
		return "", http.StatusInternalServerError, errors.New("error creating user")
	}

	// debug print
	// u, _ := json.MarshalIndent(user, "", "  ")
	// fmt.Println("Registered user:", string(u))

	// generate token include role
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("error generating token")
	}

	return token, http.StatusCreated, nil
}

func (as *AuthService) Login(req *models.LoginRequest) (string, string, error) {
	var user models.User
	if err := as.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", errors.New("invalid email or password")
		}
		return "", "", err
	}
	if err := user.CheckPassword(req.Password); err != nil {
		return "", "", errors.New("invalid email or password")
	}

	// generate token with role
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", "", errors.New("error generating token")
	}
	return token, user.Role, nil
}
