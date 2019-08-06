package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/logger"
	"poetryAdmin/worker/app/redis"
	"poetryAdmin/worker/core/grabs"
	"reflect"
	"runtime"
	"strings"
	"time"
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
    if err = grabs.Run();err!=nil{
		goto PRINTERR
	}
    for{
         time.Sleep(10 * time.Second)
	}
	return
 PRINTERR:
 	logrus.Debug("err:",err)
	return
}


func test1()  {
	reply, err := redis.Set("test3", "4", "EX", "3600", "NX")

	logrus.Info("test1 reply:",reply)
	logrus.Info("test1 err:",err)

	s := reflect.ValueOf(reply).String()
	lower := strings.ToLower(s)
	logrus.Info("lower:",lower)

}
