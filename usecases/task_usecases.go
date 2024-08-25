package usecases

import (
	"context"
	"time"

	"github.com/sofc-t/task_manager/task8/models"
)

type TaskUsecase struct {
	taskRepository models.TaskRepository
	timeOut        time.Duration
}

func NewTaskUsecase(tasRepository models.TaskRepository, timeOut time.Duration) *TaskUsecase {
	return &TaskUsecase{
		taskRepository: tasRepository,
		timeOut:        timeOut,
	}
}

func (t TaskUsecase) Fetch(ctx context.Context) ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, t.timeOut)
	defer cancel()
	return t.taskRepository.FetchTasks(ctx)
}

func (t TaskUsecase) Find(ctx context.Context, id int) (models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, t.timeOut)
	defer cancel()
	return t.taskRepository.FindTask(ctx, id)

}

func (t TaskUsecase) Update(ctx context.Context, id int, title string) (models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, t.timeOut)
	defer cancel()
	return t.taskRepository.UpdateTask(ctx, id, title)

}

func (t TaskUsecase) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, t.timeOut)
	defer cancel()
	return t.taskRepository.DeleteTask(ctx, id)

}

func (t TaskUsecase) Create(ctx context.Context, task models.Task) (models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, t.timeOut)
	defer cancel()
	return t.taskRepository.InsertTask(ctx, task)

}
