package models_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/sofc-t/task_manager/task8/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserInitialization(t *testing.T) {
	id := primitive.NewObjectID()
	name := "Uncle Bob"
	password := "securepassword"
	email := "uncle.bob@example.com"
	token := "some-token"
	role := "admin"
	refreshToken := "refresh-token"
	createdAt := time.Now()
	updatedAt := time.Now()
	userID := "user-id"

	user := models.User{
		ID:             id,
		Name:           &name,
		Password:       &password,
		Email:          &email,
		Token:          &token,
		Role:           &role,
		RefreshToken:   &refreshToken,
		CreatedAt:      createdAt,
		UpdatedAT:      updatedAt,
		UserID:         userID,
	}

	assert.Equal(t, id, user.ID, "User ID should match")
	assert.Equal(t, name, *user.Name, "User Name should match")
	assert.Equal(t, password, *user.Password, "User Password should match")
	assert.Equal(t, email, *user.Email, "User Email should match")
	assert.Equal(t, token, *user.Token, "User Token should match")
	assert.Equal(t, role, *user.Role, "User Role should match")
	assert.Equal(t, refreshToken, *user.RefreshToken, "User RefreshToken should match")
	assert.WithinDuration(t, createdAt, user.CreatedAt, time.Second, "User CreatedAt should be close to the initialized value")
	assert.WithinDuration(t, updatedAt, user.UpdatedAT, time.Second, "User UpdatedAT should be close to the initialized value")
	assert.Equal(t, userID, user.UserID, "User UserID should match")
}

func TestUserJSONSerialization(t *testing.T) {
	id := primitive.NewObjectID()
	name := "Uncle Bob"
	password := "securepassword"
	email := "uncle.bob@example.com"
	token := "some-token"
	role := "admin"
	refreshToken := "refresh-token"
	createdAt := time.Now()
	updatedAt := time.Now()
	userID := "user-id"

	user := models.User{
		ID:             id,
		Name:           &name,
		Password:       &password,
		Email:          &email,
		Token:          &token,
		Role:           &role,
		RefreshToken:   &refreshToken,
		CreatedAt:      createdAt,
		UpdatedAT:      updatedAt,
		UserID:         userID,
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Error marshalling user to JSON: %v", err)
	}

	expectedJSON := `{
		"_id":"` + id.Hex() + `",
		"name":"Uncle Bob",
		"password":"securepassword",
		"email":"uncle.bob@example.com",
		"token":"some-token",
		"role":"admin",
		"refresh_token":"refresh-token",
		"created_at":"` + createdAt.Format(time.RFC3339Nano) + `",
		"updated_at":"` + updatedAt.Format(time.RFC3339Nano) + `",
		"user_id":"user-id"
	}`

	assert.JSONEq(t, expectedJSON, string(userJSON), "JSON output should match expected output")
}
