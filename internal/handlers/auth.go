package handlers

import (
	"fmt"
	"meche/templates/pages"
	"net/http"

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

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

// HandleLogout logs out the user
func HandleLogout(c echo.Context) error {
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

	return pages.Dashboard(user).Render(c.Request().Context(), c.Response().Writer)
}
