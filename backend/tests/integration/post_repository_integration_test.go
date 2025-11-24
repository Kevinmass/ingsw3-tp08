package integration

import (
	"database/sql"
	"testing"

	"ingsw3-tp08/internal/models"
	"ingsw3-tp08/internal/repository"

	"github.com/stretchr/testify/suite"
)

type PostRepositoryIntegrationTestSuite struct {
	suite.Suite
	db        *sql.DB
	repo      repository.PostRepository
	userRepo  repository.UserRepository
	cleanupDB func()
}

func (suite *PostRepositoryIntegrationTestSuite) SetupTest() {
	// Setup database for each test
	db, cleanup, err := SetupTestDB()
	suite.Require().NoError(err)

	suite.db = db
	suite.repo = repository.NewPostgreSQLPostRepository(db)
	suite.userRepo = repository.NewPostgreSQLUserRepository(db)
	suite.cleanupDB = cleanup

	// Clean tables before each test
	err = CleanupTestDB(db)
	suite.Require().NoError(err)
}

func (suite *PostRepositoryIntegrationTestSuite) TearDownTest() {
	// Clean up after each test
	if suite.cleanupDB != nil {
		suite.cleanupDB()
	}
}

func (suite *PostRepositoryIntegrationTestSuite) TestCreate_Success() {
	// Create a user first
	user := &models.User{
		Email:    "user@example.com",
		Password: "password",
		Username: "testuser",
	}
	err := suite.userRepo.Create(user)
	suite.NoError(err)

	// Prepare test data
	post := &models.Post{
		Title:   "Test Post",
		Content: "Test Content",
		UserID:  user.ID,
	}

	// Execute
	err = suite.repo.Create(post)

	// Assert
	suite.NoError(err)
	suite.NotZero(post.ID)

	// Verify can find the post
	found, err := suite.repo.FindByID(post.ID)
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(post.ID, found.ID)
	suite.Equal(post.Title, found.Title)
	suite.Equal(post.Content, found.Content)
	suite.Equal(post.UserID, found.UserID)
	suite.Equal(user.Username, found.Username)
}

func (suite *PostRepositoryIntegrationTestSuite) TestFindAll_Empty() {
	// Execute
	posts, err := suite.repo.FindAll()

	// Assert
	suite.NoError(err)
	suite.Len(posts, 0)
}

func (suite *PostRepositoryIntegrationTestSuite) TestFindAll_WithPosts() {
	// Create user
	user := &models.User{
		Email:    "user@example.com",
		Password: "password",
		Username: "testuser",
	}
	err := suite.userRepo.Create(user)
	suite.NoError(err)

	// Create posts
	post1 := &models.Post{
		Title:   "Post 1",
		Content: "Content 1",
		UserID:  user.ID,
	}
	err = suite.repo.Create(post1)
	suite.NoError(err)

	post2 := &models.Post{
		Title:   "Post 2",
		Content: "Content 2",
		UserID:  user.ID,
	}
	err = suite.repo.Create(post2)
	suite.NoError(err)

	// Execute
	posts, err := suite.repo.FindAll()

	// Assert
	suite.NoError(err)
	suite.Len(posts, 2)
	// Should be ordered DESC by created_at
	suite.Equal(post2.Title, posts[0].Title) // Most recent first
	suite.Equal(post1.Title, posts[1].Title)
}

func (suite *PostRepositoryIntegrationTestSuite) TestFindByID_Exists() {
	// Create user and post
	user := &models.User{
		Email:    "user@example.com",
		Password: "password",
		Username: "testuser",
	}
	err := suite.userRepo.Create(user)
	suite.NoError(err)

	post := &models.Post{
		Title:   "Test Post",
		Content: "Test Content",
		UserID:  user.ID,
	}
	err = suite.repo.Create(post)
	suite.NoError(err)

	// Execute
	found, err := suite.repo.FindByID(post.ID)

	// Assert
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(post.ID, found.ID)
	suite.Equal(post.Title, found.Title)
	suite.Equal(user.Username, found.Username)
}

func (suite *PostRepositoryIntegrationTestSuite) TestFindByID_NotExists() {
	// Execute
	found, err := suite.repo.FindByID(99999)

	// Assert
	suite.NoError(err)
	suite.Nil(found)
}

