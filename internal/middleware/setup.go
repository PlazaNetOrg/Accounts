package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"plazanet-accounts/internal/db"
	"plazanet-accounts/internal/models"
)

func SetupRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		
		var user models.User
		if err := db.DB.First(&user, userID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		if c.Request.URL.Path == "/setup/display-name" || 
	   c.Request.URL.Path == "/setup/recommendations" ||
	   c.Request.URL.Path == "/api/setup/display-name" ||
	   c.Request.URL.Path == "/api/setup/recommendations" {
		c.Set("setup_status", user.SetupStatus)
		c.Next()
		return
	}

	if user.SetupStatus != "completed" {
		var redirect string
		switch user.SetupStatus {
		case "not_started":
			redirect = "/setup/display-name"
		case "display_name_set":
			redirect = "/setup/recommendations"
			}

			log.Printf("[SETUP] User %d (status: %s) redirected to %s", userID, user.SetupStatus, redirect)

			if c.GetHeader("HX-Request") != "" {
				c.Header("HX-Redirect", redirect)
				c.Status(http.StatusOK)
				c.Abort()
				return
			}

			c.Redirect(http.StatusFound, redirect)
			c.Abort()
			return
		}

		c.Set("setup_status", user.SetupStatus)
		c.Next()
	}
}
