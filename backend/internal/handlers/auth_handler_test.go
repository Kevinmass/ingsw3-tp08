package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ingsw3-tp08/internal/models"
	"ingsw3-tp08/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Register_Success(t *testing.T) {
	// ARRANGE
	mockAuthService := new(mocks.MockAuthService)
	authHandler := NewAuthHandler(mockAuthService)

	req := models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Username: "testuser",
	}

	expectedUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
	}

	mockAuthService.On("Register", &req).Return(expectedUser, nil)

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// ACT
	authHandler.Register(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusCreated, w.Code)
	mockAuthService.AssertExpectations(t)

	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.Email, response.Email)
	assert.Equal(t, expectedUser.Username, response.Username)
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	// ARRANGE
	mockAuthService := new(mocks.MockAuthService)
	authHandler := NewAuthHandler(mockAuthService)

	httpReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString("invalid json"))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// ACT
	authHandler.Register(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	assert.Equal(t, "JSON inválido", response["error"])

	mockAuthService.AssertNotCalled(t, "Register", mock.Anything)
}

func TestAuthHandler_Register_ServiceError(t *testing.T) {
	// ARRANGE
	mockAuthService := new(mocks.MockAuthService)
	authHandler := NewAuthHandler(mockAuthService)

	req := models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Username: "testuser",
	}

	mockAuthService.On("Register", &req).Return(nil, assert.AnError)

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// ACT
	authHandler.Register(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, assert.AnError.Error(), response["error"])

	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	// ARRANGE
	mockAuthService := new(mocks.MockAuthService)
	authHandler := NewAuthHandler(mockAuthService)

	creds := models.Credentials{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
	}

	mockAuthService.On("Login", &creds).Return(expectedUser, nil)

	body, _ := json.Marshal(creds)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// ACT
	authHandler.Login(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusOK, w.Code)
	mockAuthService.AssertExpectations(t)

	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.Email, response.Email)
	assert.Equal(t, expectedUser.Username, response.Username)
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	// ARRANGE
	mockAuthService := new(mocks.MockAuthService)
	authHandler := NewAuthHandler(mockAuthService)

	httpReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString("invalid json"))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// ACT
	authHandler.Login(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	assert.Equal(t, "JSON inválido", response["error"])

	mockAuthService.AssertNotCalled(t, "Login", mock.Anything)
}

func TestAuthHandler_Login_ServiceError(t *testing.T) {
	// ARRANGE
	mockAuthService := new(mocks.MockAuthService)
	authHandler := NewAuthHandler(mockAuthService)

	creds := models.Credentials{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockAuthService.On("Login", &creds).Return(nil, assert.AnError)

	body, _ := json.Marshal(creds)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// ACT
	authHandler.Login(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	assert.Equal(t, assert.AnError.Error(), response["error"])

	mockAuthService.AssertExpectations(t)
}