func (suite *PostRepositoryIntegrationTestSuite) TestDelete_Success() {
	// Create user and post
	user := &models.User{
		Email:    "user@example.com",
		Password: "password",
		Username: "testuser",
	}
	err := suite.userRepo.Create(user)
	suite.NoError(err)

	post := &models.Post{
		Title:   "Test Post",
		Content: "Test Content",
		UserID:  user.ID,
	}
	err = suite.repo.Create(post)
	suite.NoError(err)

	// Execute delete
	err = suite.repo.Delete(post.ID)

	// Assert no error
	suite.NoError(err)

	// Verify post is gone
	found, err := suite.repo.FindByID(post.ID)
	suite.NoError(err)
	suite.Nil(found)
}

func (suite *PostRepositoryIntegrationTestSuite) TestCreateComment_Success() {
	// Create user and post
	user := &models.User{
		Email:    "user@example.com",
		Password: "password",
		Username: "testuser",
	}
	err := suite.userRepo.Create(user)
	suite.NoError(err)

	post := &models.Post{
		Title:   "Test Post",
		Content: "Test Content",
		UserID:  user.ID,
	}
	err = suite.repo.Create(post)
	suite.NoError(err)

	// Prepare comment
	comment := &models.Comment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: "Test Comment",
	}

	// Execute
	err = suite.repo.CreateComment(comment)

	// Assert
	suite.NoError(err)
	suite.NotZero(comment.ID)

	// Verify can find the comment
	comments, err := suite.repo.FindCommentsByPostID(post.ID)
	suite.NoError(err)
	suite.Len(comments, 1)
	suite.Equal(comment.ID, comments[0].ID)
	suite.Equal(comment.Content, comments[0].Content)
	suite.Equal(user.Username, comments[0].Username)
}

func (suite *PostRepositoryIntegrationTestSuite) TestFindCommentsByPostID_Empty() {
	// Create user and post
	user := &models.User{
		Email:    "user@example.com",
		Password: "password",
		Username: "testuser",
	}
	err := suite.userRepo.Create(user)
	suite.NoError(err)

	post := &models.Post{
		Title:   "Test Post",
		Content: "Test Content",
		UserID:  user.ID,
	}
	err = suite.repo.Create(post)
	suite.NoError(err)

	// Execute
	comments, err := suite.repo.FindCommentsByPostID(post.ID)

	// Assert
	suite.NoError(err)
	suite.Len(comments, 0)
}

func (suite *PostRepositoryIntegrationTestSuite) TestDeleteComment_Success() {
	// Create user and post
	user := &models.User{
		Email:    "user@example.com",
		Password: "password",
		Username: "testuser",
	}
	err := suite.userRepo.Create(user)
	suite.NoError(err)

	post := &models.Post{
		Title:   "Test Post",
		Content: "Test Content",
		UserID:  user.ID,
	}
	err = suite.repo.Create(post)
	suite.NoError(err)

	// Create comment
	comment := &models.Comment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: "Test Comment",
	}
	err = suite.repo.CreateComment(comment)
	suite.NoError(err)

	// Execute delete
	err = suite.repo.DeleteComment(post.ID, comment.ID, user.ID)

	// Assert
	suite.NoError(err)

	// Verify comment is gone
	comments, err := suite.repo.FindCommentsByPostID(post.ID)
	suite.NoError(err)
	suite.Len(comments, 0)
}

func (suite *PostRepositoryIntegrationTestSuite) TestDeleteComment_NotAuthorized() {
	// Create two users
	user1 := &models.User{
		Email:    "user1@example.com",
		Password: "password",
		Username: "user1",
	}
	err := suite.userRepo.Create(user1)
	suite.NoError(err)

	user2 := &models.User{
		Email:    "user2@example.com",
		Password: "password",
		Username: "user2",
	}
	err = suite.userRepo.Create(user2)
	suite.NoError(err)

	// Create post by user1
	post := &models.Post{
		Title:   "Test Post",
		Content: "Test Content",
		UserID:  user1.ID,
	}
	err = suite.repo.Create(post)
	suite.NoError(err)

	// Create comment by user1
	comment := &models.Comment{
		PostID:  post.ID,
		UserID:  user1.ID,
		Content: "Test Comment",
	}
	err = suite.repo.CreateComment(comment)
	suite.NoError(err)

	// Try to delete as user2 (should fail)
	err = suite.repo.DeleteComment(post.ID, comment.ID, user2.ID)

	// Assert
	suite.Error(err)

	// Verify comment still exists
	comments, err := suite.repo.FindCommentsByPostID(post.ID)
	suite.NoError(err)
	suite.Len(comments, 1)
}

func TestPostRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(PostRepositoryIntegrationTestSuite))
}
