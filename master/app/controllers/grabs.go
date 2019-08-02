package controllers

import "net/http"

//一键抓取列表
func Grabs(w http.ResponseWriter, r *http.Request) {
	base.DisplayHtmlLayOut(w, "grab-list.html", nil, nil)
}

//执行抓取, 写入etcd
func GrabsImpl(w http.ResponseWriter, r *http.Request) {

}
