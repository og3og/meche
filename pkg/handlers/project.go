package handlers

import (
	"meche/pkg/models"
	"meche/pkg/storage"
	"meche/templates/pages"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// ProjectRequest represents the request body for creating or updating a project
type ProjectRequest struct {
	Name        string `form:"name"`
	Description string `form:"description"`
}

// Validate performs custom validation on the request
func (r *ProjectRequest) Validate() *ValidationError {
	errors := make(map[string]string)

	if strings.TrimSpace(r.Name) == "" {
		errors["name"] = "Name cannot be empty"
	}
	if strings.TrimSpace(r.Description) == "" {
		errors["description"] = "Description cannot be empty"
	}

	if len(errors) > 0 {
		return &ValidationError{Errors: errors}
	}
	return nil
}

// CreateProject handles the creation of a new project
func CreateProject(projectStorage storage.ProjectStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		orgID := c.Param("orgID")
		var req ProjectRequest
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "Invalid form data")
		}

		// Trim whitespace
		req.Name = strings.TrimSpace(req.Name)
		req.Description = strings.TrimSpace(req.Description)

		// Validate request
		if validationErr := req.Validate(); validationErr != nil {
			// Return error box HTML
			if err := pages.ErrorBox(validationErr.Errors).Render(c.Request().Context(), c.Response().Writer); err != nil {
				return err
			}
			c.Response().Header().Set("HX-Retarget", "#error-box")
			return nil
		}

		// Create new project
		project := models.NewProject(req.Name, req.Description, orgID)

		// Store project
		if err := projectStorage.CreateProject(project); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to create project")
		}

		// Return the project list item HTML
		return pages.ProjectList([]*models.Project{project}).Render(c.Request().Context(), c.Response().Writer)
	}
}

// ListProjects returns a list of all projects for an organization
func ListProjects(projectStorage storage.ProjectStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		orgID := c.Param("orgID")
		projects, err := projectStorage.ListProjects(orgID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch projects")
		}
		return pages.ProjectList(projects).Render(c.Request().Context(), c.Response().Writer)
	}
}

// ShowProject displays the details of a specific project
func ShowProject(projectStorage storage.ProjectStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		project, err := projectStorage.GetProject(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Project not found")
		}

		return pages.ProjectDetails(project).Render(c.Request().Context(), c.Response().Writer)
	}
}

// NewProjectForm renders the project creation form
func NewProjectForm(c echo.Context) error {
	orgID := c.Param("orgID")
	return pages.ProjectForm(orgID).Render(c.Request().Context(), c.Response().Writer)
}

// CancelProjectForm clears the project form
func CancelProjectForm(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

// EditProjectForm renders the project edit form
func EditProjectForm(projectStorage storage.ProjectStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		project, err := projectStorage.GetProject(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Project not found")
		}
		return pages.EditProjectForm(project).Render(c.Request().Context(), c.Response().Writer)
	}
}

// UpdateProject handles project updates
func UpdateProject(projectStorage storage.ProjectStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		project, err := projectStorage.GetProject(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Project not found")
		}

		var req ProjectRequest
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "Invalid form data")
		}

		// Trim whitespace
		req.Name = strings.TrimSpace(req.Name)
		req.Description = strings.TrimSpace(req.Description)

		// Validate request
		if validationErr := req.Validate(); validationErr != nil {
			// Return error box HTML
			if err := pages.ErrorBox(validationErr.Errors).Render(c.Request().Context(), c.Response().Writer); err != nil {
				return err
			}
			c.Response().Header().Set("HX-Retarget", "#error-box")
			return nil
		}

		project.Update(req.Name, req.Description)
		if err := projectStorage.UpdateProject(project); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to update project")
		}

		// Return a redirect response
		c.Response().Header().Set("HX-Redirect", "/organizations/"+project.OrgID+"/projects/"+project.ID)
		return c.NoContent(http.StatusOK)
	}
}

// DeleteProject handles project deletion
func DeleteProject(projectStorage storage.ProjectStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		orgID := c.Param("orgID")

		if err := projectStorage.DeleteProject(id); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to delete project")
		}

		// Return redirect response
		return c.Redirect(http.StatusSeeOther, "/organizations/"+orgID)
	}
}
