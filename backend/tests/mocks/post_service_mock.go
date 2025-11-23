package mocks

import (
	"ingsw3-tp08/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockPostService is a mock implementation of PostService
type MockPostService struct {
	mock.Mock
}

// CreatePost mocks the CreatePost method
func (m *MockPostService) CreatePost(req *models.CreatePostRequest, userID int) (*models.Post, error) {
	args := m.Called(req, userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Post), args.Error(1)
}

// GetAllPosts mocks the GetAllPosts method
func (m *MockPostService) GetAllPosts() ([]*models.Post, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Post), args.Error(1)
}

// GetPostByID mocks the GetPostByID method
func (m *MockPostService) GetPostByID(id int) (*models.Post, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Post), args.Error(1)
}

// DeletePost mocks the DeletePost method
func (m *MockPostService) DeletePost(postID int, userID int) error {
	args := m.Called(postID, userID)
	return args.Error(0)
}

// CreateComment mocks the CreateComment method
func (m *MockPostService) CreateComment(postID int, req *models.CreateCommentRequest, userID int) (*models.Comment, error) {
	args := m.Called(postID, req, userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Comment), args.Error(1)
}

// GetCommentsByPostID mocks the GetCommentsByPostID method
func (m *MockPostService) GetCommentsByPostID(postID int) ([]*models.Comment, error) {
	args := m.Called(postID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*models.Comment), args.Error(1)
}

// DeleteComment mocks the DeleteComment method
func (m *MockPostService) DeleteComment(postID int, commentID int, userID int) error {
	args := m.Called(postID, commentID, userID)
	return args.Error(0)
}
