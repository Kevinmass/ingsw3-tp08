package services

import (
	"ingsw3-tp08/internal/models"
)

// PostServiceInterface defines the contract for post services
type PostServiceInterface interface {
	CreatePost(req *models.CreatePostRequest, userID int) (*models.Post, error)
	GetAllPosts() ([]*models.Post, error)
	GetPostByID(id int) (*models.Post, error)
	DeletePost(postID int, userID int) error
	CreateComment(postID int, req *models.CreateCommentRequest, userID int) (*models.Comment, error)
	GetCommentsByPostID(postID int) ([]*models.Comment, error)
	DeleteComment(postID int, commentID int, userID int) error
}

// Ensure PostService implements the interface
var _ PostServiceInterface = (*PostService)(nil)
