package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"plazanet-accounts/internal/db"
	"plazanet-accounts/internal/i18n"
	"plazanet-accounts/internal/models"
)

func renderSettingsError(c *gin.Context, message string) {
	errorHTML := fmt.Sprintf(`<div class="p-4 bg-red-900 border border-red-700 rounded-lg text-red-100">%s</div>`, message)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusBadRequest, errorHTML)
}

func renderSettingsSuccess(c *gin.Context, message string) {
	successHTML := fmt.Sprintf(`<div class="p-4 bg-green-900 border border-green-700 rounded-lg text-green-100">%s</div>`, message)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, successHTML)
}

func ShowSettings(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	renderPage(c, "settings", "page_title.settings", gin.H{
		"IsAuthenticated":    true,
		"CurrentDisplayName": user.DisplayName,
		"CurrentLanguage":    user.Language,
		"CurrentStatusPrivacy": user.StatusPrivacy,
	})
}

type UpdateDisplayNameInput struct {
	DisplayName string `form:"display_name" binding:"required,max=50"`
}

func ApiUpdateDisplayName(c *gin.Context) {
	userID := c.GetUint("user_id")
	lang := c.GetString("language")
	if lang == "" {
		lang = "en"
	}

	var input UpdateDisplayNameInput
	if err := c.ShouldBind(&input); err != nil {
		renderSettingsError(c, i18n.Get(lang, "settings.error_generic"))
		return
	}

	if err := db.DB.Model(&models.User{}).Where("id = ?", userID).Update("display_name", input.DisplayName).Error; err != nil {
		renderSettingsError(c, i18n.Get(lang, "settings.error_generic"))
		return
	}

	renderSettingsSuccess(c, i18n.Get(lang, "settings.success_display_name"))
}

type UpdateLanguageInput struct {
	Language string `form:"language" binding:"required,oneof=en pl"`
}

func ApiUpdateLanguage(c *gin.Context) {
	userID := c.GetUint("user_id")
	lang := c.GetString("language")
	if lang == "" {
		lang = "en"
	}

	var input UpdateLanguageInput
	if err := c.ShouldBind(&input); err != nil {
		renderSettingsError(c, i18n.Get(lang, "settings.error_generic"))
		return
	}

	if err := db.DB.Model(&models.User{}).Where("id = ?", userID).Update("language", input.Language).Error; err != nil {
		renderSettingsError(c, i18n.Get(lang, "settings.error_generic"))
		return
	}

	renderSettingsSuccess(c, i18n.Get(lang, "settings.success_language"))
}

type UpdateStatusPrivacyInput struct {
	StatusPrivacy string `form:"status_privacy" binding:"required,oneof=everyone friends private"`
}

func ApiUpdateStatusPrivacy(c *gin.Context) {
	userID := c.GetUint("user_id")
	lang := c.GetString("language")
	if lang == "" {
		lang = "en"
	}

	var input UpdateStatusPrivacyInput
	if err := c.ShouldBind(&input); err != nil {
		renderSettingsError(c, i18n.Get(lang, "settings.error_generic"))
		return
	}

	if err := db.DB.Model(&models.User{}).Where("id = ?", userID).Update("status_privacy", input.StatusPrivacy).Error; err != nil {
		renderSettingsError(c, i18n.Get(lang, "settings.error_generic"))
		return
	}

	renderSettingsSuccess(c, i18n.Get(lang, "settings.success_status_privacy"))
}

type UpdatePasswordInput struct {
	CurrentPassword string `form:"current_password" binding:"required"`
	NewPassword     string `form:"new_password" binding:"required,min=6"`
	ConfirmPassword string `form:"confirm_password" binding:"required,min=6"`
}

func ApiUpdatePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	lang := c.GetString("language")
	if lang == "" {
		lang = "en"
	}

	var input UpdatePasswordInput
	if err := c.ShouldBind(&input); err != nil {
		renderSettingsError(c, i18n.Get(lang, "settings.error_generic"))
		return
	}

	if input.NewPassword != input.ConfirmPassword {
		renderSettingsError(c, i18n.Get(lang, "settings.error_password_mismatch"))
		return
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		renderSettingsError(c, i18n.Get(lang, "settings.error_generic"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		renderSettingsError(c, i18n.Get(lang, "settings.error_invalid_password"))
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		renderSettingsError(c, i18n.Get(lang, "settings.error_generic"))
		return
	}

	if err := db.DB.Model(&user).Update("password", string(hashed)).Error; err != nil {
		renderSettingsError(c, i18n.Get(lang, "settings.error_generic"))
		return
	}

	renderSettingsSuccess(c, i18n.Get(lang, "settings.success_password"))
}
