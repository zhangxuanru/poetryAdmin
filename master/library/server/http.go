package server

import (
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"poetryAdmin/master/library/config"
	"poetryAdmin/master/library/router"
	"strconv"
	"time"
)

func InitHttpServer() (err error) {
	logrus.Info("开始启动HTTP服务,监听端口是:", config.G_Conf.ApiPort)
	var (
		serverMux *http.ServeMux
		server    *http.Server
	)
	serverMux = http.NewServeMux()
	router.InitRouter(serverMux)

	logger := &logrus.Logger{}
	w := logger.Writer()
	defer w.Close()
	server = &http.Server{
		Addr:         ":" + strconv.Itoa(config.G_Conf.ApiPort),
		Handler:      serverMux,
		ReadTimeout:  time.Duration(config.G_Conf.ApiReadTimeOut) * time.Millisecond,
		WriteTimeout: time.Duration(config.G_Conf.ApiWriteTimeOut) * time.Millisecond,
		ErrorLog:     log.New(w, "", 0),
	}
	err = server.ListenAndServe()
	if err != nil {
		logrus.Debug("启动HTTP服务错误:", err)
	}
	return err
}
