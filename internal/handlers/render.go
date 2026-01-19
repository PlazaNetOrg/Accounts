package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"plazanet-accounts/internal/config"
)

func renderPage(c *gin.Context, page string, title string, data gin.H) {
	if data == nil {
		data = gin.H{}
	}

	lang := c.GetString("language")
	if lang == "" {
		lang = "en"
	}

	data["Lang"] = lang
	data["Title"] = title
	data["Page"] = page
	data["ServerName"] = config.Cfg.ServerName
	data["PlazaNetName"] = config.Cfg.PlazaNetName

	c.HTML(http.StatusOK, "base.html", data)
}

