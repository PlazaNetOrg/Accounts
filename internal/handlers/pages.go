package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"plazanet-accounts/internal/config"
	"plazanet-accounts/internal/db"
	"plazanet-accounts/internal/models"
)

func ShowLoginPage(c *gin.Context) {
	renderPage(c, "login", "Login", gin.H{
		"IsAuthenticated": false,
	})
}

func ShowRegisterPage(c *gin.Context) {
	renderPage(c, "register", "Register", gin.H{
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

	renderPage(c, "dashboard", "Dashboard", gin.H{
		"IsAuthenticated": true,
		"DisplayName":     user.DisplayName,
		"PlazaNetURL":     config.Cfg.PlazaNetURL,
	})
}

func ShowPalEditor(c *gin.Context) {
	renderPage(c, "pal-editor", "Pal Editor", gin.H{
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

	renderPage(c, "setup-display-name", "Set Display Name", gin.H{
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

	renderPage(c, "setup-recommendations", "What's Next?", gin.H{
		"IsAuthenticated":  true,
		"DisplayName":      user.DisplayName,
		"PlazaNetName":     config.Cfg.PlazaNetName,
		"PlazaNetEnabled":  config.Cfg.PlazaNetEnabled,
		"PlazaNetURL":      config.Cfg.PlazaNetURL,
		"GamePlazaEnabled": config.Cfg.GamePlazaEnabled,
	})
}

