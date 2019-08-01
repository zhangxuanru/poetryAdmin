package controllers

import (
	"fmt"
	"net/http"
	"poetryAdmin/master/library/config"
)

func Index(writer http.ResponseWriter, request *http.Request) {

	//file := config.G_Conf.ViewDir + "index.html"
	//must := template.Must(template.ParseFiles(file))
	//
	data := make(map[string]interface{})
	data["static"] = config.G_Conf.StaticDomain

	err := DisplayLayOut(writer, "index.html", data, nil)

	fmt.Println("-----test----", err)
	//must.Execute(writer, data)
}
