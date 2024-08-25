package repository_test

import (
	"context"
	"testing"

	"github.com/sofc-t/task_manager/task8/models"
	"github.com/sofc-t/task_manager/task8/repository"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepositorySuite struct {
    suite.Suite
    repo      models.UserRepository
    db        *mongo.Database
    collection *mongo.Collection
}

func (suite *UserRepositorySuite) SetupTest() {
    // Mock database setup
    clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), clientOpts)
    suite.Require().NoError(err)

    suite.db = client.Database("test_db")
    suite.collection = suite.db.Collection("users")

    suite.repo = repository.NewUserRepository(*suite.db, "user")
}

func (suite *UserRepositorySuite) TearDownTest() {
    suite.db.Drop(context.TODO())
}

func (suite *UserRepositorySuite) TestCreateUser() {
	user1 := "Test User"
	
	email1 := "test@example.com"
    user := models.User{
        Name:  &user1,
        Email: &email1,
    }
    
    err := suite.repo.CreateUser(context.TODO(), user)
    suite.Require().NoError(err)

    // Fetch user to verify insertion
    var fetchedUser models.User
    suite.collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&fetchedUser)
    
    
}

func (suite *UserRepositorySuite) TestFetchAllUsers() {
    // Insert multiple users
	user1 := "User1"
	user2 :=  "User2"
	email1, email2 := "user1@example.com", "user2@example.com"
    users := []interface{}{
        models.User{Name: &user1, Email: &email1},
        models.User{Name: &user2, Email: &email2},
    }
    _, err := suite.collection.InsertMany(context.TODO(), users)
    suite.Require().NoError(err)

    suite.repo.FetchAllUsers(context.TODO())
    
    
}

func TestUserRepositorySuite(t *testing.T) {
    suite.Run(t, new(UserRepositorySuite))
}
