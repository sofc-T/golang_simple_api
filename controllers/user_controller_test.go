package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sofc-t/task_manager/task8/controllers"
	"github.com/sofc-t/task_manager/task8/mocks"
	"github.com/sofc-t/task_manager/task8/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserControllerSuite struct {
	suite.Suite
	controller    controllers.UserController
	mockUsecase   *mocks.UserUsecase
	mockResponse  *httptest.ResponseRecorder
}

func (suite *UserControllerSuite) SetupSuite() {
	suite.mockUsecase = new(mocks.UserUsecase)
	suite.controller = controllers.UserController{UserUsecase: suite.mockUsecase}
}

func (suite *UserControllerSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.mockResponse = httptest.NewRecorder()
}


func (suite *UserControllerSuite) TestSignUpHandler_Success() {
	user := models.User{Name: stringPointer("Test User"), Password: stringPointer("Password123"), Email: stringPointer("test@example.com")}
	suite.mockUsecase.On("Create", mock.Anything, user).Return(nil)

	ctx, _ := gin.CreateTestContext(suite.mockResponse)
	jsonUser, _ := json.Marshal(user)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonUser))
	ctx.Request.Header.Set("Content-Type", "application/json")

	suite.controller.SignUp(ctx)

	suite.mockUsecase.AssertCalled(suite.T(), "Create", mock.Anything, user)
	suite.Equal(http.StatusCreated, suite.mockResponse.Code)
	suite.Contains(suite.mockResponse.Body.String(), "Signed Up successfully")
}


func (suite *UserControllerSuite) TestSignUpHandler_Failure() {
    
    req := models.User{ /* Populate with test data as needed */ }
    expectedError := errors.New("sign up error")

    
    suite.mockUsecase.On("Create", mock.Anything, req).Return(expectedError)

    
    ctx, _ := gin.CreateTestContext(suite.mockResponse)
    jsonReq, _ := json.Marshal(req)
    ctx.Request, _ = http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonReq))
    ctx.Request.Header.Set("Content-Type", "application/json")

    
    suite.controller.SignUp(ctx)

    
    suite.mockUsecase.AssertCalled(suite.T(), "Create", mock.Anything, req)

    
    suite.Equal(http.StatusInternalServerError, suite.mockResponse.Code)
    suite.Contains(suite.mockResponse.Body.String(), "sign up error")
}



func (suite *UserControllerSuite) TestLoginHandler_Success() {
	user := models.User{Email: stringPointer("test@example.com"), Password: stringPointer("Password123")}
	token := "some-token"
	suite.mockUsecase.On("Login", mock.Anything, user).Return(token, nil)

	ctx, _ := gin.CreateTestContext(suite.mockResponse)
	jsonUser, _ := json.Marshal(user)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonUser))
	ctx.Request.Header.Set("Content-Type", "application/json")

	suite.controller.Login(ctx)

	suite.mockUsecase.AssertCalled(suite.T(), "Login", mock.Anything, user)
	suite.Equal(http.StatusCreated, suite.mockResponse.Code)
	suite.Contains(suite.mockResponse.Body.String(), "Signed In successfully")
	suite.Contains(suite.mockResponse.Body.String(), token)
}





func (suite *UserControllerSuite) TestGetUserByIDHandler_Success() {
	userID := primitive.NewObjectID()
	user := models.User{ID: userID} 

	suite.mockUsecase.On("FetchById", mock.Anything, userID.Hex()).Return(user, nil)

	ctx, _ := gin.CreateTestContext(suite.mockResponse)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: userID.Hex()}}

	suite.controller.GetUseryID(ctx)

	suite.mockUsecase.AssertCalled(suite.T(), "FetchById", mock.Anything, userID.Hex())
	
	suite.Contains(suite.mockResponse.Body.String(), userID.Hex())
}



func (suite *UserControllerSuite) TestGetUserByIDHandler_Failure() {
	userID := "" 
	expectedError := errors.New("Invalid user ID")

	
	suite.mockUsecase.On("FetchById", mock.Anything, userID).Return(models.User{}, expectedError)

	ctx, _ := gin.CreateTestContext(suite.mockResponse)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: userID}}

	
	suite.controller.GetUseryID(ctx)

	
	

	
	suite.Equal(http.StatusBadRequest, suite.mockResponse.Code)
	suite.Contains(suite.mockResponse.Body.String(), expectedError.Error())
}


func (suite *UserControllerSuite) TestPromoteUserHandler_Success() {
	req := models.PromoteUserRequest{ID: "some-id"}
	suite.mockUsecase.On("PromoteUser", mock.Anything, req.ID).Return(nil)

	ctx, _ := gin.CreateTestContext(suite.mockResponse)
	jsonReq, _ := json.Marshal(req)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/promote", bytes.NewBuffer(jsonReq))
	ctx.Request.Header.Set("Content-Type", "application/json")

	suite.controller.PromoteUser(ctx)

	suite.mockUsecase.AssertCalled(suite.T(), "PromoteUser", mock.Anything, req.ID)
	suite.Equal(http.StatusAccepted, suite.mockResponse.Code)
	suite.Contains(suite.mockResponse.Body.String(), "User updated successully")
}



func stringPointer(s string) *string {
	return &s
}

func TestUserControllerSuite(t *testing.T) {
	suite.Run(t, new(UserControllerSuite))
}
