package handlers

import (
	"log"
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
	Status     string `json:"status" binding:"omitempty,oneof=online playing"`
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

const clientStalenessThreshold = 40 * time.Second
const offlineThreshold = 60 * time.Second

func clientTypePriority(clientType, status string) int {
	if clientType == "gameplaza" {
		if status == "playing" {
			return 5
		}
		return 2
	}

	priorities := map[string]int{"mobile": 4, "web": 3, "unknown": 0}
	if p, ok := priorities[clientType]; ok {
		return p
	}
	return 0
}

func shouldUpdateClientType(incomingType, incomingStatus, currentType, currentStatus string, lastClientUpdateAt *time.Time) bool {
	if currentType == "" {
		return true
	}
	
	// Check if current client is stale (hasn't sent heartbeat in 40+ seconds)
	// If lastClientUpdateAt is nil, the current state is stale (allow takeover)
	if lastClientUpdateAt != nil {
		timeSince := time.Since(*lastClientUpdateAt)
		log.Printf("[STALENESS] Current client %s/%s - Time since last update: %.1fs (threshold: %.1fs)", 
			currentType, currentStatus, timeSince.Seconds(), clientStalenessThreshold.Seconds())
		if timeSince > clientStalenessThreshold {
			log.Printf("[STALENESS] Current client stale - allowing %s/%s to take over", incomingType, incomingStatus)
			return true
		}
	} else {
		log.Printf("[STALENESS] No last update timestamp for %s/%s - allowing %s/%s to take over", 
			currentType, currentStatus, incomingType, incomingStatus)
		return true
	}
	
	incomingPriority := clientTypePriority(incomingType, incomingStatus)
	currentPriority := clientTypePriority(currentType, currentStatus)
	
	// Allow same client to refresh, or higher priority to take over
	result := incomingPriority >= currentPriority
	log.Printf("[PRIORITY] %s/%s (p%d) vs %s/%s (p%d) -> %v", 
		incomingType, incomingStatus, incomingPriority,
		currentType, currentStatus, currentPriority, result)
	
	return result
}

func buildPresenceUpdates(status, game, clientType string, wasActuallyUpdated bool) map[string]interface{} {
	now := time.Now()
	updates := map[string]interface{}{
		"current_status": status,
		"client_type":    clientType,
		"last_seen_at":   now,
	}
	
	if wasActuallyUpdated {
		updates["last_client_update_at"] = now
	}
	
	if status == "playing" && game != "" {
		updates["current_game"] = game
	} else {
		updates["current_game"] = ""
	}
	return updates
}

func resolveHeartbeatState(input PresenceHeartbeatInput, user models.User) (string, string, string, bool) {
	// Check if this client should be allowed to update at all
	incomingStatus := "online"
	if input.Status != "" {
		incomingStatus = input.Status
	} else if input.Game != "" {
		incomingStatus = "playing"
	}
	
	if !shouldUpdateClientType(input.ClientType, incomingStatus, user.ClientType, user.CurrentStatus, user.LastClientUpdateAt) {
		return user.CurrentStatus, user.CurrentGame, user.ClientType, false
	}
	
	status := "online"
	game := ""

	if input.Status != "" {
		status = input.Status
		if status == "playing" {
			if input.Game != "" {
				game = input.Game
			} else if user.CurrentGame != "" {
				game = user.CurrentGame
			}
		}
	} else if input.Game != "" {
		status = "playing"
		game = input.Game
	}

	return status, game, input.ClientType, true
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
	db.DB.Select("client_type, current_status, last_client_update_at").First(&user, userID)

	clientType := input.ClientType
	wasActuallyUpdated := true
	if !shouldUpdateClientType(input.ClientType, input.Status, user.ClientType, user.CurrentStatus, user.LastClientUpdateAt) {
		clientType = user.ClientType
		wasActuallyUpdated = false
	}

	updates := buildPresenceUpdates(input.Status, input.Game, clientType, wasActuallyUpdated)
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
	db.DB.Select("client_type, current_game, current_status, last_client_update_at").First(&user, userID)

	status, game, clientType, wasActuallyUpdated := resolveHeartbeatState(input, user)
	updates := buildPresenceUpdates(status, game, clientType, wasActuallyUpdated)

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
	currentStatus := user.CurrentStatus
	clientType := user.ClientType
	
	// Mark user as offline if they haven't sent any heartbeat in 60+ seconds
	if user.LastSeenAt != nil && time.Since(*user.LastSeenAt) > offlineThreshold {
		currentStatus = "offline"
	}
	
	response := UserPresenceResponse{
		UserID:        user.ID,
		Username:      user.Username,
		DisplayName:   user.DisplayName,
		CurrentStatus: currentStatus,
		ClientType:    clientType,
	}

	if currentStatus == "playing" {
		response.CurrentGame = user.CurrentGame
	}

	if currentStatus == "offline" && user.LastSeenAt != nil {
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
