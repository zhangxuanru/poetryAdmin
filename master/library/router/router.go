package router

import (
	"net/http"
	"poetryAdmin/master/app/controllers"
)

func InitRouter(mux *http.ServeMux) {
	mux.HandleFunc("/admin", CallAction(controllers.Index))
}

func CallAction(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		go before()
		handlerFunc(writer, request)
	}
}

//记录请求日志,以后完善
func before() {

}
