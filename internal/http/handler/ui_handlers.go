package handler

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"net/http"
	"taletracker.com/internal/model"
)

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func Home(c echo.Context) error {
	return c.Render(http.StatusOK, "home", "")
}

func List(c echo.Context) error {
	return c.Render(http.StatusOK, "list", []model.Tale{})
}
