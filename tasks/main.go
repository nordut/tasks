package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"tasks/internal/config"
	"tasks/internal/handlers"
	"tasks/internal/storage/postgres"

	"html/template"
	"io"
	"log"
)

type TemplateRegistry struct {
	templates *template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	cfg := config.MustLoad()

	conn, err := postgres.NewConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app := &handlers.App{Storage: conn}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/static", "static")
	e.Renderer = &TemplateRegistry{
		templates: template.Must(template.ParseGlob("view/*.html")),
	}

	e.GET("/", app.HomeHandler)
	e.GET("/tasks/new", func(c echo.Context) error {
		return c.File("view/new_task.html")
	})

	e.POST("/tasks/create", app.CreateTaskHandler)
	e.POST("/tasks/delete/:id", app.DeleteTaskHandler)

	e.GET("/tasks/edit/:id", app.EditTaskFormHandler)
	e.POST("/tasks/update/:id", app.UpdateTaskHandler)

	e.Logger.Fatal(e.Start(cfg.HTTPServer.Address))
}
