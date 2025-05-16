package handlers

import (
	"meche/pkg/models"
	"meche/pkg/storage"
	"meche/templates/pages"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// TaskRequest represents the request body for creating or updating a task
type TaskRequest struct {
	Name        string `form:"name"`
	Description string `form:"description"`
}

// Validate performs custom validation on the request
func (r *TaskRequest) Validate() *ValidationError {
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

// CreateTask handles the creation of a new task
func CreateTask(taskStorage storage.TaskStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		orgID := c.Param("orgID")
		projectID := c.Param("projectID")

		// Parse form data
		name := c.FormValue("name")
		description := c.FormValue("description")
		dueAtStr := c.FormValue("due_at")

		// Create new task
		task := models.NewTask(name, description, orgID, projectID)

		// Parse and set due date if provided
		if dueAtStr != "" {
			dueAt, err := time.Parse("2006-01-02T15:04", dueAtStr)
			if err != nil {
				return c.String(http.StatusBadRequest, "Invalid due date format")
			}
			task.SetDueDate(&dueAt)
		}

		// Save task
		if err := taskStorage.CreateTask(task); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to create task")
		}

		// Get updated task list
		tasks, err := taskStorage.ListTasks(projectID, orgID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch tasks")
		}

		// Filter out archived tasks
		activeTasks := make([]*models.Task, 0)
		for _, t := range tasks {
			if t.ArchivedAt == nil {
				activeTasks = append(activeTasks, t)
			}
		}

		// Return the updated task list
		return pages.TaskList(activeTasks, orgID, projectID).Render(c.Request().Context(), c.Response().Writer)
	}
}

// ListTasks returns a list of all tasks for a project
func ListTasks(taskStorage storage.TaskStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		orgID := c.Param("orgID")
		projectID := c.Param("projectID")
		tasks, err := taskStorage.ListTasks(projectID, orgID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch tasks")
		}

		// Filter out archived tasks
		activeTasks := make([]*models.Task, 0)
		for _, task := range tasks {
			if task.ArchivedAt == nil {
				activeTasks = append(activeTasks, task)
			}
		}

		return pages.TaskList(activeTasks, orgID, projectID).Render(c.Request().Context(), c.Response().Writer)
	}
}

// ListArchivedTasks returns a list of archived tasks for a project
func ListArchivedTasks(taskStorage storage.TaskStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		orgID := c.Param("orgID")
		projectID := c.Param("projectID")
		tasks, err := taskStorage.ListTasks(projectID, orgID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch tasks")
		}

		// Filter for archived tasks only
		archivedTasks := make([]*models.Task, 0)
		for _, task := range tasks {
			if task.ArchivedAt != nil {
				archivedTasks = append(archivedTasks, task)
			}
		}

		return pages.ArchivedTasksList(archivedTasks, orgID, projectID).Render(c.Request().Context(), c.Response().Writer)
	}
}

// validateTaskAccess checks if the task belongs to the specified organization and project
func validateTaskAccess(taskStorage storage.TaskStorage, taskID, orgID, projectID string) (*models.Task, error) {
	task, err := taskStorage.GetTask(taskID, orgID)
	if err != nil {
		return nil, err
	}
	if task.ProjectID != projectID {
		return nil, echo.NewHTTPError(http.StatusForbidden, "Task does not belong to this project")
	}
	return task, nil
}

// ShowTask displays the details of a specific task
func ShowTask(taskStorage storage.TaskStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		orgID := c.Param("orgID")
		projectID := c.Param("projectID")

		task, err := validateTaskAccess(taskStorage, id, orgID, projectID)
		if err != nil {
			return err
		}

		return pages.TaskDetails(task).Render(c.Request().Context(), c.Response().Writer)
	}
}

// NewTaskForm renders the task creation form
func NewTaskForm(c echo.Context) error {
	orgID := c.Param("orgID")
	projectID := c.Param("projectID")
	return pages.TaskForm(orgID, projectID).Render(c.Request().Context(), c.Response().Writer)
}

// CancelTaskForm clears the task form
func CancelTaskForm(c echo.Context) error {
	return c.HTML(http.StatusOK, "<div id=\"task-form\"></div>")
}

// EditTaskForm renders the task edit form
func EditTaskForm(taskStorage storage.TaskStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		orgID := c.Param("orgID")
		projectID := c.Param("projectID")

		task, err := validateTaskAccess(taskStorage, id, orgID, projectID)
		if err != nil {
			return err
		}
		return pages.EditTaskForm(task).Render(c.Request().Context(), c.Response().Writer)
	}
}

