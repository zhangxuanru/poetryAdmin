package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/master/library/config"
	"poetryAdmin/master/library/logger"
	"poetryAdmin/master/library/server"
	"runtime"
)

var confFile string

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initConfFile() {
	dir, _ := os.Getwd()
	confFile = dir + "/master/conf/conf.json"
	logrus.Info("加载配置文件:",confFile)
}

func main() {
	var (
		err error
	)
	initEnv()
	logger.InitLogger()
	initConfFile()
	if err = config.InitConfig(confFile); err != nil {
		goto PRINTERR
	}
	if err = server.InitHttpServer(); err != nil {
		goto PRINTERR
	}
	return
PRINTERR:
	logrus.Debug("err:",err)
	fmt.Println(err)
	return
}
