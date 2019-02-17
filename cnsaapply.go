package main

import (
	"CNSA-Apply/controller"
	"io"
	"text/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Template is a custom html/template renderer for Echo framework
type Template struct {
	templates *template.Template
}

// Render renders a template document
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	t := &Template{
		templates: template.Must(template.New("").Delims("[[", "]]").ParseGlob("view/*.html")),
	}

	e := echo.New()

	// Set middlewares
	// Logger: loging all request and responses
	// Recover: Recover main thread if it fails
	e.Use(middleware.Logger(), middleware.Recover())

	// Set template renderer
	// We uses standard golang template
	e.Renderer = t

	// Set static serve files
	e.Static("/assets", "static")

	// Hanle requests
	// Filter by path
	e.GET("/login", controller.Login)

	// Start web server
	e.Start(":80")
}
