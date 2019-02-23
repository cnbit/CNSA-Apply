package main

import (
	"CNSA-Apply/controller"
	"io"
	"text/template"

	"github.com/ipfans/echo-session"
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

	// Session 설정
	store := session.NewCookieStore([]byte("secret"))
	e.Use(session.Sessions("CASESSION", store))

	// Set template renderer
	// We uses standard golang template
	e.Renderer = t

	// Set static serve files
	e.Static("/assets", "static")

	// Handle requests
	// Filter by path
	// ================ 학생 페이지 ================
	// 메인 페이지
	e.GET("/", controller.Index, controller.AuthAPI)

	// Login Page
	e.GET("/login", controller.Login)
	e.POST("/login", controller.LoginPost)
	// Logout
	e.GET("/logout", controller.Logout)

	// ================ 학생 API ================
	// 신청하기
	e.POST("/api/apply", controller.ApplyAPI, controller.AuthAPI)
	// 자신의 신청내역 가져오기
	e.GET("/api/getApplys", controller.GetApplysAPI, controller.AuthAPI)

	// ================ 교사 페이지 ================
	// Login Page
	e.GET("/admin/login", controller.AdminLogin)
	e.POST("/admin/login", controller.AdminLoginPost)
	// Logout
	e.GET("/admin/logout", controller.AdminLogout)

	a := e.Group("/admin")
	a.Use(controller.AdminAuthAPI)

	// 메인 페이지
	a.GET("/", controller.AdminIndex)

	// ================ 교사 API ================

	// Start web server
	e.Start(":80")
}
