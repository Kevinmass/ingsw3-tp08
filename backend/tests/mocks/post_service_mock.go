package mocks

import (
	"ingsw3-tp08/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockPostService es un mock del PostService para testing
type MockPostService struct {
	mock.Mock
}

// CreatePost simula la creación de un post
func (m *MockPostService) CreatePost(req *models.CreatePostRequest, userID int) (*models.Post, error) {
	args := m.Called(req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

// GetAllPosts simula obtener todos los posts
func (m *MockPostService) GetAllPosts() ([]*models.Post, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Post), args.Error(1)
}

// GetPostByID simula obtener un post por ID
func (m *MockPostService) GetPostByID(id int) (*models.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

// DeletePost simula eliminar un post
func (m *MockPostService) DeletePost(postID int, userID int) error {
	args := m.Called(postID, userID)
	return args.Error(0)
}

// CreateComment simula crear un comentario
func (m *MockPostService) CreateComment(postID int, req *models.CreateCommentRequest, userID int) (*models.Comment, error) {
	args := m.Called(postID, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Comment), args.Error(1)
}

// GetCommentsByPostID simula obtener comentarios por post ID
func (m *MockPostService) GetCommentsByPostID(postID int) ([]*models.Comment, error) {
	args := m.Called(postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Comment), args.Error(1)
}

// DeleteComment simula eliminar un comentario
func (m *MockPostService) DeleteComment(postID int, commentID int, userID int) error {
	args := m.Called(postID, commentID, userID)
	return args.Error(0)
}
