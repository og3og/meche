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

// OrganizationRequest represents the request body for creating or updating an organization
type OrganizationRequest struct {
	Name        string `form:"name"`
	Description string `form:"description"`
}

// ValidationError represents field-specific validation errors
type ValidationError struct {
	Errors map[string]string `json:"errors"`
}

// Validate performs custom validation on the request
func (r *OrganizationRequest) Validate() *ValidationError {
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

// CreateOrganization handles the creation of a new organization
func CreateOrganization(orgStorage storage.OrganizationStorage, memberStorage storage.OrganizationMemberStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req OrganizationRequest
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

		// Get user from context
		user, ok := c.Get("user").(goth.User)
		if !ok {
			return c.String(http.StatusUnauthorized, "User not authenticated")
		}

		// Create new organization
		org := models.NewOrganization(req.Name, req.Description)

		// Store organization
		if err := orgStorage.CreateOrganization(org); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to create organization")
		}

		// Create organization membership for the owner
		member := models.NewOrganizationMember(org.ID, user.UserID, models.RoleOwner)
		if err := memberStorage.CreateMember(member); err != nil {
			// If we fail to create the membership, we should clean up the organization
			_ = orgStorage.DeleteOrganization(org.ID)
			return c.String(http.StatusInternalServerError, "Failed to create organization membership")
		}

		// Return a redirect response to the specific organization view
		c.Response().Header().Set("HX-Redirect", "/organizations/"+org.ID)
		return c.NoContent(http.StatusOK)
	}
}

// ListOrganizations returns a list of all organizations
func ListOrganizations(orgStorage storage.OrganizationStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		orgs, err := orgStorage.ListOrganizations()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch organizations")
		}
		return pages.OrganizationList(orgs).Render(c.Request().Context(), c.Response().Writer)
	}
}

// NewOrganizationForm renders the organization creation form
func NewOrganizationForm(c echo.Context) error {
	return pages.OrganizationForm().Render(c.Request().Context(), c.Response().Writer)
}

// CancelOrganizationForm clears the organization form
func CancelOrganizationForm(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

// DeleteOrganization handles organization deletion
func DeleteOrganization(orgStorage storage.OrganizationStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if err := orgStorage.DeleteOrganization(id); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to delete organization")
		}

		// Return empty response with 200 status to trigger HTMX swap
		c.Response().Header().Set("HX-Trigger", "organization-deleted")
		return c.NoContent(http.StatusOK)
	}
}

// EditOrganizationForm renders the organization edit form
func EditOrganizationForm(orgStorage storage.OrganizationStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		org, err := orgStorage.GetOrganization(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Organization not found")
		}
		return pages.OrganizationEdit(org).Render(c.Request().Context(), c.Response().Writer)
	}
}

// UpdateOrganization handles organization updates
func UpdateOrganization(orgStorage storage.OrganizationStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		org, err := orgStorage.GetOrganization(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Organization not found")
		}

		var req OrganizationRequest
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

		org.Update(req.Name, req.Description)
		if err := orgStorage.UpdateOrganization(org); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to update organization")
		}

		// Return a redirect response
		c.Response().Header().Set("HX-Redirect", "/organizations/"+id)
		return c.NoContent(http.StatusOK)
	}
}

// ShowOrganization displays the details of a specific organization
func ShowOrganization(orgStorage storage.OrganizationStorage, projectStorage storage.ProjectStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		org, err := orgStorage.GetOrganization(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Organization not found")
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

		// Get all projects for the organization and sort them alphabetically
		projects, err := projectStorage.ListProjects(org.ID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch projects")
		}

		// Sort projects alphabetically by name
		sort.Slice(projects, func(i, j int) bool {
			return strings.ToLower(projects[i].Name) < strings.ToLower(projects[j].Name)
		})

		return pages.OrganizationDetails(user, organizations, org, projects).Render(c.Request().Context(), c.Response().Writer)
	}
}
