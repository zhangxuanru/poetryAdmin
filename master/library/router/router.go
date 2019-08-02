package router

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"poetryAdmin/master/app/controllers"
	"poetryAdmin/master/app/logic"
	"strings"
)

func InitRouter(mux *http.ServeMux) {
	mux.HandleFunc("/admin", CallAction(controllers.Index))     //后台首页
	mux.HandleFunc("/welcome", CallAction(controllers.WelCome)) //欢迎页
	mux.HandleFunc("/login", CallAction(controllers.Login))     //登录，退出
	mux.HandleFunc("/grabs", CallAction(controllers.Grabs))     //一键抓取
}

func CallAction(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Info("recover err:", err)
			}
		}()
		before(writer, request)
		handlerFunc(writer, request)
	}
}

//记录请求日志,以后完善
func before(writer http.ResponseWriter, request *http.Request) {
	if request.RequestURI == "/login" {
		return
	}
	adminCookie, _ := request.Cookie(logic.LoginCookieName)
	userCookie, _ := request.Cookie(logic.LoginCookieUserName)
	passCookie, _ := request.Cookie(logic.LoginCookiePassword)
	if adminCookie.Value == "" || userCookie.Value == "" || passCookie.Value == "" {
		http.Redirect(writer, request, "/login", 301)
	}
	login := logic.NewUserLogin(userCookie.Value, passCookie.Value, writer)
	if strings.Compare(login.LoginDataMd5(), adminCookie.Value) == 0 {
		return
	}
	logrus.Debug("login 301....", "adminCookie:", adminCookie)
	http.Redirect(writer, request, "/login", 301)
}
