package controllers

import (
	"net/http"
	"poetryAdmin/master/library/config"
	"poetryAdmin/master/library/tools"
)

func DisplayLayOut(w http.ResponseWriter, fileName string, data interface{}, layout []string) (err error) {
	tpl := tools.NewTpl()
	tpl.ViewPath = config.G_Conf.ViewDir
	if len(layout) == 0 {
		layout = []string{
			"/public/header.html",
		}
	}
	tpl.FileName = fileName
	tpl.LayOutFiles = layout
	tpl.Data = data
	tpl.Writer = w
	err = tpl.Display()
	return
}
