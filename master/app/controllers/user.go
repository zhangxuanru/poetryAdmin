package controllers

import "net/http"

func Login(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		err := base.DisplayHtmlLayOut(w, "login.html", nil, nil)
		base.PrintError(w, err)
		return
	}
	username := req.PostFormValue("username")
	password := req.PostFormValue("password")
}
