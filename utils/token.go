package Utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type SignedDetails struct{
	Email 	string
	Name 	string
	Uid 	string
	Role	string
	jwt.StandardClaims 
}

var jwt_secret_key string

func ImportJWTSecretKey() (string, error) {
    log.Println("Loading environment variables from .env file")
    err := godotenv.Load(".env")
    if err != nil {
        log.Println("Error loading .env file:", err)
        return "", errors.New("secret key not defined in env")
    }

    jwt_secret_key = os.Getenv("jwt_secret_key")
    log.Println("Retrieved JWT secret key from environment variable")

    if jwt_secret_key == "" {
        log.Println("JWT secret key is not defined in the environment")
        return "", errors.New("secret key not defined in env")
    }

    return jwt_secret_key, nil
}


func GenerateTokens(email string, name string, UID string, role string) (*string, *string, error) {

    if jwt_secret_key == "" {
        ImportJWTSecretKey()
    }

    claims := &SignedDetails{
        email, name, UID, role, jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 5).Unix(),
        },
    }

    refreshClaims := &SignedDetails{
        email, name, UID, role, jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 168).Unix(),
        },
    }

    log.Println("Generating token with claims:", claims)
    token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwt_secret_key))
    if err != nil {
        log.Println("Error generating token:", err)
        return nil, nil, errors.New("couldn't generate token")
    }
    log.Println("Generated token:", token)

    log.Println("Generating refresh token with claims:", refreshClaims)
    refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(jwt_secret_key))
    if err != nil {
        log.Println("Error generating refresh token:", err)
        return nil, nil, errors.New("couldn't generate refresh token")
    }
    log.Println("Generated refresh token:", refreshToken)

    return &token, &refreshToken, nil
}

func ValidateToken(signedToken string) (*SignedDetails , error){
	if jwt_secret_key == ""{
		ImportJWTSecretKey()
	}

	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func (token *jwt.Token) (interface{}, error ){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
			return []byte(jwt_secret_key), nil
		},
	)

	if err != nil{
		return nil, errors.New("wrong Credentails")
	}

	claims, ok := token.Claims.(*SignedDetails)
	
	if !ok {
		return nil, errors.New("wrong Credentails")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("expired Token")
	}


	return claims, nil

} 



func UpdateAllTokens(signedToken string, signedRefreshToken string) primitive.D{
	

	var updateObj primitive.D
	
	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})


	return updateObj
}


