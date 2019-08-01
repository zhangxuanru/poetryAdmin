package tools

import (
	"fmt"
	"html/template"
	"net/http"
)

type Tpl struct {
	FileName    string
	LayOutFiles []string
	Data        interface{}
	Writer      http.ResponseWriter
	ViewPath    string
}

func (t *Tpl) Display() (err error) {
	t.LayOutFiles = append(t.LayOutFiles, t.FileName)
	for key, fileName := range t.LayOutFiles {
		t.LayOutFiles[key] = t.ViewPath + fileName
	}
	fmt.Println(t.LayOutFiles)

	must := template.Must(template.ParseFiles(t.LayOutFiles...))
	err = must.Execute(t.Writer, t.Data)
	return err
}

func NewTpl() *Tpl {
	return new(Tpl)
}
