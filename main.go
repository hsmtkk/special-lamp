package main

import (
	"io"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

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
	msg := "login name & password"

	ses, _ := session.Get("hello-session", c)

	logined, _ := ses.Values["login"].(bool)
	lname, _ := ses.Values["name"].(string)
	if logined {
		msg = "logined: " + lname
	}

	ht := helloTemplate{
		Title:   "Session",
		Message: msg,
	}
	return c.Render(http.StatusOK, "hello", ht)
}

func helloPost(c echo.Context) error {
	msg := "login name & password"

	ses, _ := session.Get("hello-session", c)
	ses.Values["login"] = false
	ses.Values["name"] = ""
	nm := c.FormValue("name")
	pw := c.FormValue("pass")
	if nm == pw {
		ses.Values["login"] = true
		ses.Values["name"] = nm
	}
	ses.Save(c.Request(), c.Response())

	logined := ses.Values["login"].(bool)
	lname := ses.Values["name"].(string)
	if logined {
		msg = "logined: " + lname
	}

	ht := helloTemplate{
		Title:   "Session",
		Message: msg,
	}
	return c.Render(http.StatusOK, "hello", ht)
}
