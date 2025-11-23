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

// Test CreatePost - Success case
func TestPostHandler_CreatePost_Success(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	handler := handlers.NewPostHandler(mockPostService)

	createPostRequest := models.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post content",
	}

	expectedPost := &models.Post{
		ID:      1,
		Title:   "Test Post",
		Content: "This is a test post content",
	}

	mockPostService.On("CreatePost", mock.AnythingOfType("*models.CreatePostRequest"), 1).Return(expectedPost, nil)

	// ACT
	body, _ := json.Marshal(createPostRequest)
	req := httptest.NewRequest("POST", "/api/posts", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "1")
	rr := httptest.NewRecorder()
	handler.CreatePost(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response models.Post
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost.Title, response.Title)
	assert.Equal(t, expectedPost.Content, response.Content)

	mockPostService.AssertExpectations(t)
}

// Test CreatePost - Invalid JSON
func TestPostHandler_CreatePost_InvalidJSON(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	handler := handlers.NewPostHandler(mockPostService)

	// ACT - Send invalid JSON
	invalidJSON := `{"title": "Test Post", "invalid": }`
	req := httptest.NewRequest("POST", "/api/posts", bytes.NewBufferString(invalidJSON))
	req.Header.Set("X-User-ID", "1")
	rr := httptest.NewRecorder()
	handler.CreatePost(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "JSON inválido", response["error"])
}

// Test CreatePost - Service returns error
func TestPostHandler_CreatePost_ServiceError(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	handler := handlers.NewPostHandler(mockPostService)

	createPostRequest := models.CreatePostRequest{
		Title:   "Test Post",
		Content: "This is a test post content",
	}

	serviceError := errors.New("error creating post")
	mockPostService.On("CreatePost", mock.AnythingOfType("*models.CreatePostRequest"), 1).Return(nil, serviceError)

	// ACT
	body, _ := json.Marshal(createPostRequest)
	req := httptest.NewRequest("POST", "/api/posts", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "1")
	rr := httptest.NewRecorder()
	handler.CreatePost(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, serviceError.Error(), response["error"])

	mockPostService.AssertExpectations(t)
}

// Test GetAllPosts - Success case
func TestPostHandler_GetAllPosts_Success(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	handler := handlers.NewPostHandler(mockPostService)

	expectedPosts := []*models.Post{
		{
			ID:      1,
			Title:   "First Post",
			Content: "First post content",
		},
		{
			ID:      2,
			Title:   "Second Post",
			Content: "Second post content",
		},
	}

	mockPostService.On("GetAllPosts").Return(expectedPosts, nil)

	// ACT
	req := httptest.NewRequest("GET", "/api/posts", nil)
	rr := httptest.NewRecorder()
	handler.GetAllPosts(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusOK, rr.Code)

	var response []*models.Post
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response))
	assert.Equal(t, expectedPosts[0].Title, response[0].Title)
	assert.Equal(t, expectedPosts[1].Title, response[1].Title)

	mockPostService.AssertExpectations(t)
}

// Test GetAllPosts - Service returns error
func TestPostHandler_GetAllPosts_ServiceError(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	handler := handlers.NewPostHandler(mockPostService)

	serviceError := errors.New("database error")
	mockPostService.On("GetAllPosts").Return(nil, serviceError)

	// ACT
	req := httptest.NewRequest("GET", "/api/posts", nil)
	rr := httptest.NewRecorder()
	handler.GetAllPosts(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, serviceError.Error(), response["error"])

	mockPostService.AssertExpectations(t)
}

// Test GetPostByID - Success case
func TestPostHandler_GetPostByID_Success(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	handler := handlers.NewPostHandler(mockPostService)

	expectedPost := &models.Post{
		ID:      1,
		Title:   "Test Post",
		Content: "Test post content",
	}

	mockPostService.On("GetPostByID", 1).Return(expectedPost, nil)

	// ACT
	req := httptest.NewRequest("GET", "/api/posts/1", nil)
	rr := httptest.NewRecorder()
	handler.GetPostByID(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.Post
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost.ID, response.ID)
	assert.Equal(t, expectedPost.Title, response.Title)

	mockPostService.AssertExpectations(t)
}

// Test GetPostByID - Post not found
func TestPostHandler_GetPostByID_NotFound(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	handler := handlers.NewPostHandler(mockPostService)

	mockPostService.On("GetPostByID", 999).Return(nil, nil)

	// ACT
	req := httptest.NewRequest("GET", "/api/posts/999", nil)
	rr := httptest.NewRecorder()
	handler.GetPostByID(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusNotFound, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "post no encontrado", response["error"])

	mockPostService.AssertExpectations(t)
}

