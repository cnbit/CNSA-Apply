package controller

import (
	"CNSA-Apply/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	session "github.com/ipfans/echo-session"
	"github.com/labstack/echo"
)

// AuthAPI 로그인 인증 middleware
func AuthAPI(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 로그인이 되어 있지 않으면 login page로 redirect
		session := session.Default(c)
		if session.Get("studentNumber") == nil {
			return c.Redirect(http.StatusMovedPermanently, "/login")
		}

		return next(c)
	}
}

// Login : Login Page
func Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

// LoginPost : Check a Login Data
func LoginPost(c echo.Context) error {
	isSuccessed, name, gender := models.Login(c.FormValue("loginID"), c.FormValue("loginPassword"))

	// Login 성공 시
	if isSuccessed {
		// Session에 학번 저장
		session := session.Default(c)
		session.Set("studentNumber", c.FormValue("loginID"))
		session.Set("name", name)
		session.Set("gender", gender)
		session.Save()

		return c.Redirect(http.StatusMovedPermanently, "/")
	}

	// Login 실패 시
	return c.Redirect(http.StatusMovedPermanently, "/login?error=loginFailed")
}

// Logout : 로그아웃 - 세션 초기화
func Logout(c echo.Context) error {
	// Session 초기화
	session := session.Default(c)
	session.Clear()
	session.Save()

	// 로그인 페이지로 빠이빠이
	return c.Redirect(http.StatusMovedPermanently, "/login")
}

// Index : Main Page
func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

// SelectForm : 신청하기 - 면학실 선택
func SelectForm(c echo.Context) error {
	return c.Render(http.StatusOK, "selectForm", nil)
}

// SelectTime : 신청하기 - 시간대 선택
func SelectTime(c echo.Context) error {
	return c.Render(http.StatusOK, "selectTime", nil)
}

// SelectArea : 신청하기 - 구역 선택
func SelectArea(c echo.Context) error {
	session := session.Default(c)
	return c.Render(http.StatusOK, "selectArea", map[string]interface{}{
		"gender": session.Get("gender").(int),
	})
}

// SelectAreaPOST : 신청하기 - 구역 선택
func SelectAreaPOST(c echo.Context) error {
	// mon-7 형식 -> mon과 7로 분리
	temp := strings.Split(c.FormValue("time"), "-")
	day, period := temp[0], temp[1]

	index := 0
	// 요일별로 구분
	if day == "mon" {
		index = 0
	} else if day == "tue" {
		index = 1
	} else if day == "wed" {
		index = 2
	} else if day == "thr" {
		index = 3
	} else if day == "fri" {
		index = 4
	}

	area := c.FormValue("area")
	if area == "" {
		// 구역을 선택하지 않았을 경우
		return c.Redirect(http.StatusMovedPermanently, "/apply/selectArea")
	}

	return c.Redirect(http.StatusMovedPermanently, "/apply/selectSeat"+area+"?date="+models.GetTimeTableDays()[index].Format("2006-01-02")+"&period="+period)
}

// SelectSeatA : 신청하기 - 좌석 선택
func SelectSeatA(c echo.Context) error {
	session := session.Default(c)
	return c.Render(http.StatusOK, "selectSeatA", map[string]interface{}{
		"gender": session.Get("gender").(int),
	})
}

// SelectSeatB : 신청하기 - 좌석 선택
func SelectSeatB(c echo.Context) error {
	session := session.Default(c)
	return c.Render(http.StatusOK, "selectSeatB", map[string]interface{}{
		"gender": session.Get("gender").(int),
	})
}

// SelectSeatC : 신청하기 - 좌석 선택
func SelectSeatC(c echo.Context) error {
	session := session.Default(c)
	return c.Render(http.StatusOK, "selectSeatC", map[string]interface{}{
		"gender": session.Get("gender").(int),
	})
}

// SelectSeatD : 신청하기 - 좌석 선택
func SelectSeatD(c echo.Context) error {
	session := session.Default(c)
	return c.Render(http.StatusOK, "selectSeatD", map[string]interface{}{
		"gender": session.Get("gender").(int),
	})
}

// SelectSeatE : 신청하기 - 좌석 선택
func SelectSeatE(c echo.Context) error {
	session := session.Default(c)
	return c.Render(http.StatusOK, "selectSeatE", map[string]interface{}{
		"gender": session.Get("gender").(int),
	})
}

