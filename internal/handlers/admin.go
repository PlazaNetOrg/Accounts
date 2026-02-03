package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"plazanet-accounts/internal/db"
	"plazanet-accounts/internal/models"
)

func ShowAdminDashboard(c *gin.Context) {
	var userCount int64
	db.DB.Model(&models.User{}).Count(&userCount)

	renderPage(c, "admin-dashboard", "page_title.admin", gin.H{
		"IsAuthenticated": true,
		"UserCount":       userCount,
	})
}

func ShowAdminUsers(c *gin.Context) {
	var users []models.User
	db.DB.Find(&users)

	renderPage(c, "admin-users", "page_title.admin_users", gin.H{
		"IsAuthenticated": true,
		"Users":           users,
	})
}

func ApiGetAllUsers(c *gin.Context) {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func ApiDeleteUser(c *gin.Context) {
	id := c.Param("id")
	
	if err := db.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func ApiUpdateUser(c *gin.Context) {
	id := c.Param("id")
	
	var input struct {
		Email       string `json:"email"`
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		IsAdmin     bool   `json:"is_admin"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Email = input.Email
	user.Username = input.Username
	user.DisplayName = input.DisplayName
	user.IsAdmin = input.IsAdmin

	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
