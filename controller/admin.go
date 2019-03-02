package controller

import (
	"CNSA-Apply/models"
	"net/http"
	"strconv"
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

// AdminSelectTime : 시간대 선택 Page
func AdminSelectTime(c echo.Context) error {
	return c.Render(http.StatusOK, "adminSelectTime", nil)
}

// AdminSelectForm : 면학실 선택 Page
func AdminSelectForm(c echo.Context) error {
	return c.Render(http.StatusOK, "adminSelectForm", map[string]interface{}{
		"date":   c.QueryParam("date"),
		"period": c.QueryParam("period"),
	})
}

// AdminSelectArea : 면학실 선택 Page
func AdminSelectArea(c echo.Context) error {
	return c.Render(http.StatusOK, "adminSelectArea", map[string]interface{}{
		"date":   c.QueryParam("date"),
		"period": c.QueryParam("period"),
	})
}

// AdminAViewA : 창학관 신청내역 보기 Page
func AdminAViewA(c echo.Context) error {
	return c.Render(http.StatusOK, "adminAViewA", nil)
}

// AdminAViewB : 창학관 신청내역 보기 Page
func AdminAViewB(c echo.Context) error {
	return c.Render(http.StatusOK, "adminAViewB", nil)
}

// AdminAViewC : 창학관 신청내역 보기 Page
func AdminAViewC(c echo.Context) error {
	return c.Render(http.StatusOK, "adminAViewC", nil)
}

// AdminAViewD : 창학관 신청내역 보기 Page
func AdminAViewD(c echo.Context) error {
	return c.Render(http.StatusOK, "adminAViewD", nil)
}

// AdminAViewE : 창학관 신청내역 보기 Page
func AdminAViewE(c echo.Context) error {
	return c.Render(http.StatusOK, "adminAViewE", nil)
}

// AdminAViewF : 창학관 신청내역 보기 Page
func AdminAViewF(c echo.Context) error {
	return c.Render(http.StatusOK, "adminAViewF", nil)
}

// AdminAViewG : 창학관 신청내역 보기 Page
func AdminAViewG(c echo.Context) error {
	return c.Render(http.StatusOK, "adminAViewG", nil)
}

// AdminBView : 자율관 신청내역 보기 Page
func AdminBView(c echo.Context) error {
	return c.Render(http.StatusOK, "adminBView", nil)
}

// AdminHolydays : 공휴일 관리 Page
func AdminHolydays(c echo.Context) error {
	return c.Render(http.StatusOK, "adminHolydays", nil)
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

// AdminGetApplysAPI 신청내역 가져오기 API
func AdminGetApplysAPI(c echo.Context) error {
	date, err := time.Parse("2006-01-02", c.FormValue("date"))
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}
	return c.JSON(http.StatusOK, models.GetApplys(date, c.QueryParam("period"), c.QueryParam("form"), c.QueryParam("area")))
}

// AdminAddHolydayAPI 공휴일 추가하기
func AdminAddHolydayAPI(c echo.Context) error {
	day, err := time.Parse("2006-01-02", c.FormValue("date"))
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}
	err = models.AddHolyday(day, c.FormValue("holydayName"))
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

// AdminGetApplyMountAPI : 시간대에 해당하는 인원 수 가져오기
func AdminGetApplyMountAPI(c echo.Context) error {
	day, err := time.Parse("2006-01-02", c.QueryParam("date"))
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}
	return c.String(http.StatusOK, strconv.Itoa(models.GetApplyMount(day, c.QueryParam("period"), c.QueryParam("form"))))
}
