package routers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sofc-t/task_manager/task8/controllers"
	"github.com/sofc-t/task_manager/task8/middleware"
	"github.com/sofc-t/task_manager/task8/repository"
	"github.com/sofc-t/task_manager/task8/usecases"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRouter(timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(*db, "user")
	userUsecase := usecases.NewUserUsecase(userRepo, timeout)
	userController := &controllers.UserController{
		UserUsecase: userUsecase,
	}

	group.POST("/register", userController.SignUp)
	group.POST("/login", userController.Login)
}

func NewTaskRouter(timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	taskRepo := repository.NewTaskRepository(db, "task")
	taskUsecase := usecases.NewTaskUsecase(taskRepo, timeout)
	taskController := &controllers.TaskController{
		TaskUsecase: taskUsecase,
	}

	group.GET("/tasks/:id", taskController.GetTaskHandler)
	group.PUT("/tasks/:id", taskController.UpdateTaskHandler)
	group.POST("/tasks", taskController.CreateTaskHandler)

}

func NewAdminRouter(timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	taskRepo := repository.NewTaskRepository(db, "task")
	taskUsecase := usecases.NewTaskUsecase(taskRepo, timeout)
	taskController := &controllers.TaskController{
		TaskUsecase: taskUsecase,
	}

	userRepo := repository.NewUserRepository(db, "user")
	userUsecase := usecases.NewUserUsecase(userRepo, timeout)
	userController := &controllers.UserController{
		UserUsecase: userUsecase,
	}

	group.GET("/tasks", taskController.GetAllTasksHandler)
	group.DELETE("/tasks/:id", taskController.DeleteTaskHandler)
	group.POST("/promote", userController.PromoteUser)
}

func SetUpRouter(timeout time.Duration, db mongo.Database) *gin.Engine {
	r := gin.Default()

	publicGroup := r.Group("/")
	NewUserRouter(timeout, &db, publicGroup)

	protectedGroup := r.Group("/")
	protectedGroup.Use(middleware.AuthenticationMiddleware())
	NewTaskRouter(timeout, db, protectedGroup)

	adminGroup := r.Group("/")
	adminGroup.Use(middleware.AuthenticationandAuthorizeMiddleware())
	NewAdminRouter(timeout, db, adminGroup)

	return r
}
