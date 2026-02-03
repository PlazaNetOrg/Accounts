package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"plazanet-accounts/internal/db"
	"plazanet-accounts/internal/i18n"
	"plazanet-accounts/internal/models"
)

const (
	TokenExpirationTime = 24 * time.Hour
	CookieMaxAge        = 86400
	CookieName          = "auth_token"
)

func renderAuthError(c *gin.Context, statusCode int, message string) {
	errorHTML := fmt.Sprintf(`<div class="p-4 bg-red-900 border border-red-700 rounded-lg text-red-100">%s</div>`, message)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(statusCode, errorHTML)
}

type RegisterInput struct {
	Email    string `form:"email" binding:"required,email"`
	Username string `form:"username" binding:"required,min=3,max=20"`
	Password string `form:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func ApiRegister(c *gin.Context) {
	var input RegisterInput
	lang := c.GetString("language")
	if lang == "" {
		lang = "en"
	}

	if err := c.ShouldBind(&input); err != nil {
		renderAuthError(c, http.StatusBadRequest, i18n.Get(lang, "auth.error_validation"))
		return
	}

	// Convert username to lowercase
	username := strings.ToLower(input.Username)
	email := strings.ToLower(input.Email)

	var existing models.User
	if err := db.DB.Where("username = ? OR email = ?", username, email).First(&existing).Error; err == nil {
		renderAuthError(c, http.StatusConflict, i18n.Get(lang, "auth.error_username_taken"))
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		renderAuthError(c, http.StatusInternalServerError, i18n.Get(lang, "auth.error_account_creation"))
		return
	}

	user := models.User{
		Email:       email,
		Username:    username,
		Password:    string(hashed),
		SetupStatus: "not_started",
	}

	if err := db.DB.Create(&user).Error; err != nil {
		renderAuthError(c, http.StatusInternalServerError, i18n.Get(lang, "auth.error_account_creation"))
		return
	}

	tokenString, err := generateToken(user.ID, user.Username)
	if err != nil {
		renderAuthError(c, http.StatusInternalServerError, i18n.Get(lang, "auth.error_account_creation"))
		return
	}

	setAuthCookie(c, tokenString)

	if c.GetHeader("HX-Request") != "" || strings.Contains(c.GetHeader("Accept"), "text/html") {
		// Check if there's a return_to cookie (from PlazaNet or other services)
		returnTo, err := c.Cookie("return_to")
		if err == nil && returnTo != "" {
			c.SetCookie("return_to", "", -1, "/", "", false, true)
			// Append the token as a URL parameter so PlazaNet can set its own cookie
			redirectURL := returnTo
			if strings.Contains(redirectURL, "?") {
				redirectURL += "&token=" + tokenString
			} else {
				redirectURL += "?token=" + tokenString
			}
			c.Header("HX-Redirect", redirectURL)
			c.Status(http.StatusOK)
			return
		}

		c.Header("HX-Redirect", "/setup/display-name")
		c.Status(http.StatusOK)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "User created successfully",
		"username": user.Username,
		"token":    tokenString,
	})
}

func ApiLogin(c *gin.Context) {
	var input LoginInput
	lang := c.GetString("language")
	if lang == "" {
		lang = "en"
	}

	wantsJSON := strings.Contains(c.GetHeader("Accept"), "application/json") || c.GetHeader("Accept") == "application/json"

	if err := c.ShouldBind(&input); err != nil {
		if wantsJSON {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		} else {
			renderAuthError(c, http.StatusBadRequest, i18n.Get(lang, "auth.error_validation"))
		}
		return
	}

	// Convert username to lowercase
	username := strings.ToLower(input.Username)

	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if wantsJSON {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			renderAuthError(c, http.StatusUnauthorized, i18n.Get(lang, "auth.error_invalid_credentials"))
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		if wantsJSON {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			renderAuthError(c, http.StatusUnauthorized, i18n.Get(lang, "auth.error_invalid_credentials"))
		}
		return
	}

	tokenString, err := generateToken(user.ID, user.Username)
	if err != nil {
		if wantsJSON {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		} else {
			renderAuthError(c, http.StatusInternalServerError, i18n.Get(lang, "auth.error_login_failed"))
		}
		return
	}

	setAuthCookie(c, tokenString)

	if wantsJSON {
		c.JSON(http.StatusOK, gin.H{
			"token":    tokenString,
			"username": user.Username,
		})
		return
	}

	// Check if there's a return_to cookie (from PlazaNet or other services)
	returnTo, err := c.Cookie("return_to")
	if err == nil && returnTo != "" {
		c.SetCookie("return_to", "", -1, "/", "", false, true)
		redirectURL := returnTo
		if strings.Contains(redirectURL, "?") {
			redirectURL += "&token=" + tokenString
		} else {
			redirectURL += "?token=" + tokenString
		}
		c.Header("HX-Redirect", redirectURL)
		c.Status(http.StatusOK)
		return
	}

	var setupRedirect string
	switch user.SetupStatus {
	case "not_started":
		setupRedirect = "/setup/display-name"
	case "display_name_set":
		setupRedirect = "/setup/recommendations"
	case "completed":
		setupRedirect = "/dashboard"
	default:
		setupRedirect = "/setup/display-name"
	}

	c.Header("HX-Redirect", setupRedirect)
	c.Status(http.StatusOK)
}

func Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/login")
}

func Me(c *gin.Context) {
	userID := c.GetUint("user_id")
	username := c.GetString("username")

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           userID,
		"username":     username,
		"display_name": user.DisplayName,
	})
}

func setAuthCookie(c *gin.Context, token string) {
	domain := ""
	
	c.SetCookie(
		CookieName,
		token,
		CookieMaxAge,
		"/",
		domain,
		false,
		false,
	)
}

func generateToken(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(TokenExpirationTime).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}