package services

import (
	"ingsw3-tp08/internal/models"
)

// AuthServiceInterface defines the contract for authentication services
type AuthServiceInterface interface {
	Register(req *models.RegisterRequest) (*models.User, error)
	Login(creds *models.Credentials) (*models.User, error)
}

// Ensure AuthService implements the interface
var _ AuthServiceInterface = (*AuthService)(nil)
