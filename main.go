package main

import (
	"fmt"
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
	e.GET("/hello", helloGet)
	e.POST("/hello", helloPost)

	e.File("/index.html", "static/index.html")

	e.Logger.Fatal(e.Start(":8000"))
}

type helloTemplate struct {
	Title   string
	Message string
}

func helloGet(c echo.Context) error {
	ht := helloTemplate{
		Title:   "Send values",
		Message: "type name and password:",
	}
	return c.Render(http.StatusOK, "hello", ht)
}

func helloPost(c echo.Context) error {
	nm := c.FormValue("name")
	pw := c.FormValue("password")
	ht := helloTemplate{
		Title:   "Send values",
		Message: fmt.Sprintf("name: %s, password: %s", nm, pw),
	}
	return c.Render(http.StatusOK, "hello", ht)
}
