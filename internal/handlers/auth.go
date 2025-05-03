package handlers

import (
	"fmt"
	"meche/templates/pages"
	"net/http"

	"meche/pkg/storage"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

const sessionName = "gothic-session"

// HandleLogin renders the login page
func HandleLogin(c echo.Context) error {
	return pages.Login().Render(c.Request().Context(), c.Response().Writer)
}

// HandleAuth initiates the OAuth2 authentication flow
func HandleAuth(c echo.Context) error {
	// Set the provider in the URL query
	c.Request().URL.RawQuery = "provider=google"

	// Begin the authentication process
	gothic.BeginAuthHandler(c.Response(), c.Request())
	return nil
}

// HandleCallback handles the OAuth2 callback from Google
func HandleCallback(c echo.Context) error {
	// Set the provider in the URL query
	//c.Request().URL.RawQuery = "provider=google"

	// Get the user from the session
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		fmt.Printf("Error completing auth: %v\n", err)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
	}

	// Get the session
	session, err := gothic.Store.Get(c.Request(), sessionName)
	if err != nil {
		fmt.Printf("Error getting session in callback: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error getting session")
	}

	// Store user info in session
	session.Values["user"] = user
	if err := session.Save(c.Request(), c.Response()); err != nil {
		fmt.Printf("Error saving session: %v\n", err)
		return c.String(http.StatusInternalServerError, "Error saving session")
	}

	// Debug: Print user info
	fmt.Printf("User logged in: %+v\n", user)

	return c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
}

// HandleLogout logs out the user
func HandleLogout(c echo.Context) error {
	// Get the session
	session, err := gothic.Store.Get(c.Request(), sessionName)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error getting session")
	}

	// Clear all session values
	session.Values = make(map[interface{}]interface{})

	// Set MaxAge to -1 to delete the cookie
	session.Options.MaxAge = -1

	// Save the session
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, "Error saving session")
	}

	// Call gothic logout
	gothic.Logout(c.Response(), c.Request())

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

// HandleDashboard renders the dashboard page for authenticated users
func HandleDashboard(c echo.Context) error {
	session, err := gothic.Store.Get(c.Request(), sessionName)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	user, ok := session.Values["user"].(goth.User)
	if !ok {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// Get organization storage from context
	orgStorage, ok := c.Get("organization_storage").(storage.OrganizationStorage)
	if !ok {
		return c.String(http.StatusInternalServerError, "Organization storage not found")
	}

	// Get all organizations
	organizations, err := orgStorage.ListOrganizations()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error fetching organizations: %v", err))
	}

	return pages.Dashboard(user, organizations).Render(c.Request().Context(), c.Response().Writer)
}
