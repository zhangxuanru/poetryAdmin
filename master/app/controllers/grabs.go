package controllers

import "net/http"

func Grabs(w http.ResponseWriter, r *http.Request) {
	base.DisplayHtmlLayOut(w, "grab-list.html", nil, nil)
}
