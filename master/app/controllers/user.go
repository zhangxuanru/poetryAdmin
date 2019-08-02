package controllers

import (
	"net/http"
	"poetryAdmin/master/app/logic"
)

func Login(w http.ResponseWriter, req *http.Request) {
	var (
		loginStatus bool
		userLogic   *logic.UserLogin
	)
	if req.Method == http.MethodGet {
		//userLogic = logic.NewUserLogin("", "", w)
		//userLogic.DelLoginCookie()
		err := base.DisplayHtmlLayOut(w, "login.html", nil, nil)
		base.DisplayErrorHtml(w, err)
		return
	}
	userName := req.PostFormValue("username")
	passWord := req.PostFormValue("password")
	userLogic = logic.NewUserLogin(userName, passWord, w)
	if err := userLogic.ValidateLogin(); err != nil {
		base.OutPutRespJson(w, nil, err.Error(), logic.RespFail)
		return
	}
	if loginStatus = userLogic.Login(); loginStatus == true {
		userLogic.WriteLoginCookie()
		base.OutPutRespJson(w, nil, logic.RespLoginSuccess, logic.RespSuccess)
	}
	if loginStatus == false {
		base.OutPutRespJson(w, nil, logic.RespLoginFailMsg, logic.RespFail)
	}
	return
}
