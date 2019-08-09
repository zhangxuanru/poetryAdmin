package test

import (
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/logger"
	"poetryAdmin/worker/app/models"
	"poetryAdmin/worker/app/redis"
	"strings"
)

var confFile string

func initConfigFile() (err error) {
	var (
		dir string
	)
	if dir, err = os.Getwd(); err != nil {
		return err
	}
	dir = strings.TrimRight(dir, "/test")
	confFile = dir + "/config/config.json"
	return
}

func init() {
	var err error
	if err = initConfigFile(); err != nil {
		goto PRINTERR
	}
	logger.InitLogger()
	if err = config.InitConfig(confFile); err != nil {
		goto PRINTERR
	}
	if err = redis.InitRedis(config.G_Conf.RedisHost); err != nil {
		goto PRINTERR
	}
	if err = models.InitDb(); err != nil {
		goto PRINTERR
	}
PRINTERR:
	logrus.Debug("err:", err)
	return
}
