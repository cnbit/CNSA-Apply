package controller

import (
	"CNSA-Apply/models"
	"fmt"
	"net/http"
	"time"


	session "github.com/ipfans/echo-session"
	"github.com/labstack/echo"
)

// AdminAuthAPI 로그인 인증 middleware
func AdminAuthAPI(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 로그인이 되어 있지 않으면 login page로 redirect
		session := session.Default(c)
		if session.Get("cnsanetID") == nil {
			return c.Redirect(http.StatusMovedPermanently, "/admin/login")
		}

		return next(c)
	}
}

// AdminLogin : Login Page
func AdminLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "adminLogin", nil)
}

// AdminLoginPost : Check a Login Data
func AdminLoginPost(c echo.Context) error {
	// Login 성공 시
	if models.TcrLogin(c.FormValue("loginID"), c.FormValue("loginPassword")) {
		// Session에 학번 저장
		session := session.Default(c)
		session.Set("cnsanetID", c.FormValue("loginID"))
		session.Save()

		return c.Redirect(http.StatusMovedPermanently, "/admin/")
	}

	// Login 실패 시
	return c.Redirect(http.StatusMovedPermanently, "/admin/login?error=loginFailed")
}

// AdminLogout : 로그아웃 - 세션 초기화
func AdminLogout(c echo.Context) error {
	// Session 초기화
	session := session.Default(c)
	session.Clear()
	session.Save()

	// 로그인 페이지로 빠이빠이
	return c.Redirect(http.StatusMovedPermanently, "/admin/login")
}

// AdminIndex : Main Page
func AdminIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "adminIndex", nil)
}

// AdminCancelHolydayAPI 공휴일 취소 API
func AdminCancelHolydayAPI(c echo.Context) error {
	day, err := time.Parse("2006-01-02", c.FormValue("date"))
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	err = models.DeleteHolyday(day)
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	return c.String(http.StatusOK, "success")
}
