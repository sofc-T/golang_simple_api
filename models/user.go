package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	ID primitive.ObjectID 	`json:"_id" validate:"required,min=2"`
	Name *string `json:"name" validate:"required,min=2,max=25"` 
	Password *string `json:"password" validate:"required,min=6,max=100"`
	Email *string `json:"email" validate:"required,email"`
	Token *string `json:"token"`
	Role *string `json:"role"`
	RefreshToken *string `json:"refresh_token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
	UserID string `json:"user_id"`
}


type UserRepository interface {
	CreateUser(ctx context.Context, user User) error 
	Login(ctx context.Context, user User) (string, error)
	FetchAllUsers(ctx context.Context) ([]User, error) 
	FetchByID(ctx context.Context, id string) (User, error)
	PromoteUser(ctx context.Context, id string) error 
}


type UserUsecase interface {
	Create(ctx context.Context, user User) error 
	Login(ctx context.Context, user User) (string, error)
	FetchAll(ctx context.Context) ([]User, error) 
	FetchById(ctx context.Context, id string) (User, error)
	PromoteUser(ctx context.Context, id string) error 
}
type PromoteUserRequest struct {
    ID string `json:"id" binding:"required"`
}


