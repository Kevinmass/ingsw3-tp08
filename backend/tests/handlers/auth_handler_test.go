package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"ingsw3-tp08/internal/handlers"
	"ingsw3-tp08/internal/models"
	"ingsw3-tp08/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test helper function to create request body
func createRequestBody(t *testing.T, data interface{}) *bytes.Buffer {
	t.Helper()
	body, err := json.Marshal(data)
	assert.NoError(t, err)
	return bytes.NewBuffer(body)
}

// Test Register endpoint - Success case
func TestAuthHandler_Register_Success(t *testing.T) {
	// ARRANGE
	mockAuthService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockAuthService)

	registerRequest := models.RegisterRequest{
		Email:    "test@example.com",
		Password: "123456",
		Username: "testuser",
	}

	expectedUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
	}

	mockAuthService.On("Register", mock.AnythingOfType("*models.RegisterRequest")).Return(expectedUser, nil)

	// ACT
	req := httptest.NewRequest("POST", "/api/auth/register", createRequestBody(t, registerRequest))
	rr := httptest.NewRecorder()
	handler.Register(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response models.User
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Email, response.Email)
	assert.Equal(t, expectedUser.Username, response.Username)

	mockAuthService.AssertExpectations(t)
}

// Test Register endpoint - JSON parsing error
func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	// ARRANGE
	mockAuthService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockAuthService)

	// ACT - Send invalid JSON
	invalidJSON := `{"email": "test@example.com", "invalid": }`
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(invalidJSON))
	rr := httptest.NewRecorder()
	handler.Register(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "JSON inválido", response["error"])
}

// Test Register endpoint - Service returns error
func TestAuthHandler_Register_ServiceError(t *testing.T) {
	// ARRANGE
	mockAuthService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockAuthService)

	registerRequest := models.RegisterRequest{
		Email:    "existing@example.com",
		Password: "123456",
		Username: "testuser",
	}

	serviceError := errors.New("el email ya está registrado")
	mockAuthService.On("Register", mock.AnythingOfType("*models.RegisterRequest")).Return(nil, serviceError)

	// ACT
	req := httptest.NewRequest("POST", "/api/auth/register", createRequestBody(t, registerRequest))
	rr := httptest.NewRecorder()
	handler.Register(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, serviceError.Error(), response["error"])

	mockAuthService.AssertExpectations(t)
}

// Test Login endpoint - Success case
func TestAuthHandler_Login_Success(t *testing.T) {
	// ARRANGE
	mockAuthService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockAuthService)

	credentials := models.Credentials{
		Email:    "test@example.com",
		Password: "123456",
	}

	expectedUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
	}

	mockAuthService.On("Login", mock.AnythingOfType("*models.Credentials")).Return(expectedUser, nil)

	// ACT
	req := httptest.NewRequest("POST", "/api/auth/login", createRequestBody(t, credentials))
	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.User
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Email, response.Email)
	assert.Equal(t, expectedUser.Username, response.Username)

	mockAuthService.AssertExpectations(t)
}

// Test Login endpoint - JSON parsing error
func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	// ARRANGE
	mockAuthService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockAuthService)

	// ACT - Send invalid JSON
	invalidJSON := `{"email": "test@example.com", "invalid": }`
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(invalidJSON))
	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "JSON inválido", response["error"])
}

// Test Login endpoint - Service returns error
func TestAuthHandler_Login_ServiceError(t *testing.T) {
	// ARRANGE
	mockAuthService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockAuthService)

	credentials := models.Credentials{
		Email:    "wrong@example.com",
		Password: "wrongpassword",
	}

	serviceError := errors.New("credenciales inválidas")
	mockAuthService.On("Login", mock.AnythingOfType("*models.Credentials")).Return(nil, serviceError)

	// ACT
	req := httptest.NewRequest("POST", "/api/auth/login", createRequestBody(t, credentials))
	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, serviceError.Error(), response["error"])

	mockAuthService.AssertExpectations(t)
}
