package controller

import (
	"net/http"

	"github.com/labstack/echo"
)

// Login : Login Page
func Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}
