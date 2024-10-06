package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"tasks/internal/storage/postgres"
)

type StorageInterface interface {
	GetTasks() ([]postgres.Task, error)
	CreateTask(title string, description string) error
	DeleteTask(taskID string) error
	GetTaskByID(taskID string) (postgres.Task, error)
	UpdateTask(taskID string, title string, description string, isDone bool) error
}

type App struct {
	Storage StorageInterface
}

func (a *App) HomeHandler(c echo.Context) error {
	tasks, err := a.Storage.GetTasks()
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"tasks": tasks,
	})
}

func (a *App) CreateTaskHandler(c echo.Context) error {
	title := c.FormValue("title")
	description := c.FormValue("description")

	err := a.Storage.CreateTask(title, description)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, "/")
}

func (a *App) DeleteTaskHandler(c echo.Context) error {
	taskID := c.Param("id")

	err := a.Storage.DeleteTask(taskID)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func (a *App) EditTaskFormHandler(c echo.Context) error {
	taskID := c.Param("id")

	task, err := a.Storage.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "edit_task.html", map[string]interface{}{
		"task": task,
	})
}

func (a *App) UpdateTaskHandler(c echo.Context) error {
	taskID := c.Param("id")

	title := c.FormValue("title")
	description := c.FormValue("description")
	isDone := c.FormValue("is_done") == "on" //

	err := a.Storage.UpdateTask(taskID, title, description, isDone)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, "/")
}
