package models

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

const (
	// SQLConnectionString : MySQL Connection String
	SQLConnectionString = "*"
	// SALT : SALT
	SALT = "*"
)

// Database Connection
var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", SQLConnectionString)
	if err != nil {
		panic(err)
	}
}

// Apply models
type Apply struct {
	StudentNumber string    `gorm:"type:VARCHAR(6); primary_key" json:"student-number"`
	Name          string    `gorm:"type:VARCHAR(50)" json:"name"`
	Date          time.Time `gorm:"type:DATE; primary_key; unique_index" json:"date"`
	Period        string    `gorm:"type:VARCHAR(5); primary_key; unique_index" json:"period"`
	Form          string    `gorm:"type:VARCHAR(1)" json:"form"`
	Area          string    `gorm:"type:VARCHAR(1)" json:"area"`
	Seat          string    `gorm:"type:VARCHAR(6); unique_index" json:"seat"`
}

// TableName of Apply
func (c *Apply) TableName() string {
	return "applys"
}

// User models
type User struct {
	StudentNumber string `gorm:"type:VARCHAR(6); primary_key" json:"student-number"`
	Password      string `gorm:"type:VARCHAR(100)" json:"password"`
	Name          string `gorm:"type:VARCHAR(50)" json:"name"`
	Gender        int    `gorm:"type:INT" json:"gender"`
}

// TableName of User
func (c *User) TableName() string {
	return "users"
}

// Holyday models
type Holyday struct {
	Date time.Time `gorm:"type:DATE; primary_key" json:"date"`
	Name string    `gorm:"type:VARCHAR(50)" json:"name"`
}

// TableName of Holyday
func (c *Holyday) TableName() string {
	return "holydays"
}

// Login 학생 아이디 인증(SALT)
func Login(studentNumber string, password string) (bool, string) {
	user := User{}
	err := db.Table("users").Where("student_number = ?", studentNumber).First(&user).Error
	if err != nil {
		return false, ""
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+SALT))
	if err != nil {
		return false, ""
	}

	return true, user.Name
}

// ChangePassword 비밀번호 변경
func ChangePassword(studentNumber string, password string, newPassword string) error {
	if len(newPassword) > 30 {
		// 새로운 비밀번호의 길이가 길 때
		err := errors.New("Exceed the length")
		return err
	} else if newPassword == "" {
		// 아무것도 입력하지 않았을 때
		err := errors.New("NewPassword is empty")
		return err
	}
	user := User{}
	db.Table("users").Where("student_number = ?", studentNumber).First(&user)

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+SALT)) != nil {
		// db의 현재 비밀번호와 입력된 현재 비밀번호가 일치하지 않을 때
		return errors.New("Password is incorrect")
	}

	bytes, _ := bcrypt.GenerateFromPassword([]byte(newPassword+SALT), bcrypt.DefaultCost)
	user.Password = string(bytes)

	return db.Save(&user).Error
}

// TcrLogin 교사 아이디 인증(cnsanet)
func TcrLogin(loginID string, loginPW string) bool {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, _ := client.PostForm("https://school.cnsa.hs.kr/login/userLogin", url.Values{
		"loginId": {loginID},
		"loginPw": {loginPW},
	})

	if resp.Request.Method == "GET" {
		return true
	}
	return false
}

// GetTimeTableDays : 페이지에 표시할 5일을 반환
// 월~금: 금주 월~금을 반환
// 토~일: 다음주 월~금을 반환
func GetTimeTableDays() [5]time.Time {
	var days [5]time.Time
	now := time.Now()

	if now.Weekday() > 5 {
		// 주말일 경우 시작일이 다음주 월요일
		days[0] = time.Date(now.Year(), now.Month(), now.Day()+(8-int(now.Weekday())), 0, 0, 0, 0, time.Local)
	} else {
		// 평일일 경우 시작일이 금주 월요일
		days[0] = time.Date(now.Year(), now.Month(), now.Day()+(1-int(now.Weekday())), 0, 0, 0, 0, time.Local)
	}
	days[1] = days[0].AddDate(0, 0, 1)
	days[2] = days[1].AddDate(0, 0, 1)
	days[3] = days[2].AddDate(0, 0, 1)
	days[4] = days[3].AddDate(0, 0, 1)

	return days
}

