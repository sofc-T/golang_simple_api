package repository

import (
	"context"
	"errors"
	"log"

	"github.com/sofc-t/task_manager/task8/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepository struct {
	database   mongo.Database
	collection string
}

func NewTaskRepository(db mongo.Database, collection string) models.TaskRepository {
	return &taskRepository{
		database:   db,
		collection: collection,
	}
}

func (t *taskRepository) FetchTasks(ctx context.Context) ([]models.Task, error) {

	taskCollection := t.database.Collection(t.collection)
	cursor, err := taskCollection.Find(context.TODO(), bson.D{})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tasks []models.Task

	if err := cursor.All(context.TODO(), &tasks); err != nil {
		return nil, err
	}

	return tasks, nil

}

func (t *taskRepository) FindTask(ctx context.Context, id int) (models.Task, error) {

	filter := bson.D{{Key: "id", Value: id}}
	taskCollection := t.database.Collection(t.collection)
	var task models.Task
	err := taskCollection.FindOne(context.TODO(), filter).Decode(&task)

	if err != nil {
		return task, errors.New("failed to load Data")
	}

	return task, nil

}

func (t *taskRepository) UpdateTask(ctx context.Context, id int, title string) (models.Task, error) {

	taskCollection := t.database.Collection(t.collection)
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "title", Value: title}}}}

	_, err := taskCollection.UpdateOne(context.TODO(), filter, update)
	task := models.Task{Id: id, Title: title}

	if err != nil {
		return task, errors.New("failed to load Data")
	}
	return task, nil
}

func (t *taskRepository) DeleteTask(ctx context.Context, id int) error {
	taskCollection := t.database.Collection(t.collection)
	filter := bson.D{{Key: "id", Value: id}}

	_, err := taskCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return errors.New("failed to load Data")
	}
	return nil

}

func (t *taskRepository) InsertTask(ctx context.Context, task models.Task) (models.Task, error) {

	taskCollection := t.database.Collection(t.collection)

	_, err := taskCollection.InsertOne(context.TODO(), task)
	if err != nil {
		log.Printf("task not creatd")
		return task, errors.New("failed to load Data")
	}

	return task, nil

}
