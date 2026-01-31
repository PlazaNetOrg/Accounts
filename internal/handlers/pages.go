package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"plazanet-accounts/internal/config"
	"plazanet-accounts/internal/db"
	"plazanet-accounts/internal/models"
)

func ShowLoginPage(c *gin.Context) {
	renderPage(c, "login", "page_title.login", gin.H{
		"IsAuthenticated": false,
	})
}

func ShowRegisterPage(c *gin.Context) {
	renderPage(c, "register", "page_title.register", gin.H{
		"IsAuthenticated": false,
	})
}

func ShowDashboard(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	renderPage(c, "dashboard", "page_title.dashboard", gin.H{
		"IsAuthenticated": true,
		"DisplayName":     user.DisplayName,
		"CurrentStatus":   user.CurrentStatus,
		"CurrentGame":     user.CurrentGame,
		"ClientType":      user.ClientType,
		"PlazaNetURL":     config.Cfg.PlazaNetURL,
	})
}

func ShowPalEditor(c *gin.Context) {
	renderPage(c, "pal-editor", "page_title.pal_editor", gin.H{
		"IsAuthenticated": true,
	})
}

func ShowSetupDisplayName(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	renderPage(c, "setup-display-name", "page_title.setup_display_name", gin.H{
		"IsAuthenticated": true,
		"Username":        user.Username,
	})
}

func ShowSetupRecommendations(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	renderPage(c, "setup-recommendations", "page_title.setup_recommendations", gin.H{
		"IsAuthenticated":  true,
		"DisplayName":      user.DisplayName,
		"PlazaNetName":     config.Cfg.PlazaNetName,
		"PlazaNetEnabled":  config.Cfg.PlazaNetEnabled,
		"PlazaNetURL":      config.Cfg.PlazaNetURL,
		"GamePlazaEnabled": config.Cfg.GamePlazaEnabled,
	})
}

