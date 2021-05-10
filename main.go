package main

import (
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	t := &Template{
		templates: template.Must(template.ParseGlob("template/*.html")),
	}
	e.Renderer = t
	e.GET("/hello", hello)

	e.File("/index.html", "static/index.html")

	e.Logger.Fatal(e.Start(":8000"))
}

type helloTemplate struct {
	Title string
	Items []string
}

func hello(c echo.Context) error {
	ht := helloTemplate{
		Title: "Send values",
		Items: []string{"One", "Two", "Three"},
	}
	return c.Render(http.StatusOK, "hello", ht)
}
