package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sofc-t/task_manager/task8/models"
)

type TaskController struct {
	TaskUsecase models.TaskUsecase
}

func (tc *TaskController) GetAllTasksHandler(ctx *gin.Context) {

	tasks, err := tc.TaskUsecase.Fetch(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch services"})
	}
	ctx.IndentedJSON(http.StatusAccepted, tasks)

}

func (tc *TaskController) GetTaskHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	task, err := tc.TaskUsecase.Find(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch services"})
	}

	ctx.IndentedJSON(http.StatusOK, task)
}

func (tc *TaskController) UpdateTaskHandler(ctx *gin.Context) {
	var task models.Task
	if err := ctx.BindJSON(&task); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invaid Task Information"})
		return
	}
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	result, err := tc.TaskUsecase.Update(ctx, id, task.Title)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch services"})
	}

	ctx.IndentedJSON(http.StatusAccepted, result)

}

func (tc *TaskController) DeleteTaskHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": ctx.Request.Body})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"id": id, "err": err})
		return
	}

	err = tc.TaskUsecase.Delete(ctx, id)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Could not Delete Task"})
		return
	}

	ctx.Status(http.StatusNoContent)

}

func (tc *TaskController) CreateTaskHandler(ctx *gin.Context) {
	var task models.Task
	err := ctx.BindJSON(&task)

	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid Input"})
		return
	}

	result, err := tc.TaskUsecase.Create(ctx, task)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Could not Delete Task"})
		return
	}
	ctx.IndentedJSON(http.StatusCreated, result)

}
