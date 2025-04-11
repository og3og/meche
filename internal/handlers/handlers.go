package handlers

import (
	"meche/templates"

	"github.com/labstack/echo/v4"
)

// HandleHome handles the home page route
func HandleHome(c echo.Context) error {
	component := templates.Home()
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// HandleGreet handles the greeting API endpoint
func HandleGreet(c echo.Context) error {
	return c.HTML(200, "<span class='text-green-600'>YO!a12Hello from the server!</span>")
}
