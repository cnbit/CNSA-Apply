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
	Time          string    `gorm:"type:VARCHAR(5); primary_key; unique_index" json:"time"`
	Type          string    `gorm:"type:VARCHAR(1)" json:"type"`
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
func Login(studentNumber string, password string) bool {
	user := User{}
	err := db.Table("users").Where("student_number = ?", studentNumber).First(&user).Error

	if err != nil {
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+SALT)) == nil
}

// ChangePassword 비밀번호 변경
func ChangePassword(studentNumber string, newPassword string) error {
	if len(newPassword) > 30 {
		err := errors.New("Exceed the length")
		return err
	} else if newPassword == "" {
		err := errors.New("PW is empty")
		return err
	}
	user := User{}
	db.Table("users").Where("student_number = ?", studentNumber).First(&user)
	bytes, _ := bcrypt.GenerateFromPassword([]byte(newPassword+SALT), bcrypt.DefaultCost)
	user.Password = string(bytes)

	return db.Save(&user).Error
}

// Tlogin 교사 아이디 인증(cnsanet)
func Tlogin(loginID string, loginPW string) bool {
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
		days[0] = now.AddDate(0, 0, 8-int(now.Weekday()))
	} else {
		days[0] = now.AddDate(0, 0, 1-int(now.Weekday()))
	}
	days[1] = days[0].AddDate(0, 0, 1)
	days[2] = days[1].AddDate(0, 0, 1)
	days[3] = days[2].AddDate(0, 0, 1)
	days[4] = days[3].AddDate(0, 0, 1)

	return days
}
