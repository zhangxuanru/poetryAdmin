package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/logger"
	"poetryAdmin/worker/app/models"
	"poetryAdmin/worker/app/redis"
	"poetryAdmin/worker/core/subscribe"
	"runtime"
)

var confFile string

func initConfigFile() (err error) {
	var (
		dir string
	)
	if dir, err = os.Getwd();err!=nil{
        return err
	}
	confFile = dir + "/worker/config/config.json"
	return
}

func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
    var err error
    if err = initConfigFile();err!=nil{
    	goto PRINTERR
	}
	logger.InitLogger()
	if err = config.InitConfig(confFile);err!=nil{
		goto PRINTERR
	}
    if err = redis.InitRedis(config.G_Conf.RedisHost);err!=nil{
		goto PRINTERR
	}
    if err = models.InitDb();err!=nil{
    	goto PRINTERR
	}
    if err = subscribe.InitGrasp();err!=nil{
		goto PRINTERR
	}
    logrus.Info("end .... worker")
	return
 PRINTERR:
 	logrus.Debug("err:",err)
	return
}