// Test GetPostByID - Invalid ID
func TestPostHandler_GetPostByID_InvalidID(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	handler := handlers.NewPostHandler(mockPostService)

	// ACT
	req := httptest.NewRequest("GET", "/api/posts/invalid", nil)
	rr := httptest.NewRecorder()
	handler.GetPostByID(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ID inválido", response["error"])
}

// Test DeletePost - Success case
func TestPostHandler_DeletePost_Success(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	handler := handlers.NewPostHandler(mockPostService)

	mockPostService.On("DeletePost", 1, 1).Return(nil)

	// ACT
	req := httptest.NewRequest("DELETE", "/api/posts/1", nil)
	req.Header.Set("X-User-ID", "1")
	rr := httptest.NewRecorder()
	handler.DeletePost(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Post eliminado", response["message"])

	mockPostService.AssertExpectations(t)
}

// Test DeletePost - Post not found
func TestPostHandler_DeletePost_NotFound(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	handler := handlers.NewPostHandler(mockPostService)

	serviceError := errors.New("post no encontrado")
	mockPostService.On("DeletePost", 999, 1).Return(serviceError)

	// ACT
	req := httptest.NewRequest("DELETE", "/api/posts/999", nil)
	req.Header.Set("X-User-ID", "1")
	rr := httptest.NewRecorder()
	handler.DeletePost(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusForbidden, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, serviceError.Error(), response["error"])

	mockPostService.AssertExpectations(t)
}

// Test CreateComment - Success case
func TestPostHandler_CreateComment_Success(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	handler := handlers.NewPostHandler(mockPostService)

	createCommentRequest := models.CreateCommentRequest{
		Content: "This is a test comment",
	}

	expectedComment := &models.Comment{
		ID:      1,
		PostID:  1,
		Content: "This is a test comment",
	}

	mockPostService.On("CreateComment", 1, mock.AnythingOfType("*models.CreateCommentRequest"), 1).Return(expectedComment, nil)

	// ACT
	body, _ := json.Marshal(createCommentRequest)
	req := httptest.NewRequest("POST", "/api/posts/1/comments", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "1")
	rr := httptest.NewRecorder()
	handler.CreateComment(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response models.Comment
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedComment.Content, response.Content)
	assert.Equal(t, expectedComment.PostID, response.PostID)

	mockPostService.AssertExpectations(t)
}

// Test CreateComment - Invalid JSON
func TestPostHandler_CreateComment_InvalidJSON(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	handler := handlers.NewPostHandler(mockPostService)

	// ACT - Send invalid JSON
	invalidJSON := `{"content": "Test comment", "invalid": }`
	req := httptest.NewRequest("POST", "/api/posts/1/comments", bytes.NewBufferString(invalidJSON))
	req.Header.Set("X-User-ID", "1")
	rr := httptest.NewRecorder()
	handler.CreateComment(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "JSON inválido", response["error"])
}

// Test CreateComment - Service returns error
func TestPostHandler_CreateComment_ServiceError(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	handler := handlers.NewPostHandler(mockPostService)

	createCommentRequest := models.CreateCommentRequest{
		Content: "This is a test comment",
	}

	serviceError := errors.New("error creating comment")
	mockPostService.On("CreateComment", 1, mock.AnythingOfType("*models.CreateCommentRequest"), 1).Return(nil, serviceError)

	// ACT
	body, _ := json.Marshal(createCommentRequest)
	req := httptest.NewRequest("POST", "/api/posts/1/comments", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "1")
	rr := httptest.NewRecorder()
	handler.CreateComment(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, serviceError.Error(), response["error"])

	mockPostService.AssertExpectations(t)
}

// Test GetComments - Success case
func TestPostHandler_GetComments_Success(t *testing.T) {
	// ARRANGE
	mockPostService := new(mocks.MockPostService)
	handler := handlers.NewPostHandler(mockPostService)

	expectedComments := []*models.Comment{
		{
			ID:      1,
			PostID:  1,
			Content: "First comment",
		},
		{
			ID:      2,
			PostID:  1,
			Content: "Second comment",
		},
	}

	mockPostService.On("GetCommentsByPostID", 1).Return(expectedComments, nil)

	// ACT
	req := httptest.NewRequest("GET", "/api/posts/1/comments", nil)
	rr := httptest.NewRecorder()
	handler.GetComments(rr, req)

	// ASSERT
	assert.Equal(t, http.StatusOK, rr.Code)

	var response []*models.Comment
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response))
	assert.Equal(t, expectedComments[0].Content, response[0].Content)
	assert.Equal(t, expectedComments[1].Content, response[1].Content)

	mockPostService.AssertExpectations(t)
}
