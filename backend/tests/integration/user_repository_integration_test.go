package integration

import (
	"database/sql"
	"testing"

	"ingsw3-tp08/internal/models"
	"ingsw3-tp08/internal/repository"

	_ "github.com/lib/pq"

	"github.com/stretchr/testify/suite"
)

type UserRepositoryIntegrationTestSuite struct {
	suite.Suite
	db        *sql.DB
	repo      repository.UserRepository
	cleanupDB func()
}

func (suite *UserRepositoryIntegrationTestSuite) SetupTest() {
	// Setup database for each test
	db, cleanup, err := SetupTestDB()
	suite.Require().NoError(err)

	suite.db = db
	suite.repo = repository.NewPostgreSQLUserRepository(db)
	suite.cleanupDB = cleanup

	// Clean tables before each test
	err = CleanupTestDB(db)
	suite.Require().NoError(err)
}

func (suite *UserRepositoryIntegrationTestSuite) TearDownTest() {
	// Clean up after each test
	if suite.cleanupDB != nil {
		suite.cleanupDB()
	}
}

func (suite *UserRepositoryIntegrationTestSuite) TestCreate_Success() {
	// Prepare test data
	user := &models.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
		Username: "testuser",
	}

	// Execute
	err := suite.repo.Create(user)

	// Assert
	suite.NoError(err)
	suite.NotZero(user.ID) // Should have auto-generated ID

	// Verify can find the user
	found, err := suite.repo.FindByEmail("test@example.com")
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(user.ID, found.ID)
	suite.Equal(user.Email, found.Email)
	suite.Equal(user.Username, found.Username)
}

func (suite *UserRepositoryIntegrationTestSuite) TestCreate_DuplicateEmail() {
	// Create first user
	user1 := &models.User{
		Email:    "duplicate@example.com",
		Password: "pass1",
		Username: "user1",
	}
	err := suite.repo.Create(user1)
	suite.NoError(err)

	// Try to create second user with same email (should fail due to unique constraint)
	user2 := &models.User{
		Email:    "duplicate@example.com",
		Password: "pass2",
		Username: "user2",
	}
	err = suite.repo.Create(user2)
	suite.Error(err) // Should fail due to duplicate email
}

func (suite *UserRepositoryIntegrationTestSuite) TestFindByEmail_Exists() {
	// Create user
	user := &models.User{
		Email:    "findme@example.com",
		Password: "password",
		Username: "findme",
	}
	err := suite.repo.Create(user)
	suite.NoError(err)

	// Find by email
	found, err := suite.repo.FindByEmail("findme@example.com")

	// Assert
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(user.ID, found.ID)
	suite.Equal(user.Email, found.Email)
	suite.Equal(user.Username, found.Username)
}

func (suite *UserRepositoryIntegrationTestSuite) TestFindByEmail_NotExists() {
	// Try to find non-existent user
	found, err := suite.repo.FindByEmail("nonexistent@example.com")

	// Assert
	suite.NoError(err)
	suite.Nil(found)
}

func (suite *UserRepositoryIntegrationTestSuite) TestFindByID_Exists() {
	// Create user
	user := &models.User{
		Email:    "findbyid@example.com",
		Password: "password",
		Username: "findbyid",
	}
	err := suite.repo.Create(user)
	suite.NoError(err)

	// Find by ID
	found, err := suite.repo.FindByID(user.ID)

	// Assert
	suite.NoError(err)
	suite.NotNil(found)
	suite.Equal(user.ID, found.ID)
	suite.Equal(user.Email, found.Email)
	suite.Equal(user.Username, found.Username)
}

func (suite *UserRepositoryIntegrationTestSuite) TestFindByID_NotExists() {
	// Try to find non-existent user by ID
	found, err := suite.repo.FindByID(99999)

	// Assert
	suite.NoError(err)
	suite.Nil(found)
}

func TestUserRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryIntegrationTestSuite))
}
