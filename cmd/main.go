package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"os"

	"github.com/gin-gonic/gin"

	"plazanet-accounts/internal/config"
	"plazanet-accounts/internal/db"
	"plazanet-accounts/internal/handlers"
	"plazanet-accounts/internal/i18n"
	"plazanet-accounts/internal/middleware"
)

func main() {
	cfg := config.Load()
	i18n.Init()

	db.Connect()

	r := gin.Default()
	files, err := filepath.Glob("templates/*.html")
	if err != nil {
		log.Fatalf("Cannot glob templates: %v", err)
	}
	if len(files) == 0 {
		log.Fatal("ERROR: No template files found in templates/*.html !")
	}
	log.Printf("Loaded templates: %v", files)

	tmpl := template.New("").Funcs(template.FuncMap{
		"T": func(lang, id string) string {
			return i18n.Get(lang, id)
		},
	})
	tmpl = template.Must(tmpl.ParseGlob("templates/*.html"))
	r.SetHTMLTemplate(tmpl)
	r.Static("/static", "./static")

	r.Use(middleware.SetDefaultLanguage())

	// ====================
    //  UI routes
    // ====================
	r.StaticFile("/favicon.ico", "./favicon.ico")

    r.GET("/login", handlers.ShowLoginPage)
    r.GET("/register", handlers.ShowRegisterPage)

    // Protected UI routes
    protectedUI := r.Group("/")
    protectedUI.Use(middleware.AuthRequiredUI())
    {
        protectedUI.GET("/", func(c *gin.Context) {
            c.Redirect(http.StatusFound, "/dashboard")
        })
        protectedUI.GET("/dashboard", middleware.SetupRequired(), handlers.ShowDashboard)
        protectedUI.GET("/pal-editor", middleware.SetupRequired(), handlers.ShowPalEditor)
        protectedUI.GET("/logout", handlers.Logout)
        
        // Setup routes
        protectedUI.GET("/setup/display-name", handlers.ShowSetupDisplayName)
        protectedUI.GET("/setup/recommendations", handlers.ShowSetupRecommendations)
    }
	
	// ====================
    //  API routes (/api/...)
    // ====================
    api := r.Group("/api")
    {
        // Public API
        api.POST("/register", handlers.ApiRegister)
        api.POST("/login", handlers.ApiLogin)

        // Protected API
        protectedAPI := api.Group("/")
        protectedAPI.Use(middleware.AuthRequired())
        {
            protectedAPI.GET("/me", handlers.Me)
            
            // Setup API routes
            protectedAPI.POST("/setup/display-name", handlers.ApiSetupDisplayName)
            protectedAPI.POST("/setup/recommendations", handlers.ApiSetupRecommendations)
        }
    }

	port := cfg.Port
	serverName := cfg.ServerName

	log.Printf("Starting %s on %s", serverName, port)

	if err := r.Run(port); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}