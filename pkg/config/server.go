package config

import (
	"fmt"
	"meche/internal/handlers"
	customMiddleware "meche/internal/middleware"
	orgHandlers "meche/pkg/handlers"
	"meche/pkg/storage"
	"meche/pkg/storage/memory"
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

	// Add validator middleware
	e.Validator = &CustomValidator{}

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

	// Initialize storage
	orgStorage := memory.NewMemoryOrganizationStorage()
	memberStorage := memory.NewMemoryOrganizationMemberStorage()
	projectStorage := memory.NewMemoryProjectStorage()
	taskStorage := memory.NewMemoryTaskStorage()

	// Initialize seed data
	memory.SeedData(orgStorage, projectStorage)

	// Serve static files
	e.Static("/static", "static")

	// Setup routes
	setupRoutes(e, orgStorage, memberStorage, projectStorage, taskStorage)

	return e
}

// CustomValidator implements echo.Validator interface
type CustomValidator struct{}

func (cv *CustomValidator) Validate(i interface{}) error {
	return nil // We're using struct tags for validation
}

// setupRoutes configures all routes for the application
func setupRoutes(e *echo.Echo, orgStorage storage.OrganizationStorage, memberStorage storage.OrganizationMemberStorage, projectStorage storage.ProjectStorage, taskStorage storage.TaskStorage) {
	// Auth routes
	e.GET("/", handlers.HandleHome)
	e.GET("/login", handlers.HandleLogin)
	e.GET("/auth/google", handlers.HandleAuth)
	e.GET("/auth/google/callback", handlers.HandleCallback)
	e.GET("/logout", handlers.HandleLogout)

	// Protected routes
	protected := e.Group("")
	protected.Use(customMiddleware.AuthMiddleware)
	protected.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, err := gothic.Store.Get(c.Request(), "gothic-session")
			if err != nil {
				return c.String(http.StatusUnauthorized, "User not authenticated")
			}

			user, ok := session.Values["user"].(goth.User)
			if !ok {
				return c.String(http.StatusUnauthorized, "User not authenticated")
			}

			// Add storage instances to context
			c.Set("organization_storage", orgStorage)
			c.Set("member_storage", memberStorage)
			c.Set("project_storage", projectStorage)
			c.Set("task_storage", taskStorage)
			c.Set("user", user)

			return next(c)
		}
	})

	// Dashboard route
	protected.GET("/dashboard", handlers.HandleDashboard)

	// Organization routes
	protected.POST("/organizations", orgHandlers.CreateOrganization(orgStorage, memberStorage))
	protected.GET("/organizations", orgHandlers.ListOrganizations(orgStorage))
	protected.GET("/organizations/new", orgHandlers.NewOrganizationForm)
	protected.GET("/organizations/cancel", orgHandlers.CancelOrganizationForm)
	protected.DELETE("/organizations/:id", orgHandlers.DeleteOrganization(orgStorage))
	protected.GET("/organizations/:id/edit", orgHandlers.EditOrganizationForm(orgStorage))
	protected.PUT("/organizations/:id", orgHandlers.UpdateOrganization(orgStorage))
	protected.GET("/organizations/:id", orgHandlers.ShowOrganization(orgStorage, projectStorage))
	protected.GET("/organizations/:id/overview", orgHandlers.ShowOrganizationOverview(orgStorage, projectStorage))
	protected.GET("/organizations/:id/settings", orgHandlers.ShowOrganizationSettings(orgStorage, projectStorage))

	// Project routes
	protected.POST("/organizations/:orgID/projects", orgHandlers.CreateProject(projectStorage))
	protected.GET("/organizations/:orgID/projects", orgHandlers.ListProjects(projectStorage))
	protected.GET("/organizations/:orgID/projects/new", orgHandlers.NewProjectForm(orgStorage, projectStorage))
	protected.GET("/organizations/:orgID/projects/cancel", orgHandlers.CancelProjectForm)
	protected.GET("/organizations/:orgID/projects/:id", orgHandlers.ShowProject(projectStorage, orgStorage))
	protected.GET("/organizations/:orgID/projects/:id/overview", orgHandlers.ShowProjectOverview(projectStorage))
	protected.GET("/organizations/:orgID/projects/:id/settings", orgHandlers.ShowProjectSettings(projectStorage))
	protected.GET("/organizations/:orgID/projects/:id/edit", orgHandlers.EditProjectForm(projectStorage))
	protected.PUT("/organizations/:orgID/projects/:id", orgHandlers.UpdateProject(projectStorage))
	protected.DELETE("/organizations/:orgID/projects/:id", orgHandlers.DeleteProject(projectStorage))
	protected.POST("/organizations/:orgID/projects/:id/pin", orgHandlers.PinProject(projectStorage))
	protected.POST("/organizations/:orgID/projects/:id/unpin", orgHandlers.UnpinProject(projectStorage))
	protected.POST("/organizations/:orgID/projects/:id/archive", orgHandlers.ArchiveProject(projectStorage))
	protected.POST("/organizations/:orgID/projects/:id/unarchive", orgHandlers.UnarchiveProject(projectStorage))

	// Task routes
	protected.POST("/organizations/:orgID/projects/:projectID/tasks", orgHandlers.CreateTask(taskStorage))
	protected.GET("/organizations/:orgID/projects/:projectID/tasks", orgHandlers.ListTasks(taskStorage))
	protected.GET("/organizations/:orgID/projects/:projectID/tasks/archived", orgHandlers.ListArchivedTasks(taskStorage))
	protected.GET("/organizations/:orgID/projects/:projectID/tasks/new", orgHandlers.NewTaskForm)
	protected.GET("/organizations/:orgID/projects/:projectID/tasks/cancel", orgHandlers.CancelTaskForm)
	protected.GET("/organizations/:orgID/projects/:projectID/tasks/:id", orgHandlers.ShowTask(taskStorage))
	protected.GET("/organizations/:orgID/projects/:projectID/tasks/:id/edit", orgHandlers.EditTaskForm(taskStorage))
	protected.PUT("/organizations/:orgID/projects/:projectID/tasks/:id", orgHandlers.UpdateTask(taskStorage))
	protected.DELETE("/organizations/:orgID/projects/:projectID/tasks/:id", orgHandlers.DeleteTask(taskStorage))
	protected.POST("/organizations/:orgID/projects/:projectID/tasks/:id/pin", orgHandlers.PinTask(taskStorage))
	protected.POST("/organizations/:orgID/projects/:projectID/tasks/:id/unpin", orgHandlers.UnpinTask(taskStorage))
	protected.POST("/organizations/:orgID/projects/:projectID/tasks/:id/archive", orgHandlers.ArchiveTask(taskStorage))
	protected.POST("/organizations/:orgID/projects/:projectID/tasks/:id/unarchive", orgHandlers.UnarchiveTask(taskStorage))
}