// UpdateTask handles task updates
func UpdateTask(taskStorage storage.TaskStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		orgID := c.Param("orgID")
		projectID := c.Param("projectID")

		// Get existing task
		task, err := taskStorage.GetTask(id, orgID)
		if err != nil {
			return c.String(http.StatusNotFound, "Task not found")
		}

		// Validate task access
		if task.OrgID != orgID || task.ProjectID != projectID {
			return c.String(http.StatusForbidden, "Access denied")
		}

		// Update task fields
		task.Name = c.FormValue("name")
		task.Description = c.FormValue("description")
		dueAtStr := c.FormValue("due_at")

		// Parse and set due date if provided
		if dueAtStr != "" {
			dueAt, err := time.Parse("2006-01-02T15:04", dueAtStr)
			if err != nil {
				return c.String(http.StatusBadRequest, "Invalid due date format")
			}
			task.SetDueDate(&dueAt)
		} else {
			task.SetDueDate(nil)
		}

		// Save task
		if err := taskStorage.UpdateTask(task, orgID); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to update task")
		}

		// Get updated task list
		tasks, err := taskStorage.ListTasks(projectID, orgID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch tasks")
		}

		// Filter out archived tasks
		activeTasks := make([]*models.Task, 0)
		for _, t := range tasks {
			if t.ArchivedAt == nil {
				activeTasks = append(activeTasks, t)
			}
		}

		// Return the updated task list
		return pages.TaskList(activeTasks, orgID, projectID).Render(c.Request().Context(), c.Response().Writer)
	}
}

// DeleteTask handles task deletion
func DeleteTask(taskStorage storage.TaskStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		orgID := c.Param("orgID")
		projectID := c.Param("projectID")

		// Validate task access first
		_, err := validateTaskAccess(taskStorage, id, orgID, projectID)
		if err != nil {
			return err
		}

		// Delete the task
		if err := taskStorage.DeleteTask(id, orgID); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to delete task")
		}

		// Get updated task list
		tasks, err := taskStorage.ListTasks(projectID, orgID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch tasks")
		}

		// Filter out archived tasks
		activeTasks := make([]*models.Task, 0)
		for _, task := range tasks {
			if task.ArchivedAt == nil {
				activeTasks = append(activeTasks, task)
			}
		}

		// Return the updated task list
		return pages.TaskList(activeTasks, orgID, projectID).Render(c.Request().Context(), c.Response().Writer)
	}
}

// PinTask handles pinning a task
func PinTask(taskStorage storage.TaskStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		orgID := c.Param("orgID")
		projectID := c.Param("projectID")

		task, err := validateTaskAccess(taskStorage, id, orgID, projectID)
		if err != nil {
			return err
		}

		task.Pin()
		if err := taskStorage.UpdateTask(task, orgID); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to pin task")
		}

		// Get updated task list
		tasks, err := taskStorage.ListTasks(projectID, orgID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch tasks")
		}

		// Filter out archived tasks
		activeTasks := make([]*models.Task, 0)
		for _, task := range tasks {
			if task.ArchivedAt == nil {
				activeTasks = append(activeTasks, task)
			}
		}

		// Return the updated task list
		return pages.TaskList(activeTasks, orgID, projectID).Render(c.Request().Context(), c.Response().Writer)
	}
}

// UnpinTask handles unpinning a task
func UnpinTask(taskStorage storage.TaskStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		orgID := c.Param("orgID")
		projectID := c.Param("projectID")

		task, err := validateTaskAccess(taskStorage, id, orgID, projectID)
		if err != nil {
			return err
		}

		task.Unpin()
		if err := taskStorage.UpdateTask(task, orgID); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to unpin task")
		}

		// Get updated task list
		tasks, err := taskStorage.ListTasks(projectID, orgID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch tasks")
		}

		// Filter out archived tasks
		activeTasks := make([]*models.Task, 0)
		for _, task := range tasks {
			if task.ArchivedAt == nil {
				activeTasks = append(activeTasks, task)
			}
		}

		// Return the updated task list
		return pages.TaskList(activeTasks, orgID, projectID).Render(c.Request().Context(), c.Response().Writer)
	}
}

// ArchiveTask handles archiving a task
func ArchiveTask(taskStorage storage.TaskStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		orgID := c.Param("orgID")
		projectID := c.Param("projectID")

		task, err := validateTaskAccess(taskStorage, id, orgID, projectID)
		if err != nil {
			return err
		}

		task.Archive()
		if err := taskStorage.UpdateTask(task, orgID); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to archive task")
		}

		// Get updated task list
		tasks, err := taskStorage.ListTasks(projectID, orgID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch tasks")
		}

		// Filter out archived tasks
		activeTasks := make([]*models.Task, 0)
		for _, task := range tasks {
			if task.ArchivedAt == nil {
				activeTasks = append(activeTasks, task)
			}
		}

		// Return the updated task list
		return pages.TaskList(activeTasks, orgID, projectID).Render(c.Request().Context(), c.Response().Writer)
	}
}

// UnarchiveTask handles unarchiving a task
func UnarchiveTask(taskStorage storage.TaskStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		orgID := c.Param("orgID")
		projectID := c.Param("projectID")

		task, err := validateTaskAccess(taskStorage, id, orgID, projectID)
		if err != nil {
			return err
		}

		task.Unarchive()
		if err := taskStorage.UpdateTask(task, orgID); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to unarchive task")
		}

		// Get updated task list
		tasks, err := taskStorage.ListTasks(projectID, orgID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to fetch tasks")
		}

		// Return the updated task list
		return pages.TaskList(tasks, orgID, projectID).Render(c.Request().Context(), c.Response().Writer)
	}
}
