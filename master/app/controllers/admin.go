package controllers

import (
	"net/http"
)

var base *Base

func init() {
	base = initBase()
}

func Index(writer http.ResponseWriter, request *http.Request) {
	err := base.DisplayHtmlLayOut(writer, "index.html", nil, nil)
	base.PrintError(writer, err)
}

func WelCome(writer http.ResponseWriter, req *http.Request) {
	err := base.DisplayHtmlLayOut(writer, "welcome.html", nil, []string{
		"public/header.html",
	})
	base.PrintError(writer, err)
}
