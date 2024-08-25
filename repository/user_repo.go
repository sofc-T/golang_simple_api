package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/sofc-t/task_manager/task8/models"
	Utils "github.com/sofc-t/task_manager/task8/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
	adm   = "admin"
	admin = &adm
)

type userRepository struct {
	database   mongo.Database
	collection string
}

func NewUserRepository(db mongo.Database, collection string) models.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

func verifyPassword(existingPassword string, newPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(newPassword), []byte(existingPassword))
	if err != nil {
		log.Println("Password verification failed:", err)
	} else {
		log.Println("Password verification successful")
	}

	return err == nil
}

func (u userRepository) CreateUser(ctx context.Context, user models.User) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "email", Value: user.Email}}
	userCollection := u.database.Collection(u.collection)
	var existingUser models.User
	err := userCollection.FindOne(ctx, filter).Decode(&existingUser)
	fmt.Println(existingUser, *user.Email)
	if err == nil {
		log.Println("User already exists with email:", user.Email)
		return errors.New("user already exists")
	} else if err != mongo.ErrNoDocuments {
		log.Println("Error checking for existing user:", err)
		return errors.New("internal server error")
	}

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		log.Println("Error inserting user into database:", err)
		return errors.New("internal server error")
	}

	log.Println("User created successfully with email:", user.Email)
	return nil
}

func (u userRepository) Login(ctx context.Context, user models.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	userCollection := u.database.Collection(u.collection)
	filter := bson.D{{Key: "email", Value: user.Email}}
	var existingUser models.User

	log.Println("Attempting to find user with email:", user.Email)
	err := userCollection.FindOne(ctx, filter).Decode(&existingUser)
	if err != nil {
		log.Println("Error finding user:", err)
		return "", errors.New("user does not exist")
	}

	if existingUser.Email == nil {
		log.Println("User email is nil, user does not exist")
		return "", errors.New("user does not exist")
	}
	if existingUser.Email == nil || existingUser.Name == nil || existingUser.Role == nil {

		log.Println("One or more required fields are nil", user.Name, user.Email, user.Role, user)
		return "", errors.New("invalid user data")
	}

	log.Println("Verifying password for user:", user.Email)
	ok := verifyPassword(*user.Password, *existingUser.Password)
	if !ok {
		log.Println("Password verification failed for user:", user.Email)
		return "", errors.New("wrong password")
	}

	log.Println("Generating tokens for user:", user.Email)
	token, refreshToken, err := Utils.GenerateTokens(*existingUser.Email, *existingUser.Name, existingUser.UserID, *existingUser.Role)
	if err != nil {
		log.Println("Error generating tokens:", err)
		return "", errors.New("internal server error")
	}

	updateObj := Utils.UpdateAllTokens(*token, *refreshToken)
	upsert := true
	filters := bson.M{"user_id": existingUser.UserID}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	log.Println("Updating tokens for user with UserID:", existingUser.UserID)
	_, err = userCollection.UpdateOne(
		ctx,
		filters,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)

	if err != nil {
		log.Panic("Error updating user tokens:", err)
		return "", errors.New("internal server error, user not updated")
	}

	log.Println("Fetching updated user data for UserID:", existingUser.UserID)
	// var newuser models.User
	// err = userCollection.FindOne(ctx, bson.M{"UserID": existingUser.UserID}).Decode(&newuser)
	// if err != nil {
	// 	log.Println("Error fetching updated user data:", err)
	// 	return errors.New("internal server error")
	// }

	// log.Println("Login successful for user:", user.Email)
	return *token, nil
}

func (u userRepository) FetchAllUsers(ctx context.Context) ([]models.User, error) {
	userCollection := u.database.Collection(u.collection)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var users []models.User
	cursor, err := userCollection.Find(ctx, bson.D{})
	if err != nil {
		return users, errors.New("couldnt Fetch Data")
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &users); err != nil {
		return users, errors.New("couldnt Parse Data")
	}

	return users, nil
}

func (u userRepository) FetchByID(ctx context.Context, id string) (models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	userCollection := u.database.Collection(u.collection)
	var user models.User

	iD, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, errors.New("invalid ID")
	}

	err = userCollection.FindOne(ctx, bson.M{"_id": iD}).Decode(&user)
	defer cancel()

	if err != nil {
		return user, errors.New("error User not found")
	}

	return user, nil
}

func (u userRepository) PromoteUser(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	userCollection := u.database.Collection(u.collection)

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}
	var user models.User

	err = userCollection.FindOne(ctx, bson.M{"userid": userID}).Decode(&user)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Role == admin {
		return errors.New("user is already an admin")
	}

	user.Role = admin

	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"userid": userID},
		bson.D{
			{Key: "$set", Value: bson.M{"role": user.Role}},
		},
	)
	if err != nil {
		return errors.New("failed to promote user")
	}

	return nil
}
