package controller

import (
	"net/http"

	"CNSA-Apply/models"

	session "github.com/ipfans/echo-session"
	"github.com/labstack/echo"
)

// Index : Main Page
func Index(c echo.Context) error {
	// 로그인이 되어 있지 않으면 login page로 redirect
	session := session.Default(c)
	if session.Get("studentNumber") == nil {
		return c.Redirect(http.StatusMovedPermanently, "/login")
	}

	return c.Render(http.StatusOK, "index", nil)
}

// Login : Login Page
func Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

// LoginPost : Check a Login Data
func LoginPost(c echo.Context) error {
	// Login 성공 시
	if models.Login(c.FormValue("loginID"), c.FormValue("loginPassword")) {
		// Session에 학번 저장
		session := session.Default(c)
		session.Set("studentNumber", c.FormValue("loginID"))
		session.Save()

		return c.Redirect(http.StatusMovedPermanently, "/")
	}
	// Login 실패 시
	return c.Redirect(http.StatusMovedPermanently, "/login?error=loginFailed")
}
