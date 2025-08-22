package services

import (
	"api-go-invitation/models"
	"api-go-invitation/repositories"
	"errors"
)

type AuthService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{userRepo}
}

func (s *AuthService) Register(req *models.RegisterRequest) error {
	// mapping DTO ke model
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	// hash password → langsung simpan ke field PasswordHash
	if err := user.HashPassword(req.Password); err != nil {
		return err
	}

	return s.userRepo.Create(user)
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// check password
	if err := user.CheckPassword(req.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
