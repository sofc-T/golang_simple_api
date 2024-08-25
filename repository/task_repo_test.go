package repository_test

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/sofc-t/task_manager/task8/models"
	"github.com/sofc-t/task_manager/task8/repository"
)

type TaskRepositorySuite struct {
	suite.Suite
	TaskRepository models.TaskRepository
	collection     *mongo.Collection
}

func (suite *TaskRepositorySuite) SetupSuite() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Replace with your MongoDB URI
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	database := client.Database("taskmanager_test") // Use a test database
	collection := database.Collection("tasks")
	suite.collection = collection
	collection.DeleteMany(context.TODO(), bson.D{{}})
	suite.TaskRepository = repository.NewTaskRepository(*database, "tasks")
}

func (suite *TaskRepositorySuite) SetupTest() {
	
}

// Tests InsertTask function
func (suite *TaskRepositorySuite) TestInsertTask() {
	task := models.Task{Id: 1, Title: "Test Task"}
	insertedTask, err := suite.TaskRepository.InsertTask(context.TODO(), task)
	suite.NoError(err, "no error when inserting task")
	suite.Equal(task, insertedTask)
}

// Tests FetchTasks function
func (suite *TaskRepositorySuite) TestFetchTasks() {
	task := models.Task{Id: 1, Title: "Test Task"}
	_, err := suite.TaskRepository.InsertTask(context.TODO(), task)
	suite.NoError(err, "no error when inserting task")

	tasks, err := suite.TaskRepository.FetchTasks(context.TODO())
	suite.NoError(err, "no error when fetching tasks")
	suite.Len(tasks, 1)
	suite.Equal("Test Task", tasks[0].Title)
}

// Tests FindTask function
func (suite *TaskRepositorySuite) TestFindTask() {
	task := models.Task{Id: 1, Title: "Test Task"}
	_, err := suite.TaskRepository.InsertTask(context.TODO(), task)
	suite.NoError(err, "no error when inserting task")

	foundTask, err := suite.TaskRepository.FindTask(context.TODO(), 1)
	suite.NoError(err, "no error when finding task")
	suite.Equal(task, foundTask)
}

// Tests UpdateTask function
func (suite *TaskRepositorySuite) TestUpdateTask() {
	task := models.Task{Id: 1, Title: "Test Task"}
	_, err := suite.TaskRepository.InsertTask(context.TODO(), task)
	suite.NoError(err, "no error when inserting task")

	updatedTask, err := suite.TaskRepository.UpdateTask(context.TODO(), 1, "Updated Task")
	suite.NoError(err, "no error when updating task")
	suite.Equal("Updated Task", updatedTask.Title)
}

// Tests DeleteTask function
func (suite *TaskRepositorySuite) TestDeleteTask() {
	task := models.Task{Id: 1, Title: "Test Task"}
	_, err := suite.TaskRepository.InsertTask(context.TODO(), task)
	suite.NoError(err, "no error when inserting task")

	err = suite.TaskRepository.DeleteTask(context.TODO(), 1)
	suite.NoError(err, "no error when deleting task")

	tasks, err := suite.TaskRepository.FetchTasks(context.TODO())
	suite.NoError(err, "no error when fetching tasks")
	suite.Len(tasks, 0, "no tasks should remain")
}

func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositorySuite))
}
