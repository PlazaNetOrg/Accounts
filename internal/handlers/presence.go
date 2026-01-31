package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"plazanet-accounts/internal/db"
	"plazanet-accounts/internal/models"
)

type UpdatePresenceInput struct {
	Status     string `json:"status" binding:"required,oneof=offline online playing"`
	Game       string `json:"game"`
	ClientType string `json:"client_type" binding:"required,oneof=web gameplaza mobile"`
}

type PresenceHeartbeatInput struct {
	Game       string `json:"game"`
	ClientType string `json:"client_type" binding:"required,oneof=web gameplaza mobile"`
}

type UserPresenceResponse struct {
	UserID        uint       `json:"user_id"`
	Username      string     `json:"username"`
	DisplayName   string     `json:"display_name"`
	CurrentStatus string     `json:"current_status"`
	CurrentGame   string     `json:"current_game,omitempty"`
	ClientType    string     `json:"client_type"`
	LastSeenAt    *time.Time `json:"last_seen_at,omitempty"`
}

func clientTypePriority(clientType, status string) int {
	if clientType == "gameplaza" {
		if status == "playing" {
			return 5
		}
		return 1
	}

	priorities := map[string]int{"mobile": 4, "web": 3, "unknown": 0}
	if p, ok := priorities[clientType]; ok {
		return p
	}
	return 0
}

func shouldUpdateClientType(incomingType, incomingStatus, currentType, currentStatus string) bool {
	if currentType == "" {
		return true
	}
	return clientTypePriority(incomingType, incomingStatus) >= clientTypePriority(currentType, currentStatus)
}

func buildPresenceUpdates(status, game, clientType string) map[string]interface{} {
	updates := map[string]interface{}{
		"current_status": status,
		"client_type":    clientType,
		"last_seen_at":   time.Now(),
	}
	if status == "playing" && game != "" {
		updates["current_game"] = game
	} else {
		updates["current_game"] = ""
	}
	return updates
}

func resolveHeartbeatState(input PresenceHeartbeatInput, user models.User) (string, string, string) {
	status := "online"
	game := ""

	if input.Game != "" {
		status = "playing"
		game = input.Game
	}

	clientType := input.ClientType
	if !shouldUpdateClientType(input.ClientType, status, user.ClientType, user.CurrentStatus) {
		clientType = user.ClientType
	}

	return status, game, clientType
}

func ApiUpdatePresence(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input UpdatePresenceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if input.Status == "playing" && input.Game == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game name is required when status is 'playing'"})
		return
	}

	var user models.User
	db.DB.Select("client_type, current_status").First(&user, userID)

	clientType := input.ClientType
	if !shouldUpdateClientType(input.ClientType, input.Status, user.ClientType, user.CurrentStatus) {
		clientType = user.ClientType
	}

	updates := buildPresenceUpdates(input.Status, input.Game, clientType)
	if err := db.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update presence"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Presence updated successfully",
		"current_status": input.Status,
		"current_game":   updates["current_game"],
		"client_type":    clientType,
	})
}

func ApiPresenceHeartbeat(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input PresenceHeartbeatInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var user models.User
	db.DB.Select("client_type, current_game, current_status").First(&user, userID)

	status, game, clientType := resolveHeartbeatState(input, user)
	updates := buildPresenceUpdates(status, game, clientType)

	if err := db.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update heartbeat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Heartbeat received",
		"last_seen":   updates["last_seen_at"],
		"client_type": clientType,
	})
}

func ApiGetUserPresence(c *gin.Context) {
	requestingUserID := c.GetUint("user_id")
	targetUsername := c.Param("username")

	var targetUser models.User
	if err := db.DB.Where("username = ?", targetUsername).First(&targetUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !canUserViewPresence(requestingUserID, targetUser) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to view this user's status"})
		return
	}

	c.JSON(http.StatusOK, buildPresenceResponse(targetUser))
}

func ApiGetMyPresence(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := buildPresenceResponse(user)
	response.LastSeenAt = user.LastSeenAt
	c.JSON(http.StatusOK, response)
}

func buildPresenceResponse(user models.User) UserPresenceResponse {
	response := UserPresenceResponse{
		UserID:        user.ID,
		Username:      user.Username,
		DisplayName:   user.DisplayName,
		CurrentStatus: user.CurrentStatus,
		ClientType:    user.ClientType,
	}

	if user.CurrentStatus == "playing" {
		response.CurrentGame = user.CurrentGame
	}

	if user.CurrentStatus == "offline" && user.LastSeenAt != nil {
		response.LastSeenAt = user.LastSeenAt
	}

	return response
}

func canUserViewPresence(requestingUserID uint, targetUser models.User) bool {
	if requestingUserID == targetUser.ID {
		return true
	}

	switch targetUser.StatusPrivacy {
	case "everyone", "friends":
		return true
	case "private":
		return false
	default:
		return false
	}
}
