package subscribe

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/redis"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/parse"
)

//订阅redis并开始执行抓取
func InitGrasp() (err error) {
	subReceiveChan := make(chan []byte)
	go parse.NewAnalysis(subReceiveChan).ParseSubscribeData()
	go data.NewGraspResult().PrintMsg()
	if err = redis.SubScribe(config.G_Conf.PubChannelTitle, subReceiveChan); err != nil {
		logrus.Debug("SubScribe err:", err)
		return err
	}
	return nil
}
