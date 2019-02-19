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

//Login 학생 아이디 인증(SALT)
func Login(studentNumber string, password string) bool {
	user := User{}
	err := db.Table("users").Where("student_number = ?", studentNumber).First(&user).Error

	if err != nil {
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+SALT)) == nil
}

func ChangePW(studentNumber string, newPW string) error {
	if len(newPW) > 95 {
		err := errors.New("Exceed the length")
		return err
	} else if newPW == "" {
		err := errors.New("PW is empty")
		return err
	}
	user := User{}
	db.Table("users").Where("student_number = ?", studentNumber).First(&user)
	bytes, _ := bcrypt.GenerateFromPassword([]byte(newPW+SALT), 14)
	user.Password = string(bytes)

	return db.Save(&user).Error
}

//Tlogin 교사 아이디 인증(cnsanet)
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
