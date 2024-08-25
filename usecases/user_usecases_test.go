package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sofc-t/task_manager/task8/mocks"
	"github.com/sofc-t/task_manager/task8/models"
	"github.com/sofc-t/task_manager/task8/usecases"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	userRepo *mocks.UserRepository
	usecase  usecases.UserUsecase
}

func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.userRepo = &mocks.UserRepository{}
	suite.usecase = *usecases.NewUserUsecase(suite.userRepo, 1000 * time.Second)
}

// func (suite *UserUsecaseTestSuite) TestCreateUser_Success() {
// 	// Arrange

// 	user := models.User{
// 		ID:    primitive.NewObjectID(),
// 		Name:  stringPtr(name),
// 		Email: stringPtr(email),
// 	}

// 	suite.userRepo.On("CreateUser", mock.Anything, user).Return(nil).Once()

// 	// Act
// 	err := suite.usecase.Create(context.TODO(), user)

// 	// Assert
// 	suite.NoError(err)
// 	suite.userRepo.AssertExpectations(suite.T())
// }





func (suite *UserUsecaseTestSuite) TestCreateUser_Success() {
		p := "password123"
			name:= "John Doe"
			email := "john.doe@example.com"
			amdin := "admin"
			user := models.User{
			Name:         stringPtr(name),
			Password:     stringPtr(p),
			Email:        stringPtr(email),
			Role:         stringPtr(amdin),
			CreatedAt:    time.Now(),
			UpdatedAT:    time.Now(),
			UserID:       "66bd0d1f16be6a49b496a9c7",
		}
	
		ctx :=  context.Background																																																																												()
		suite.userRepo.On("CreateUser", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("models.User")).Return(nil)
	
		err := suite.usecase.Create(ctx, user)
	
		assert.NoError(suite.T(), err)
		suite.userRepo.AssertExpectations(suite.T())
	}

func (suite *UserUsecaseTestSuite) TestLogin_Success() {
	p := "password123"
	user := models.User{
		Email:    stringPointer("test@example.com"),
		Password: stringPointer(p),
	}
	expectedToken := "someToken"

	suite.userRepo.On("Login", mock.Anything, user).Return(expectedToken, nil)

	token, err := suite.usecase.Login(context.Background(), user)
	suite.NoError(err)
	suite.Equal(expectedToken, token)
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserUsecaseTestSuite) TestLogin_Error() {
	user := models.User{
		Email:    stringPointer("test@example.com"),
		Password: stringPointer("password123"),
	}

	suite.userRepo.On("Login", mock.Anything, user).Return("", errors.New("error"))

	token, err := suite.usecase.Login(context.Background(), user)
	suite.Error(err)
	suite.Empty(token)
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserUsecaseTestSuite) TestFetchAll_Success() {
	users := []models.User{
		{Email: stringPointer("test1@example.com")},
		{Email: stringPointer("test2@example.com")},
	}

	suite.userRepo.On("FetchAllUsers", mock.Anything).Return(users, nil)

	result, err := suite.usecase.FetchAll(context.Background())
	suite.NoError(err)
	suite.Len(result, 2)
	suite.Equal(users, result)
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserUsecaseTestSuite) TestFetchById_Success() {
	ID := primitive.NewObjectID().Hex()
	user := models.User{UserID: ID}

	suite.userRepo.On("FetchByID", mock.Anything, ID).Return(user, nil)

	result,  err := suite.usecase.FetchById(context.Background(), ID)
	
	suite.NoError(err)
	
	suite.Equal(result, user)
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserUsecaseTestSuite) TestPromoteUser_Success() {
	userID := "someID"

	suite.userRepo.On("PromoteUser", mock.Anything, userID).Return(nil)

	err := suite.usecase.PromoteUser(context.Background(), userID)
	suite.NoError(err)
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserUsecaseTestSuite) TestPromoteUser_Error() {
	userID := "someID"

	suite.userRepo.On("PromoteUser", mock.Anything, userID).Return(errors.New("error"))

	err := suite.usecase.PromoteUser(context.Background(), userID)
	suite.Error(err)
	suite.userRepo.AssertExpectations(suite.T())
}

func stringPointer(s string) *string {
	return &s
}


func stringPtr(s string) *string {
	return &s
}
func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}