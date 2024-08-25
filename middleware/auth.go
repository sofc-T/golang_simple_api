package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	Utils "github.com/sofc-t/task_manager/task8/utils"
)

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			log.Println("No Authorization header provided")
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization Required"})
			ctx.Abort()
			return
		}

		parts := strings.Split(tokenHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			log.Println("Authorization header format must be Bearer <token>")
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format"})
			ctx.Abort()
			return
		}

		token := parts[1]
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization Required"})
			ctx.Abort()
			return
		}
		claims, err := Utils.ValidateToken(token)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "iNVALID tOKEN"})
			ctx.Abort()
			return
		}

		ctx.Set("email", claims.Email)
		ctx.Set("first_name", claims.Name)
		ctx.Set("uid", claims.Uid)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}

}

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			log.Println("No Authorization header provided")
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization Required"})
			ctx.Abort()
			return
		}

		parts := strings.Split(tokenHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			log.Println("Authorization header format must be Bearer <token>")
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format"})
			ctx.Abort()
			return
		}

		token := parts[1]
		log.Println("Token received:", token)

		claims, err := Utils.ValidateToken(token)
		if err != nil {
			log.Println("Token validation failed:", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Token"})
			ctx.Abort()
			return
		}

		if claims.Role != "admin" {
			log.Println("Forbidden access attempted by non-admin user:", claims.Email)
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Forbidden Content"})
			ctx.Abort()
			return
		}

		log.Println("Admin access granted for user:", claims.Email)
		ctx.Next()
	}
}

func AuthenticationandAuthorizeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Authorization header missing")
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization Required"})
			ctx.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			log.Println("Invalid or missing Bearer token")
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header"})
			ctx.Abort()
			return
		}

		claims, err := Utils.ValidateToken(token)
		if err != nil {
			log.Println("Token validation failed:", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Token"})
			ctx.Abort()
			return
		}

		if claims.Role != "admin" {
			log.Println("Unauthorized role:", claims.Role)
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Forbidden Content"})
			ctx.Abort()
			return
		}

		log.Println("Authenticated user:", claims.Email, "Role:", claims.Role)

		ctx.Set("email", claims.Email)
		ctx.Set("first_name", claims.Name)
		ctx.Set("uid", claims.Uid)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}
