package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/k1ngalph0x/payflow/auth-service/api"
	"github.com/k1ngalph0x/payflow/auth-service/config"
)

type AuthMiddleware struct {
	Config *config.Config
}


func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{Config: cfg}
}


func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context){
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == ""{
			c.JSON(http.StatusUnauthorized, gin.H{"error":"Unauthorized: Missing or invalid token format"})
			c.Abort()
			return 
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer"{
			c.JSON(http.StatusUnauthorized, gin.H{"error":"Unauthorized: Missing or invalid token format"})
			c.Abort()
			return 
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(tokenString, &api.Claims{}, func(token *jwt.Token)(interface{}, error){

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
				return nil, fmt.Errorf("unexpected signing method")
			}

			return []byte(a.Config.TOKEN.JwtKey), nil
		})

		if err !=nil || !token.Valid{
			c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid or expired token"})
			c.Abort()
			return 
		}

		claims, ok := token.Claims.(*api.Claims)
		if !ok{
			c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid token claims"})
			c.Abort()
			return 
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
		
	}
}

func (a *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc{
	return func(c *gin.Context){
	
		userRole, exists := c.Get("role")
		if !exists{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No role found"})
			c.Abort()
			return
		}

		roleString, ok := userRole.(string)

		if !ok{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No role found"})
			c.Abort()
			return
		}
		allowed:= false
		for _, allowedRole := range roles{
			if roleString == allowedRole{
				allowed = true
				break
			}
		}

		if !allowed{
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No role found"})
			c.Abort()
			return
		}

		c.Next()

	}
}