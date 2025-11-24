package mocks

import (
	"ingsw3-tp08/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockAuthService es un mock del AuthService para testing
type MockAuthService struct {
	mock.Mock
}

// Register simula el registro de usuario
func (m *MockAuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Login simula el login de usuario
func (m *MockAuthService) Login(creds *models.Credentials) (*models.User, error) {
	args := m.Called(creds)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
