package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ingsw3-tp08/internal/models"
	"ingsw3-tp08/tests/mocks"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostHandler_CreatePost_Success(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	postHandler := NewPostHandler(mockPostService)

	req := models.CreatePostRequest{
		Title:   "Test Post",
		Content: "Test Content",
	}

	expectedPost := &models.Post{
		ID:       1,
		Title:    "Test Post",
		Content:  "Test Content",
		UserID:   1,
		Username: "testuser",
	}

	mockPostService.On("CreatePost", &req, 1).Return(expectedPost, nil)

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()

	// ACT
	postHandler.CreatePost(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusCreated, w.Code)
	mockPostService.AssertExpectations(t)

	var response models.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost.ID, response.ID)
	assert.Equal(t, expectedPost.Title, response.Title)
	assert.Equal(t, expectedPost.Content, response.Content)
	assert.Equal(t, expectedPost.UserID, response.UserID)
	assert.Equal(t, expectedPost.Username, response.Username)
}

func TestPostHandler_CreatePost_InvalidJSON(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	postHandler := NewPostHandler(mockPostService)

	httpReq := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBufferString("invalid json"))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// ACT
	postHandler.CreatePost(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidJSON, response["error"])

	mockPostService.AssertNotCalled(t, "CreatePost", mock.Anything, mock.Anything)
}

func TestPostHandler_CreatePost_MissingUserID(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	postHandler := NewPostHandler(mockPostService)

	req := models.CreatePostRequest{
		Title:   "Test Post",
		Content: "Test Content",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	// No X-User-ID

	w := httptest.NewRecorder()

	// ACT
	postHandler.CreatePost(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrUserNotAuthenticated, response["error"])

	mockPostService.AssertNotCalled(t, "CreatePost", mock.Anything, mock.Anything)
}

func TestPostHandler_CreatePost_InvalidUserID(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	postHandler := NewPostHandler(mockPostService)

	req := models.CreatePostRequest{
		Title:   "Test Post",
		Content: "Test Content",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/posts", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", "abc") // Invalid int

	w := httptest.NewRecorder()

	// ACT
	postHandler.CreatePost(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidUserID, response["error"])

	mockPostService.AssertNotCalled(t, "CreatePost", mock.Anything, mock.Anything)
}

func TestPostHandler_GetAllPosts_Success(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	postHandler := NewPostHandler(mockPostService)

	expectedPosts := []*models.Post{
		{ID: 1, Title: "Post 1", Content: "Content 1"},
		{ID: 2, Title: "Post 2", Content: "Content 2"},
	}

	mockPostService.On("GetAllPosts").Return(expectedPosts, nil)

	httpReq := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
	w := httptest.NewRecorder()

	// ACT
	postHandler.GetAllPosts(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusOK, w.Code)
	mockPostService.AssertExpectations(t)

	var response []*models.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Post 1", response[0].Title)
}

func TestPostHandler_GetAllPosts_ServiceError(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	postHandler := NewPostHandler(mockPostService)

	mockPostService.On("GetAllPosts").Return(nil, assert.AnError)

	httpReq := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
	w := httptest.NewRecorder()

	// ACT
	postHandler.GetAllPosts(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPostService.AssertExpectations(t)
}

func TestPostHandler_GetPostByID_Success(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	postHandler := NewPostHandler(mockPostService)

	expectedPost := &models.Post{
		ID:      1,
		Title:   "Test Post",
		Content: "Test Content",
	}

	mockPostService.On("GetPostByID", 1).Return(expectedPost, nil)

	httpReq := httptest.NewRequest(http.MethodGet, "/api/posts/1", nil)
	httpReq = mux.SetURLVars(httpReq, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	// ACT
	postHandler.GetPostByID(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusOK, w.Code)
	mockPostService.AssertExpectations(t)

	var response models.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost.ID, response.ID)
}

func TestPostHandler_GetPostByID_InvalidID(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	postHandler := NewPostHandler(mockPostService)

	httpReq := httptest.NewRequest(http.MethodGet, "/api/posts/abc", nil) // Invalid ID
	w := httptest.NewRecorder()

	// ACT
	postHandler.GetPostByID(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidID, response["error"])

	mockPostService.AssertNotCalled(t, "GetPostByID", mock.Anything)
}

func TestPostHandler_DeletePost_Success(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	postHandler := NewPostHandler(mockPostService)

	mockPostService.On("DeletePost", 1, 1).Return(nil)

	httpReq := httptest.NewRequest(http.MethodDelete, "/api/posts/1", nil)
	httpReq.Header.Set("X-User-ID", "1")
	httpReq = mux.SetURLVars(httpReq, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	// ACT
	postHandler.DeletePost(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusOK, w.Code)
	mockPostService.AssertExpectations(t)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Post eliminado", response["message"])
}

func TestPostHandler_DeletePost_PermissionDenied(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	postHandler := NewPostHandler(mockPostService)

	mockPostService.On("DeletePost", 1, 2).Return(assert.AnError) // Different user

	httpReq := httptest.NewRequest(http.MethodDelete, "/api/posts/1", nil)
	httpReq.Header.Set("X-User-ID", "2")
	httpReq = mux.SetURLVars(httpReq, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	// ACT
	postHandler.DeletePost(w, httpReq)

	// ASSERT
	assert.Equal(t, http.StatusForbidden, w.Code)
	mockPostService.AssertExpectations(t)
}
