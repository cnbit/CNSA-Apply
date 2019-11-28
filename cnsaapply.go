package main

import (
	"CNSA-Apply/controller"
	"html/template"
	"io"

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
		templates: template.Must(template.New("").Delims("[[", "]]").ParseGlob("view/*/*.html")),
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
	// Login Page
	e.GET("/login", controller.Login)
	e.POST("/login", controller.LoginPost)
	// Logout
	e.GET("/logout", controller.Logout)

	// 메인 페이지
	e.GET("/", controller.Index, controller.AuthAPI)

	// 신청하기 - 면학실 선택
	e.GET("/apply/", controller.SelectForm, controller.AuthAPI)
	// 신청하기 - 시간대 선택
	e.GET("/apply/selectTime", controller.SelectTime, controller.AuthAPI)
	// 신청하기 - 구역 선택
	e.GET("/apply/selectArea", controller.SelectArea, controller.AuthAPI)
	e.POST("/apply/selectArea", controller.SelectAreaPOST, controller.AuthAPI)
	// 신청하기 - 구역 선택
	e.GET("/apply/selectSeatA", controller.SelectSeatA, controller.AuthAPI)
	e.GET("/apply/selectSeatB", controller.SelectSeatB, controller.AuthAPI)
	e.GET("/apply/selectSeatC", controller.SelectSeatC, controller.AuthAPI)
	e.GET("/apply/selectSeatD", controller.SelectSeatD, controller.AuthAPI)
	e.GET("/apply/selectSeatE", controller.SelectSeatE, controller.AuthAPI)
	e.GET("/apply/selectSeatF", controller.SelectSeatF, controller.AuthAPI)
	e.GET("/apply/selectSeatG", controller.SelectSeatG, controller.AuthAPI)
	// 신청하기 - 신청완료
	e.GET("/apply/applySuccess", controller.ApplySuccess, controller.AuthAPI)

	// 내정보
	e.GET("/user/", controller.MyPage, controller.AuthAPI)
	// 내정보 - 신청내역
	e.GET("/user/history", controller.ApplyHistory, controller.AuthAPI)
	// 내정보 - 계정관리
	e.GET("/user/account", controller.Account, controller.AuthAPI)
	e.POST("/user/account", controller.AccountPOST, controller.AuthAPI)
	e.GET("/user/changeSuccess", controller.ChangeSuccess, controller.AuthAPI)

	// 공지사항
	e.GET("/notice/", controller.Notices, controller.AuthAPI)
	// 공지사항 내용
	e.GET("/notice/content", controller.NoticeContent, controller.AuthAPI)

	// ================ 학생 API ================
	// 신청하기
	e.POST("/api/apply", controller.ApplyAPI, controller.AuthAPI)
	// 자신의 신청내역 가져오기
	e.GET("/api/getApplys", controller.GetApplysAPI, controller.AuthAPI)
	// 구역 신청내역 가져오기
	e.GET("/api/getApplysByArea", controller.GetApplysByAreaAPI, controller.AuthAPI)
	// 시간대에 해당하는 인원 수 가져오기
	e.GET("/api/getApplyMount", controller.GetApplyMountAPI, controller.AuthAPI)
	// 구역 신청 인원 수 가져오기
	e.GET("/api/getApplyMountByArea", controller.GetApplyMountByAreaAPI, controller.AuthAPI)
	// 신청 한도가 넘은 자율관 신청 시간들을 가져오기
	e.GET("/api/getDatesByOverCount", controller.GetDatesByOverCountAPI, controller.AuthAPI)
	// 신청 취소하기
	e.POST("/api/cancelApply", controller.CancelApplyAPI, controller.AuthAPI)

	// 공휴일 정보 가져오기
	e.GET("/api/getHolydays", controller.GetHolydaysAPI, controller.AuthAPI)

	// 공지 가져오기
	e.GET("/api/getNotices", controller.GetNoticesAPI, controller.AuthAPI)
	// id 값의 공지 가져오기
	e.GET("/api/getNoticeByID", controller.GetNoticeByIDAPI, controller.AuthAPI)

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

	// 신청현황 - 시간대 선택
	a.GET("/apply/", controller.AdminSelectTime)
	// 신청현황 - 면학실 선택
	a.GET("/apply/selectForm", controller.AdminSelectForm)
	// 신청현황 - 창학관 구역 선택
	a.GET("/apply/selectArea", controller.AdminSelectArea)
	// 신청현황 - 창학관 신청현황 보기
	a.GET("/apply/a-viewA", controller.AdminAViewA)
	a.GET("/apply/a-viewB", controller.AdminAViewB)
	a.GET("/apply/a-viewC", controller.AdminAViewC)
	a.GET("/apply/a-viewD", controller.AdminAViewD)
	a.GET("/apply/a-viewE", controller.AdminAViewE)
	a.GET("/apply/a-viewF", controller.AdminAViewF)
	a.GET("/apply/a-viewG", controller.AdminAViewG)
	// 신청현황 - 자율관 신청현황 보기
	a.GET("/apply/b-view", controller.AdminBView)

	// 공휴일 관리
	a.GET("/holydays", controller.AdminHolydays)

	// ================ 교사 API ================
	// 당일의 학생들 신청내역 가져오기 by period, form, area
	a.GET("/api/getApplys", controller.AdminGetApplysAPI)
	// 당일의 학생들 신청 인원 가져오기 by period, form
	a.GET("/api/getApplyMount", controller.AdminGetApplyMountAPI)

	// 공휴일 추가
	a.POST("/api/addHolyday", controller.AdminAddHolydayAPI)
	// 공휴일 삭제
	a.POST("/api/cancelHolyday", controller.AdminCancelHolydayAPI)
	// 공휴일 모든 정보 가져오기
	a.GET("/api/getAllHolydays", controller.AdminGetAllHolydaysAPI)
	// 공휴일 일주일 치 정보 가져오기
	a.GET("/api/getHolydays", controller.AdminGetHolydaysAPI)

	// Start web server
	e.Start(":80")
}
