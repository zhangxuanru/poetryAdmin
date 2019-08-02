package router

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"poetryAdmin/master/app/controllers"
)

func InitRouter(mux *http.ServeMux) {
	mux.HandleFunc("/admin", CallAction(controllers.Index))
	mux.HandleFunc("/welcome", CallAction(controllers.WelCome))
	mux.HandleFunc("/login", CallAction(controllers.Login))
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
	//return
	//if request.RequestURI == "/login" {
	//	return
	//}
	//logrus.Info("----before---")
	//adminCookie, _ := request.Cookie(logic.LoginCookieName)
	//userCookie, _ := request.Cookie(logic.LoginCookieUserName)
	//passCookie, _ := request.Cookie(logic.LoginCookiePassword)
	//if adminCookie.Value == "" || userCookie.Value == "" || passCookie.Value == "" {
	//	logrus.Info("cookie err:", adminCookie, "--", userCookie, "--", passCookie)
	//	http.Redirect(writer, request, "/login", 301)
	//}
	//login := logic.NewUserLogin(userCookie.Value, passCookie.Value, writer)
	//if strings.Compare(login.LoginDataMd5(), adminCookie.Value) == 0 {
	//	logrus.Info("md5==md5")
	//	//return
	//}
	//logrus.Info("md5:", login.LoginDataMd5(), "adminc:", adminCookie.Value)
	//http.Redirect(writer, request, "/login", 301)
}