// SelectSeatF : 신청하기 - 좌석 선택
func SelectSeatF(c echo.Context) error {
	session := session.Default(c)
	return c.Render(http.StatusOK, "selectSeatF", map[string]interface{}{
		"gender": session.Get("gender").(int),
	})
}

// SelectSeatG : 신청하기 - 좌석 선택
func SelectSeatG(c echo.Context) error {
	session := session.Default(c)
	return c.Render(http.StatusOK, "selectSeatG", map[string]interface{}{
		"gender": session.Get("gender").(int),
	})
}

// ApplySuccess : 신청하기 - 신청완료
func ApplySuccess(c echo.Context) error {
	return c.Render(http.StatusOK, "applySuccess", nil)
}

// ApplyAPI 신청정보 등록
func ApplyAPI(c echo.Context) error {
	session := session.Default(c)
	day, err := time.Parse("2006-01-02", c.FormValue("date"))
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	err = models.AddApply(session.Get("studentNumber").(string), session.Get("name").(string), session.Get("gender").(int), day, c.FormValue("period"), c.FormValue("form"), c.FormValue("area"), c.FormValue("seat"))
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

// GetApplysAPI 자신의 신청내역 가져오기
func GetApplysAPI(c echo.Context) error {
	return c.JSON(http.StatusOK, models.GetApplysByStudentNumber(session.Default(c).Get("studentNumber").(string)))
}

// GetApplysByAreaAPI 구역 신청내역 가져오기
func GetApplysByAreaAPI(c echo.Context) error {
	day, err := time.Parse("2006-01-02", c.QueryParam("date"))
	if err != nil {
		// TODO: http code 추후 정리
		return c.String(http.StatusOK, err.Error())
	}
	return c.JSON(http.StatusOK, models.GetApplys(day, c.QueryParam("period"), c.QueryParam("form"), c.QueryParam("area")))
}

// GetApplyMountAPI : 시간대에 해당하는 인원 수 가져오기
func GetApplyMountAPI(c echo.Context) error {
	day, err := time.Parse("2006-01-02", c.QueryParam("date"))
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}
	return c.String(http.StatusOK, strconv.Itoa(models.GetApplyMount(day, c.QueryParam("period"), c.QueryParam("form"))))
}

// GetApplyMountByAreaAPI 구역 신청 인원 수 가져오기
func GetApplyMountByAreaAPI(c echo.Context) error {
	day, err := time.Parse("2006-01-02", c.QueryParam("date"))
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}
	return c.String(http.StatusOK, strconv.Itoa(models.GetApplyMountByArea(day, c.QueryParam("period"), c.QueryParam("area"))))
}

// GetDatesByOverCountAPI 신청 한도가 넘은 자율관 신청시간들을 가져오기
func GetDatesByOverCountAPI(c echo.Context) error {
	times := models.GetDatesByOverCount(session.Default(c).Get("gender").(int))

	return c.JSON(http.StatusOK, times)
}

// CancelApplyAPI 신청 취소
func CancelApplyAPI(c echo.Context) error {
	session := session.Default(c)
	day, err := time.Parse("2006-01-02", c.FormValue("date"))
	if err != nil {
		// TODO: http code 추후 정리
		return c.String(http.StatusOK, err.Error())
	}

	err = models.DeleteApply(session.Get("studentNumber").(string), day, c.FormValue("period"))
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

// MyPage : 내 정보 page
func MyPage(c echo.Context) error {
	return c.Render(http.StatusOK, "myPage", nil)
}

// ApplyHistory : 신청내역 page
func ApplyHistory(c echo.Context) error {
	return c.Render(http.StatusOK, "applyHistory", nil)
}

// Account : ChangePassword page
func Account(c echo.Context) error {
	return c.Render(http.StatusOK, "account", nil)
}

// AccountPOST : Check a Password and change
func AccountPOST(c echo.Context) error {
	session := session.Default(c)
	err := models.ChangePassword(session.Get("studentNumber").(string), c.FormValue("loginPassword"), c.FormValue("newPassword"))
	// 비번 변경 실패
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}
	// 비밀번호 변경 성공
	return c.String(http.StatusOK, "success")
}

// ChangeSuccess : 변경 완료 페이지
func ChangeSuccess(c echo.Context) error {
	return c.Render(http.StatusOK, "changeSuccess", nil)
}

// GetHolydaysAPI : 공휴일 정보 가져오기 API
func GetHolydaysAPI(c echo.Context) error {
	return c.JSON(http.StatusOK, models.GetTimeTableHolydays())
}
