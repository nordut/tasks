package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"tasks/internal/storage/postgres"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock StorageInterface
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetTasks() ([]postgres.Task, error) {
	args := m.Called()
	return args.Get(0).([]postgres.Task), args.Error(1)
}

func (m *MockStorage) CreateTask(title string, description string) error {
	args := m.Called(title, description)
	return args.Error(0)
}

func (m *MockStorage) DeleteTask(taskID string) error {
	args := m.Called(taskID)
	return args.Error(0)
}

func (m *MockStorage) GetTaskByID(taskID string) (postgres.Task, error) {
	args := m.Called(taskID)
	return args.Get(0).(postgres.Task), args.Error(1)
}

func (m *MockStorage) UpdateTask(taskID string, title string, description string, isDone bool) error {
	args := m.Called(taskID, title, description, isDone)
	return args.Error(0)
}

type TemplateRenderer struct{}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	_, err := w.Write([]byte("Test Task"))
	return err
}

func Test_HomeHandler(t *testing.T) {
	e := echo.New()
	e.Renderer = &TemplateRenderer{}

	mockStorage := new(MockStorage)
	app := &App{Storage: mockStorage}

	mockStorage.On("GetTasks").Return([]postgres.Task{
		{Title: "Test Task", Description: "Test Description", IsDone: false, ID: 1},
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, app.HomeHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Test Task")
	}

	mockStorage.AssertExpectations(t)
}

func TestCreateTaskHandler(t *testing.T) {
	e := echo.New()
	e.Renderer = &TemplateRenderer{}

	mockStorage := new(MockStorage)
	app := &App{Storage: mockStorage}

	form := bytes.NewBufferString("title=New Task&description=New Description")
	req := httptest.NewRequest(http.MethodPost, "/tasks/create", form) // Исправлено на /tasks/create
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockStorage.On("CreateTask", "New Task", "New Description").Return(nil)

	if assert.NoError(t, app.CreateTaskHandler(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, "/", rec.Header().Get(echo.HeaderLocation))
	}

	mockStorage.AssertExpectations(t)
}

func TestApp_DeleteTaskHandler(t *testing.T) {
	e := echo.New()
	e.Renderer = &TemplateRenderer{}

	mockStorage := new(MockStorage)
	app := &App{Storage: mockStorage}
	taskID := "1"
	req := httptest.NewRequest(http.MethodPost, "/tasks/delete/"+taskID, nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/tasks/delete/:id")
	c.SetParamNames("id")
	c.SetParamValues(taskID)

	mockStorage.On("DeleteTask", taskID).Return(nil) // Возвращаем nil, чтобы имитировать успешное удаление

	if assert.NoError(t, app.DeleteTaskHandler(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, "/", rec.Header().Get(echo.HeaderLocation))
	}

	mockStorage.AssertExpectations(t)
}

func TestApp_EditTaskHandler(t *testing.T) {
	e := echo.New()
	e.Renderer = &TemplateRenderer{}

	mockStorage := new(MockStorage)
	app := &App{Storage: mockStorage}

	taskID := "1"
	req := httptest.NewRequest(http.MethodGet, "/tasks/edit/"+taskID, nil) // Измените на GET, если это форма редактирования
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/tasks/edit/:id")
	c.SetParamNames("id")
	c.SetParamValues(taskID)

	resTask := postgres.Task{ID: 1, Title: "Test Task", Description: "Test Description", IsDone: false}
	mockStorage.On("GetTaskByID", taskID).Return(resTask, nil)

	if assert.NoError(t, app.EditTaskFormHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	mockStorage.AssertExpectations(t)
}

func TestApp_UpdateTaskHandler(t *testing.T) {
	e := echo.New()
	e.Renderer = &TemplateRenderer{}

	mockStorage := new(MockStorage)
	app := &App{Storage: mockStorage}

	taskID := "1"
	form := bytes.NewBufferString("title=Changed Task&description=Changed Description&is_done=on")
	req := httptest.NewRequest(http.MethodPost, "/tasks/update/"+taskID, form)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tasks/update/:id")
	c.SetParamNames("id")
	c.SetParamValues(taskID)

	mockStorage.On("UpdateTask", taskID, "Changed Task", "Changed Description", true).Return(nil)

	if assert.NoError(t, app.UpdateTaskHandler(c)) {
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, "/", rec.Header().Get(echo.HeaderLocation))
	}

	mockStorage.AssertExpectations(t)
}
