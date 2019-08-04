package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/master/library/config"
	"poetryAdmin/master/library/logger"
	"poetryAdmin/master/library/redis"
	"poetryAdmin/master/library/server"
	"poetryAdmin/master/library/validate"
	"runtime"
)

var confFile string

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initConfFile() {
	dir, _ := os.Getwd()
	confFile = dir + "/master/conf/conf.json"
	logrus.Info("加载配置文件:", confFile)
}

func init() {
	initEnv()
	initConfFile()
}

func main() {
	var (
		err error
	)
	logger.InitLogger()
	if err = config.InitConfig(confFile); err != nil {
		goto PRINTERR
	}
	if err = redis.InitRedis(config.G_Conf.RedisHost); err != nil {
		goto PRINTERR
	}
	if err = validate.InitValidate(); err != nil {
		goto PRINTERR
	}
	if err = server.InitHttpServer(); err != nil {
		goto PRINTERR
	}
	return
PRINTERR:
	logrus.Debug("err:", err)
	fmt.Println(err)
	return
}
