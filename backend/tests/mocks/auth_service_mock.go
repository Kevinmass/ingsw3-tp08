package mocks

import (
	"ingsw3-tp08/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

// Register mocks the Register method
func (m *MockAuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	args := m.Called(req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.User), args.Error(1)
}

// Login mocks the Login method
func (m *MockAuthService) Login(creds *models.Credentials) (*models.User, error) {
	args := m.Called(creds)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.User), args.Error(1)
}
