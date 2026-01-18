package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"plazanet-accounts/internal/db"
	"plazanet-accounts/internal/models"
)

type SetupDisplayNameInput struct {
	DisplayName string `form:"display_name" json:"display_name" binding:"required,min=1,max=50"`
}

type SetupPalCreatorInput struct {
	PalConfig string `form:"pal_config" json:"pal_config" binding:"required"`
}

func ApiSetupDisplayName(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input SetupDisplayNameInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
		return
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := db.DB.Model(&user).Updates(map[string]interface{}{
		"display_name":  input.DisplayName,
		"setup_status":  "display_name_set",
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update display name"})
		return
	}

	if c.GetHeader("HX-Request") != "" {
		c.Header("HX-Redirect", "/setup/recommendations")
		c.Status(http.StatusOK)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Display name set successfully",
		"display_name": input.DisplayName,
	})
}

func ApiSetupRecommendations(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := db.DB.Model(&user).Updates(map[string]interface{}{
		"setup_status": "completed",
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete setup"})
		return
	}

	if c.GetHeader("HX-Request") != "" {
		c.Header("HX-Redirect", "/dashboard")
		c.Status(http.StatusOK)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Setup complete!",
	})
}
