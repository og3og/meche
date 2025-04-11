package config

import (
	"fmt"
	"meche/internal/handlers"
	customMiddleware "meche/internal/middleware"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

// NewServer creates and configures a new Echo server instance
func NewServer() *echo.Echo {
	e := echo.New()

	// Setup middleware
	customMiddleware.SetupMiddleware(e)

	// Initialize the session store
	key := []byte("your-secret-key") // Replace with a secure key in production
	store := sessions.NewCookieStore(key)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30, // 30 days
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteLaxMode,
	}
	gothic.Store = store

	// Debug: Print environment variables
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	fmt.Printf("GOOGLE_CLIENT_ID: %s\n", clientID)
	fmt.Printf("GOOGLE_CLIENT_SECRET: %s\n", clientSecret)

	// Initialize Google OAuth provider
	goth.UseProviders(
		google.New(
			clientID,
			clientSecret,
			"http://localhost:3000/auth/google/callback",
		),
	)

	// Serve static files
	e.Static("/static", "static")

	// Setup routes
	setupRoutes(e)

	return e
}

// setupRoutes configures all routes for the application
func setupRoutes(e *echo.Echo) {
	e.GET("/", handlers.HandleHome)
	e.GET("/login", handlers.HandleLogin)
	e.GET("/auth/google", handlers.HandleAuth)
	e.GET("/auth/google/callback", handlers.HandleCallback)
	e.GET("/logout", handlers.HandleLogout)

	// Protected routes
	dashboard := e.Group("/dashboard")
	dashboard.Use(customMiddleware.AuthMiddleware)
	dashboard.GET("", handlers.HandleDashboard)
}
