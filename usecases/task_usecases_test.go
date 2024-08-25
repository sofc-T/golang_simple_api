package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/mock"
	"github.com/sofc-t/task_manager/task8/models"
	"github.com/sofc-t/task_manager/task8/usecases"
	"github.com/sofc-t/task_manager/task8/mocks"
)

type TaskUsecaseTestSuite struct {
	suite.Suite
	usecase         *usecases.TaskUsecase
	taskRepository  *mocks.TaskRepository
	ctx             context.Context
	timeout         time.Duration
}

func (suite *TaskUsecaseTestSuite) SetupTest() {
	suite.timeout = 1000 * time.Second
	suite.taskRepository = mocks.NewTaskRepository(suite.T())
	suite.usecase = usecases.NewTaskUsecase(suite.taskRepository, suite.timeout)
	suite.ctx = context.Background()
}

// TestFetch tests both successful and failed fetch operations
func (suite *TaskUsecaseTestSuite) TestFetch() {
	tasks := []models.Task{
		{Id: 1, Title: "Task 1", TaskId: 1},
		{Id: 2, Title: "Task 2", TaskId: 2},
	}

	suite.taskRepository.On("FetchTasks", mock.Anything).Return(tasks, nil).Once()
	result, err := suite.usecase.Fetch(suite.ctx)
	suite.NoError(err)
	suite.ElementsMatch(tasks, result)
	suite.taskRepository.AssertExpectations(suite.T())

	suite.taskRepository.On("FetchTasks", mock.Anything).Return(nil, errors.New("fetch error")).Once()
	result, err = suite.usecase.Fetch(suite.ctx)
	suite.Error(err)
	suite.Nil(result)
	suite.taskRepository.AssertExpectations(suite.T())
}

// TestFind tests both successful and failed find operations
func (suite *TaskUsecaseTestSuite) TestFind() {
	task := models.Task{Id: 1, Title: "Task 1", TaskId: 1}

	suite.taskRepository.On("FindTask", mock.Anything, 1).Return(task, nil).Once()
	result, err := suite.usecase.Find(suite.ctx, 1)
	suite.NoError(err)
	suite.Equal(task, result)
	suite.taskRepository.AssertExpectations(suite.T())

	suite.taskRepository.On("FindTask", mock.Anything, 1).Return(models.Task{}, errors.New("find error")).Once()
	result, err = suite.usecase.Find(suite.ctx, 1)
	suite.Error(err)
	suite.Equal(models.Task{}, result)
	suite.taskRepository.AssertExpectations(suite.T())
}

// TestUpdate tests both successful and failed update operations
func (suite *TaskUsecaseTestSuite) TestUpdate() {
	task := models.Task{Id: 1, Title: "Updated Task", TaskId: 1}

	suite.taskRepository.On("UpdateTask", mock.Anything, 1, "Updated Task").Return(task, nil).Once()
	result, err := suite.usecase.Update(suite.ctx, 1, "Updated Task")
	suite.NoError(err)
	suite.Equal(task, result)
	suite.taskRepository.AssertExpectations(suite.T())

	suite.taskRepository.On("UpdateTask", mock.Anything, 1, "Updated Task").Return(models.Task{}, errors.New("update error")).Once()
	result, err = suite.usecase.Update(suite.ctx, 1, "Updated Task")
	suite.Error(err)
	suite.Equal(models.Task{}, result)
	suite.taskRepository.AssertExpectations(suite.T())
}

// TestDelete tests both successful and failed delete operations
func (suite *TaskUsecaseTestSuite) TestDelete() {
	suite.taskRepository.On("DeleteTask", mock.Anything, 1).Return(nil).Once()
	err := suite.usecase.Delete(suite.ctx, 1)
	suite.NoError(err)
	suite.taskRepository.AssertExpectations(suite.T())

	suite.taskRepository.On("DeleteTask", mock.Anything, 1).Return(errors.New("delete error")).Once()
	err = suite.usecase.Delete(suite.ctx, 1)
	suite.Error(err)
	suite.taskRepository.AssertExpectations(suite.T())
}

// TestCreate tests both successful and failed create operations
func (suite *TaskUsecaseTestSuite) TestCreate() {
	task := models.Task{Id: 1, Title: "New Task", TaskId: 1}

	suite.taskRepository.On("InsertTask", mock.Anything, task).Return(task, nil).Once()
	result, err := suite.usecase.Create(suite.ctx, task)
	suite.NoError(err)
	suite.Equal(task, result)
	suite.taskRepository.AssertExpectations(suite.T())

	suite.taskRepository.On("InsertTask", mock.Anything, task).Return(models.Task{}, errors.New("create error")).Once()
	result, err = suite.usecase.Create(suite.ctx, task)
	suite.Error(err)
	suite.Equal(models.Task{}, result)
	suite.taskRepository.AssertExpectations(suite.T())
}

func TestTaskUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseTestSuite))
}
