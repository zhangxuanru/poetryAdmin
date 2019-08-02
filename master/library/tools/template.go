package tools

import (
	"html/template"
	"net/http"
)

type Tpl struct {
	FileName    string
	LayOutFiles []string
	Data        interface{}
	Writer      http.ResponseWriter
	ViewPath    string
	LayOutPath  string
}

func (t *Tpl) Display() (err error) {
	for key, fileName := range t.LayOutFiles {
		//此处注意， 如果在方法中修改属性数据下标对应的值，则外面的值也会改变***
		t.LayOutFiles[key] = t.LayOutPath + fileName
	}
	t.LayOutFiles = append(t.LayOutFiles, t.ViewPath+t.FileName)
	must := template.Must(template.ParseFiles(t.LayOutFiles...))
	err = must.ExecuteTemplate(t.Writer, t.FileName, t.Data)
	return err
}

func NewTpl() *Tpl {
	return new(Tpl)
}
