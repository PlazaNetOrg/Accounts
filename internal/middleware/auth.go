package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"plazanet-accounts/internal/db"
	"plazanet-accounts/internal/models"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == authHeader {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
				return
			}
		} else {
			var err error
			tokenStr, err = c.Cookie("auth_token")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
				return
			}
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in token"})
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username in token"})
			return
		}

		var user models.User
		if err := db.DB.Select("language").Where("id = ?", uint(userID)).First(&user).Error; err != nil {
			log.Printf("[AUTH] Failed to fetch user language: %v", err)
			user.Language = "en"
		}

		c.Set("user_id", uint(userID))
		c.Set("username", username)
		c.Set("language", user.Language)
		log.Printf("[AUTH] Set context - user_id: %d, username: %s, language: %s", uint(userID), username, user.Language)
		c.Next()
	}
}

func SetDefaultLanguage() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("language"); !exists {
			c.Set("language", "en")
		}
		c.Next()
	}
}