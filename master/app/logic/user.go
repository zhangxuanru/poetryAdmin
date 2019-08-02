package logic

import (
	"errors"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"poetryAdmin/master/library/tools"
	"poetryAdmin/master/library/validate"
)

type UserLogin struct {
	UserName string `validate:"required,gte=5,lte=20"`
	PassWord string `validate:"required,gte=5,lte=20"`
	w        http.ResponseWriter
}

func NewUserLogin(userName, passWord string, w http.ResponseWriter) *UserLogin {
	return &UserLogin{UserName: userName, PassWord: passWord, w: w}
}

//数据验证
func (u *UserLogin) ValidateLogin() (err error) {
	var (
		errs validator.ValidationErrors
	)
	if err = validate.ValidateStruct(u); err != nil {
		errs = validate.ErrToValidationErrors(err)
	}
	for _, err := range errs {
		return errors.New(err.Translate(validate.G_Translator))
	}
	return nil
}

//数据验证，登录
func (u *UserLogin) Login() (ret bool) {
	if u.UserName == "admin" && u.PassWord == "abc123456" {
		return true
	}
	return false
}

//根据登录信息生成加密串
func (u *UserLogin) LoginDataMd5() string {
	loginStr := fmt.Sprintf("%s:%s", u.UserName, u.PassWord)
	return tools.Md5(loginStr)
}

//登录记录cookie
func (u *UserLogin) WriteLoginCookie() {
	u.PassWord = tools.Md5(u.PassWord)
	cookie := http.Cookie{
		Name:  LoginCookieName,
		Value: u.LoginDataMd5(),
		Path:  "/",
	}
	userCookie := http.Cookie{
		Name:  LoginCookieUserName,
		Value: u.UserName,
		Path:  "/",
	}
	passCookie := http.Cookie{
		Name:  LoginCookiePassword,
		Value: u.PassWord,
		Path:  "/",
	}
	http.SetCookie(u.w, &cookie)
	http.SetCookie(u.w, &userCookie)
	http.SetCookie(u.w, &passCookie)
}

//删除登录的cookie
func (u *UserLogin) DelLoginCookie() {
	cookie := http.Cookie{
		Name:   LoginCookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	userCookie := http.Cookie{
		Name:   LoginCookieUserName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	passCookie := http.Cookie{
		Name:   LoginCookiePassword,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(u.w, &cookie)
	http.SetCookie(u.w, &userCookie)
	http.SetCookie(u.w, &passCookie)
}
