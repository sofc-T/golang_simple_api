package models_test

import (
	"encoding/json"
	"testing"

	"github.com/sofc-t/task_manager/task8/models"
	"github.com/stretchr/testify/assert"
)

func TestTaskInitialization(t *testing.T) {
	task := models.Task{
		Id:     1,
		Title:  "Test Task",
		TaskId: 101,
	}

	assert.Equal(t, 1, task.Id, "Task ID should be 1")
	assert.Equal(t, "Test Task", task.Title, "Task title should be 'Test Task'")
	assert.Equal(t, 101, task.TaskId, "Task TaskId should be 101")
}

func TestTaskJSONSerialization(t *testing.T) {

	task := models.Task{
		Id:     1,
		Title:  "Test Task",
		TaskId: 101,
	}

	expectedJSON := `{"_id":1,"title":"Test Task","task_id":101}`
	taskJSON, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Error marshalling task to JSON: %v", err)
	}

	assert.JSONEq(t, expectedJSON, string(taskJSON), "JSON output should match expected output")
}
