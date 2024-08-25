package controllers_test

import (
	"bytes"
	"encoding/json"
	
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sofc-t/task_manager/task8/controllers"
	"github.com/sofc-t/task_manager/task8/mocks"
	"github.com/sofc-t/task_manager/task8/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TaskControllerSuite struct {
	suite.Suite
	controller    controllers.TaskController
	mockUsecase   *mocks.TaskUsecase
	mockResponse  *httptest.ResponseRecorder
}

func (suite *TaskControllerSuite) SetupSuite() {
	suite.mockUsecase = new(mocks.TaskUsecase)
	suite.controller = controllers.TaskController{TaskUsecase: suite.mockUsecase}
}

func (suite *TaskControllerSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.mockResponse = httptest.NewRecorder()
}

func (suite *TaskControllerSuite) TestGetAllTasksHandler_Success() {
	tasks := []models.Task{
		{Id: 1, Title: "Test Task 1"},
		{Id: 2, Title: "Test Task 2"},
	}
	suite.mockUsecase.On("Fetch", mock.Anything).Return(tasks, nil)

	ctx, _ := gin.CreateTestContext(suite.mockResponse)

	suite.controller.GetAllTasksHandler(ctx)

	suite.mockUsecase.AssertCalled(suite.T(), "Fetch", mock.Anything)
	suite.Equal(http.StatusAccepted, suite.mockResponse.Code)
	suite.Contains(suite.mockResponse.Body.String(), "Test Task 1")
	suite.Contains(suite.mockResponse.Body.String(), "Test Task 2")
}
func (suite *TaskControllerSuite) TestGetTaskHandler_Success() {
    taskId := 1
    task := models.Task{Id: taskId, Title: "Test Task"} 

    suite.mockUsecase.On("Find", mock.Anything, taskId).Return(task, nil)

    ctx, _ := gin.CreateTestContext(suite.mockResponse)
    ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(taskId)}}

    suite.controller.GetTaskHandler(ctx)

    suite.mockUsecase.AssertCalled(suite.T(), "Find", mock.Anything, taskId)

    suite.Equal(http.StatusOK, suite.mockResponse.Code)

}

func (suite *TaskControllerSuite) TestGetTaskHandler_InvalidId() {
	ctx, _ := gin.CreateTestContext(suite.mockResponse)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}

	suite.controller.GetTaskHandler(ctx)

	suite.Equal(http.StatusBadRequest, suite.mockResponse.Code)
	suite.Contains(suite.mockResponse.Body.String(), "Invalid ID")
}
func (suite *TaskControllerSuite) TestUpdateTaskHandler_Success() {
    taskId := 1
    task := models.Task{Id: taskId, Title: "Updated Task"}

    suite.mockUsecase.On("Update", mock.Anything, taskId, task.Title).Return(task, nil)
    ctx, _ := gin.CreateTestContext(suite.mockResponse)
    ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(taskId)}}

    jsonTask, _ := json.Marshal(task)
    ctx.Request, _ = http.NewRequest(http.MethodPut, "/tasks/"+strconv.Itoa(taskId), bytes.NewBuffer(jsonTask))
    ctx.Request.Header.Set("Content-Type", "application/json")

    suite.controller.UpdateTaskHandler(ctx)
    suite.mockUsecase.AssertCalled(suite.T(), "Update", mock.Anything, taskId, task.Title)
    suite.Equal(http.StatusAccepted, suite.mockResponse.Code)
    suite.Contains(suite.mockResponse.Body.String(), "Updated Task")
}


func (suite *TaskControllerSuite) TestUpdateTaskHandler_InvalidId() {
	ctx, _ := gin.CreateTestContext(suite.mockResponse)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}

	task := models.Task{Title: "Updated Task"}
	jsonTask, _ := json.Marshal(task)
	ctx.Request, _ = http.NewRequest(http.MethodPut, "/tasks/invalid", bytes.NewBuffer(jsonTask))
	ctx.Request.Header.Set("Content-Type", "application/json")

	suite.controller.UpdateTaskHandler(ctx)

	suite.Equal(http.StatusBadRequest, suite.mockResponse.Code)
	suite.Contains(suite.mockResponse.Body.String(), "Invalid ID")
}
func (suite *TaskControllerSuite) TestDeleteTaskHandler_Success() {
    taskId := 1
    
    suite.mockUsecase.On("Delete", mock.Anything, taskId).Return(nil)

    ctx, _ := gin.CreateTestContext(suite.mockResponse)
    ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(taskId)}}

    ctx.Request, _ = http.NewRequest(http.MethodDelete, "/tasks/"+strconv.Itoa(taskId), nil)

    suite.controller.DeleteTaskHandler(ctx)

    suite.mockUsecase.AssertCalled(suite.T(), "Delete", mock.Anything, taskId)

}



func (suite *TaskControllerSuite) TestDeleteTaskHandler_InvalidId() {
	suite.mockUsecase.On("Delete", mock.Anything, mock.AnythingOfType("int")).Return(nil)

suite.mockResponse = httptest.NewRecorder()
ctx, _ := gin.CreateTestContext(suite.mockResponse)
ctx.Request, _ = http.NewRequest(http.MethodDelete, "/tasks/invalid", nil)

suite.controller.DeleteTaskHandler(ctx)

suite.Equal(http.StatusBadRequest, suite.mockResponse.Code)


}
func (suite *TaskControllerSuite) TestCreateTaskHandler_Success() {
	task := models.Task{Id: 1, Title: "New Task"}
	suite.mockUsecase.On("Create", mock.Anything, mock.Anything).Return(models.Task{Id: 1, Title: "Test Task"}, nil)


	ctx, _ := gin.CreateTestContext(suite.mockResponse)

	jsonTask, _ := json.Marshal(task)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonTask))
	ctx.Request.Header.Set("Content-Type", "application/json")

	suite.controller.CreateTaskHandler(ctx)

	suite.mockUsecase.AssertCalled(suite.T(), "Create", mock.Anything, task)
	suite.Equal(http.StatusCreated, suite.mockResponse.Code)
	suite.Contains(suite.mockResponse.Body.String(), "Test Task")
}

func (suite *TaskControllerSuite) TestCreateTaskHandler_InvalidInput() {
	ctx, _ := gin.CreateTestContext(suite.mockResponse)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer([]byte("{invalid json}")))
	ctx.Request.Header.Set("Content-Type", "application/json")

	suite.controller.CreateTaskHandler(ctx)

	suite.Equal(http.StatusBadRequest, suite.mockResponse.Code)
	suite.Contains(suite.mockResponse.Body.String(), "Invalid Input")
}

func TestTaskControllerSuite(t *testing.T) {
	suite.Run(t, new(TaskControllerSuite))
}
