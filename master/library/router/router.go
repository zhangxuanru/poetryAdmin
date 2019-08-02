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
	var (
		cookieData *http.Cookie
		err        error
	)
	if request.RequestURI == "/login" {
		return
	}
	if cookieData, err = request.Cookie("poetry"); err != nil {
		http.Redirect(writer, request, "/login", 301)
		return
	}
	logrus.Info(cookieData.Value)
}
