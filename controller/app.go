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

// Periods json models
type Periods struct {
	Form  string   `json:"form"`
	Times []string `json:"times"`
}

// SelectTimePOST : 시간대 선택 완료
func SelectTimePOST(c echo.Context) error {
	var periods Periods
	if err := c.Bind(&periods); err != nil {
		return err
	}

	if periods.Form == "A" {
		// 창학관 신청

		// 신청 시간으로 redirect
		return c.Redirect(http.StatusMovedPermanently, "/apply/selectArea?date="+models.GetTimeTableDays()[index].Format("2006-01-02")+"&period="+period)
	} else {
		// 자율관 신청

		days := models.GetTimeTableDays()
		// mon-7,tue-EP1 형식 -> , 단위로 구분
		for _, period := range periods.Times {
			// mon-7 형식 -> mon과 7로 분리
			temp := strings.Split(period, "-")
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

			session := session.Default(c)
			models.AddApply(session.Get("studentNumber").(string), session.Get("name").(string), days[index], period, "B", "", "")
		}

		return c.Redirect(http.StatusMovedPermanently, "/apply/success")
	}
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

	err = models.AddApply(session.Get("studentNumber").(string), session.Get("name").(string), day, c.FormValue("period"), c.FormValue("form"), c.FormValue("area"), c.FormValue("seat"))
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

// GetApplysAPI 자신의 신청내역 가져오기
func GetApplysAPI(c echo.Context) error {
	return c.JSON(http.StatusOK, models.GetApplysByStudentNumber(session.Default(c).Get("studentNumber").(string)))
}

// GetApplyMountOfAreaAPI 구역 신청 인원 수 가져오기
func GetApplyMountOfAreaAPI(c echo.Context) error {
	day, err := time.Parse("2006-01-02", c.QueryParam("date"))
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}
	return c.String(http.StatusOK, strconv.Itoa(models.GetApplyMountOfArea(day, c.QueryParam("period"), c.QueryParam("area"))))
}

// CancelApplyAPI 신청 취소
func CancelApplyAPI(c echo.Context) error {
	session := session.Default(c)
	day, err := time.Parse("2006-01-02", c.FormValue("date"))
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	err = models.DeleteApply(session.Get("studentNumber").(string), day, c.FormValue("period"))
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

// ChangePassword : ChangePassword page
func ChangePassword(c echo.Context) error {
	return c.Render(http.StatusOK, "changePassword", nil)
}

// ChangePasswordPost : Check a Password and change
func ChangePasswordPost(c echo.Context) error {
	// 새로운 비밀번호와 새로운 비밀번호 확인이 다를 때
	if c.FormValue("newPassword") != c.FormValue("newPasswordCheck") {
		return c.Redirect(http.StatusMovedPermanently, "/user/changePassword?status=equalError")
	}
	session := session.Default(c)
	err := models.ChangePassword(session.Get("studentNumber").(string), c.FormValue("loginPassword"), c.FormValue("newPassword"))
	// 비번 변경 실패
	if err != nil {
		return c.Redirect(http.StatusMovedPermanently, "/user/changePassword?status="+err.Error())
	}
	// 비밀번호 변경 성공
	return c.Redirect(http.StatusMovedPermanently, "/user/changePassword?status=success")

}

// GetHolydaysAPI : 공휴일 정보 가져오기 API
func GetHolydaysAPI(c echo.Context) error {
	return c.JSON(http.StatusOK, models.GetTimeTableHolydays())
}

// GetApplyMountAPI : 시간대에 해당하는 인원 수 가져오기
func GetApplyMountAPI(c echo.Context) error {
	day, err := time.Parse("2006-01-02", c.QueryParam("date"))
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}
	return c.String(http.StatusOK, strconv.Itoa(models.GetApplyMount(day, c.QueryParam("period"), c.QueryParam("form"))))
}
