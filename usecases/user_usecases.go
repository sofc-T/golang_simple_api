package usecases

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/sofc-t/task_manager/task8/models"
	Utils "github.com/sofc-t/task_manager/task8/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	UserRepositry models.UserRepository
	TimeOut       time.Duration
}

var (
	adminCreated bool
	adm          = "admin"
	gue          = "user"
	admin        = &adm
	guest        = &gue
)

func NewUserUsecase(UserRepositry models.UserRepository, TimeOut time.Duration) *UserUsecase {
	return &UserUsecase{
		UserRepositry: UserRepositry,
		TimeOut:       TimeOut,
	}
}

func hashPassword(password string) (string, error) {
	log.Println("hashing", password)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		log.Println("couldnt Generate Password")
		return "", errors.New("couldnt Generate Password")
	}
	log.Println("hashed")
	return string(bytes), nil
}

func (u UserUsecase) Create(ctx context.Context, user models.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.TimeOut)
	defer cancel()

	if !adminCreated {
		adminCreated = true
		user.Role = admin
		log.Println("Assigned Admin role to new user")
	} else {
		user.Role = guest
		log.Println("Assigned role : user")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAT = time.Now()
	user.ID = primitive.NewObjectID()
	user.UserID = user.ID.Hex()
	log.Println("id generated")

	password, err := hashPassword(*user.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return errors.New("couldn't parse password")
	}
	user.Password = &password

	token, refreshToken, err := Utils.GenerateTokens(*user.Email, *user.Name, user.UserID, *user.Role)
	if err != nil {
		log.Println("Error generating tokens:", err)
		return errors.New("internal server error")
	}

	user.Token = token
	user.RefreshToken = refreshToken

	return u.UserRepositry.CreateUser(ctx, user)

}

func (u UserUsecase) Login(ctx context.Context, user models.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, u.TimeOut)
	defer cancel()
	return u.UserRepositry.Login(ctx, user)

}

func (u UserUsecase) FetchAll(ctx context.Context) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.TimeOut)
	defer cancel()
	return u.UserRepositry.FetchAllUsers(ctx)

}

func (u UserUsecase) FetchById(ctx context.Context, id string) (models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.TimeOut)
	defer cancel()
	return u.UserRepositry.FetchByID(ctx, id)

}

func (u UserUsecase) PromoteUser(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.TimeOut)
	defer cancel()
	return u.UserRepositry.PromoteUser(ctx, id)

}
