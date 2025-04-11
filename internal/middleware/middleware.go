package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth/gothic"
)

// SetupMiddleware configures all middleware for the application
func SetupMiddleware(e *echo.Echo) {
	// Basic middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// Security middleware
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self' 'unsafe-inline' 'unsafe-eval' https: data:",
	}))

	// Request timeout middleware
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	}))

	// Request ID middleware
	e.Use(middleware.RequestID())

	// Gzip compression
	e.Use(middleware.Gzip())
}

// AuthMiddleware checks if a user is authenticated
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := gothic.Store.Get(c.Request(), "gothic-session")
		if err != nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		// Check if user is in session
		if _, ok := session.Values["user"]; !ok {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		return next(c)
	}
}
