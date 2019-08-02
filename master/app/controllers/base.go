package controllers

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"poetryAdmin/master/app/logic"
	"poetryAdmin/master/library/config"
	"poetryAdmin/master/library/tools"
)

type Base struct {
	Err    error
	Writer http.ResponseWriter
}

func initBase() *Base {
	return new(Base)
}

//显示模板并且加上layout页面
func (b *Base) DisplayHtmlLayOut(w http.ResponseWriter, fileName string, data map[string]interface{}, layout []string) (err error) {
	var layOutFiles = []string{
		"public/header.html",
		"public/footer.html",
	}
	tpl := &tools.Tpl{
		ViewPath:   config.G_Conf.ViewDir,
		LayOutPath: config.G_Conf.ViewDir,
		FileName:   fileName,
		Writer:     w,
	}
	if len(layout) == 0 {
		layout = layOutFiles
	}
	if data == nil {
		data = make(map[string]interface{})
	}
	data["static"] = config.G_Conf.StaticDomain
	data["baseUrl"] = config.G_Conf.BaseDomain
	tpl.LayOutFiles = layout
	tpl.Data = data
	err = tpl.Display()
	return
}

//显示单独的页面，不带layout
func (b *Base) DisplayHtml(w http.ResponseWriter, fileName string, data interface{}) (err error) {
	tpl := &tools.Tpl{
		Writer:   w,
		Data:     data,
		FileName: fileName,
		ViewPath: config.G_Conf.ViewDir,
	}
	return tpl.Display()
}

//显示错误，后期可用于集中处理错误，像跳转等逻辑啥的
func (b *Base) DisplayErrorHtml(w http.ResponseWriter, err error) {
	if err != nil {
		layOut := []string{"public/header.html"}
		data := make(map[string]interface{})
		data["error"] = err.Error()
		b.DisplayHtmlLayOut(w, "error.html", data, layOut)
		return
	}
}

//输出json到response
func (b *Base) OutPutRespJson(w http.ResponseWriter, data interface{}, msg string, code int) {
	resp := logic.NewRespData(data, msg, code)
	err := b.OutPutJson(w, resp)
	if err != nil {
		logrus.Debug("err:", err, "data:", resp)
		logrus.Info("err:", err, "data:", resp)
	}
	return
}

//输出JSON
func (b *Base) OutPutJson(w http.ResponseWriter, data interface{}) (err error) {
	var outPutData []byte
	outPutData, err = config.G_Json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(outPutData)
	return err
}
