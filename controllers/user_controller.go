package controllers

import (
	// "net/http"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sofc-t/task_manager/task8/models"
)

type UserController struct {
	UserUsecase models.UserUsecase
}

func (u UserController) SignUp(ctx *gin.Context) {
	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid Credentials"})
		return
	}

	err := u.UserUsecase.Create(ctx, user)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return

	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "Signed Up successfully"})

}

func (u UserController) Login(ctx *gin.Context) {
	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid Credentials"})
		return
	}

	token, err := u.UserUsecase.Login(ctx, user)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return

	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "Signed In successfully", "token": token})

}

func (u UserController) GetUseryID(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := u.UserUsecase.FetchById(ctx, id)
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusBadRequest, user)

}

func (u UserController) PromoteUser(ctx *gin.Context) {
	var req models.PromoteUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := req.ID
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err := u.UserUsecase.PromoteUser(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": "User updated successully"})

}
