package controller

import (
	"net/http"

	"CNSA-Apply/models"

	session "github.com/ipfans/echo-session"
	"github.com/labstack/echo"
)

// Index : Main Page
func Index(c echo.Context) error {
	session := session.Default(c)
	if session.Get("loginID") == nil {
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
	if models.Login(c.FormValue("loginID"), c.FormValue("loginPassword")) {
		session := session.Default(c)
		session.Set("loginID", c.FormValue("loginID"))
		session.Save()

		return c.Redirect(http.StatusMovedPermanently, "/")
	}
	return c.Redirect(http.StatusMovedPermanently, "/login?error=loginFailed")
}