// AddApply 비어있는 좌석에 신청
// 같은 사람이 같은 시간에 신청은 선택할 때 방지
// 발생 가능한 오류는 비슷한 시간대에 동일한 좌석에 신청
func AddApply(studentNumber string, name string, day time.Time, period string, form string, seat string) error {
	apply := Apply{
		StudentNumber: studentNumber,
		Name:          name,
		Date:          day,
		Period:        period,
		Form:          form,
		Seat:          seat,
	}

	err := db.Save(&apply).Error
	if err.Error()[:9] != "Error 1062" {
		err = errors.New("The seat has been applied")
	}

	return err
}

// GetApplysByStudentNumber 페이지에 표시될 5일에 해당하는 신청내역을 가져옴
func GetApplysByStudentNumber(studentNumber string) []Apply {
	applys := []Apply{}
	db.Table("applys").Where("student_number = ? AND date >= ?", studentNumber, GetTimeTableDays()[0]).Find(&applys)
	return applys
}

// DeleteApply 좌석 신청 정보 삭제
// 요청 시간이 시작 이후일 경우 에러 리턴
func DeleteApply(studentNumber string, day time.Time, period string) error {
	// studyDay: 신청한 면학시간
	var studyDay time.Time
	if period == "7" {
		// 7교시
		studyDay = time.Date(day.Year(), day.Month(), day.Day(), 3, 30, 0, 0, time.Local)
	} else if period == "CAS" {
		// CAS
		studyDay = time.Date(day.Year(), day.Month(), day.Day(), 4, 50, 0, 0, time.Local)
	} else if period == "EP1" {
		// EP1
		studyDay = time.Date(day.Year(), day.Month(), day.Day(), 7, 20, 0, 0, time.Local)
	} else if period == "EP2" {
		// EP2
		studyDay = time.Date(day.Year(), day.Month(), day.Day(), 9, 00, 0, 0, time.Local)
	}

	if time.Now().After(studyDay) {
		// 지금이 면학신청한 시간 이후면 Error 반환
		return errors.New("It's already over")
	}

	err := db.Table("applys").Where("student_number = ? AND date = ? AND period = ?", studentNumber, day.Format("2006-01-02"), period).Delete(Apply{}).Error
	return err
}

// GetTimeTableHolydays 페이지에 표시될 5일에 해당하는 공휴일 정보를 가져옴
func GetTimeTableHolydays() []Holyday {
	holydays := []Holyday{}
	db.Table("holydays").Where("date >= ? AND date <= ?", GetTimeTableDays()[0], GetTimeTableDays()[4].Format("2006-01-02")).Find(&holydays)
	return holydays
}

// DeleteHolyday 공휴일을 삭제함
func DeleteHolyday(holyday time.Time) error {
	return db.Table("holydays").Where("date = ?", holyday).Delete(Holyday{}).Error
}

// GetApplyMount 특정 시간의 신청 수를 반환함
func GetApplyMount(day time.Time, period string, form string) int {
	var count int
	db.Table("applys").Where("date = ? AND period = ? AND form = ?", day.Format("2006-01-02"), period, form).Count(&count)
	return count
}

// AddHolyday 공휴일 추가
// error 반환
func AddHolyday(day time.Time, name string) error {
	return db.Save(&Holyday{Date: day, Name: name}).Error
}

// GetApplys 신청내역 확인
func GetApplys(day time.Time, period string, form string) []Apply {
	applys := []Apply{}
	db.Table("applys").Where("date = ? AND period = ? AND form = ?", day.Format("2006-01-02"), period, form).Find(&applys)

	return applys
}

// GetHolydays 모든 공휴일 정보 가져오기
func GetHolydays() []Holyday {
	holydays := []Holyday{}
	db.Table("holydays").Where("date >= ?", time.Now().Format("2006-01-02")).Find(&holydays)

	return holydays
}
