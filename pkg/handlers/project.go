package handlers

import (
	"meche/pkg/models"
	"meche/pkg/storage"
	"meche/templates/pages"
	"net/http"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
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

		// Return a redirect response to the project view page
		c.Response().Header().Set("HX-Redirect", "/organizations/"+orgID+"/projects/"+project.ID)
		return c.NoContent(http.StatusOK)
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
func ShowProject(projectStorage storage.ProjectStorage, orgStorage storage.OrganizationStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		project, err := projectStorage.GetProject(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Project not found")
		}

		// Get the current user from the session
		user, ok := c.Get("user").(goth.User)
		if !ok {
			return c.String(http.StatusUnauthorized, "User not authenticated")
		}

		// Get all organizations for the user
		organizations, err := orgStorage.ListOrganizations()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch organizations")
		}

		// Get all projects for the current organization and sort them alphabetically
		projects, err := projectStorage.ListProjects(project.OrgID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch projects")
		}

		// Sort projects alphabetically by name
		sort.Slice(projects, func(i, j int) bool {
			return strings.ToLower(projects[i].Name) < strings.ToLower(projects[j].Name)
		})

		return pages.ProjectDetails(user, organizations, project, projects).Render(c.Request().Context(), c.Response().Writer)
	}
}

// ShowProjectOverview displays the project overview tab
func ShowProjectOverview(projectStorage storage.ProjectStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		project, err := projectStorage.GetProject(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Project not found")
		}
		return pages.ProjectOverview(project).Render(c.Request().Context(), c.Response().Writer)
	}
}

// ShowProjectSettings displays the project settings tab
func ShowProjectSettings(projectStorage storage.ProjectStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		project, err := projectStorage.GetProject(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Project not found")
		}
		return pages.ProjectSettings(project).Render(c.Request().Context(), c.Response().Writer)
	}
}

// NewProjectForm renders the project creation page
func NewProjectForm(orgStorage storage.OrganizationStorage, projectStorage storage.ProjectStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		orgID := c.Param("orgID")

		// Get the current user from the session
		user, ok := c.Get("user").(goth.User)
		if !ok {
			return c.String(http.StatusUnauthorized, "User not authenticated")
		}

		// Get all organizations for the user
		organizations, err := orgStorage.ListOrganizations()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch organizations")
		}

		// Get all projects for the current organization
		projects, err := projectStorage.ListProjects(orgID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch projects")
		}

		return pages.ProjectNew(user, organizations, orgID, projects).Render(c.Request().Context(), c.Response().Writer)
	}
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

// PinProject handles pinning a project
func PinProject(projectStorage storage.ProjectStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		project, err := projectStorage.GetProject(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Project not found")
		}

		project.Pin()
		if err := projectStorage.UpdateProject(project); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to pin project")
		}

		return pages.ProjectSettings(project).Render(c.Request().Context(), c.Response().Writer)
	}
}

// UnpinProject handles unpinning a project
func UnpinProject(projectStorage storage.ProjectStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		project, err := projectStorage.GetProject(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Project not found")
		}

		project.Unpin()
		if err := projectStorage.UpdateProject(project); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to unpin project")
		}

		return pages.ProjectSettings(project).Render(c.Request().Context(), c.Response().Writer)
	}
}

// ArchiveProject handles archiving a project
func ArchiveProject(projectStorage storage.ProjectStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		project, err := projectStorage.GetProject(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Project not found")
		}

		project.Archive()
		if err := projectStorage.UpdateProject(project); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to archive project")
		}

		return c.NoContent(http.StatusOK)
	}
}

// UnarchiveProject handles unarchiving a project
func UnarchiveProject(projectStorage storage.ProjectStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		project, err := projectStorage.GetProject(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Project not found")
		}

		project.Unarchive()
		if err := projectStorage.UpdateProject(project); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to unarchive project")
		}

		return c.NoContent(http.StatusOK)
	}
}
